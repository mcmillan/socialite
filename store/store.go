package store

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

const (
	decay = 12 * time.Hour
)

var redisClient redis.Conn

type Link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int64  `json:"score"`
}

func (l *Link) ID() string {
	h := sha1.New()
	io.WriteString(h, strings.ToLower(l.URL))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SinceID() (sinceID string, err error) {
	verifyConnection()

	sinceID, err = redis.String(redisClient.Do("GET", "socialite:since_id"))

	if err == redis.ErrNil {
		sinceID = ""
		err = nil
	}

	return
}

func SetSinceID(sinceID string) (err error) {
	verifyConnection()

	_, err = redisClient.Do("SET", "socialite:since_id", sinceID)

	return
}

func Popular() (links []Link, err error) {
	verifyConnection()
	scoredLinks, err := redis.Values(redisClient.Do(
		"ZREVRANGEBYSCORE",
		"socialite:urls",
		"+inf",
		"-inf",
		"WITHSCORES",
		"LIMIT",
		0,
		10))

	if err != nil {
		log.WithFields(log.Fields{"package": "store"}).Error(err)
		return
	}

	for len(scoredLinks) > 0 {
		var (
			id    string
			url   string
			title string
			score int64
		)

		scoredLinks, err = redis.Scan(scoredLinks, &id, &score)

		if err != nil {
			log.WithFields(log.Fields{
				"package": "store",
				"id":      id,
			}).Error(err)
			return
		}

		url, err = redis.String(redisClient.Do("HGET", metadataKey(id), "url"))

		if err != nil {
			log.WithFields(log.Fields{
				"package": "store",
				"id":      id,
			}).Error(err)
			return
		}

		title, err = redis.String(redisClient.Do("HGET", metadataKey(id), "title"))

		if err != nil {
			log.WithFields(log.Fields{
				"package": "store",
				"id":      id,
			}).Error(err)
			return
		}

		links = append(links, Link{URL: url, Title: title, Score: score})
	}

	return
}

func Add(link Link) (err error) {
	verifyConnection()

	expiryDate := time.Now().Add(decay).Unix()

	redisClient.Send("MULTI")
	redisClient.Send("ZINCRBY", "socialite:urls", 1, link.ID())
	redisClient.Send("ZADD", "socialite:expiry", expiryDate, link.ID())
	redisClient.Send("HSET", metadataKey(link.ID()), "url", link.URL)
	redisClient.Send("HSET", metadataKey(link.ID()), "title", link.Title)
	_, err = redisClient.Do("EXEC")

	if err != nil {
		log.WithFields(log.Fields{
			"id":      link.ID(),
			"package": "store",
		}).Error(err)
	}

	return
}

func Expire() (err error) {
	verifyConnection()

	expireAt := time.Now().Unix()
	expiredLinks, err := redis.Values(redisClient.Do(
		"ZRANGEBYSCORE",
		"socialite:expiry",
		"-inf",
		expireAt,
		"WITHSCORES"))

	if err != nil {
		log.WithFields(log.Fields{"package": "store"}).Error(err)
		return
	}

	for len(expiredLinks) > 0 {
		var (
			id     string
			expiry int64
		)

		expiredLinks, err = redis.Scan(expiredLinks, &id, &expiry)

		if err != nil {
			log.WithFields(log.Fields{
				"package": "store",
				"id":      id,
				"expiry":  expiry,
			}).Error(err)
			return
		}

		redisClient.Send("MULTI")
		redisClient.Send("ZREM", "socialite:urls", id)
		redisClient.Send("ZREM", "socialite:expiry", id)
		redisClient.Send("DEL", metadataKey(id))
		_, err = redisClient.Do("EXEC")

		if err != nil {
			log.WithFields(log.Fields{
				"package": "store",
				"id":      id,
				"expiry":  expiry,
			}).Error(err)
			return
		}

		log.WithFields(log.Fields{
			"package": "store",
			"id":      id,
			"expiry":  expiry,
		}).Info("Expired")
	}

	return
}

func Close() {
	if redisClient != nil {
		redisClient.Close()
	}
}

func metadataKey(id string) string {
	return fmt.Sprintf("socialite:urls:%s:data", id)
}

func verifyConnection() {
	if redisClient != nil {
		return
	}

	log.WithField("package", "store").Info("Connecting to Redis...")
	rc, err := redis.DialURL(os.Getenv("REDIS_URL"))

	if err != nil {
		// Should this be fatal?
		log.WithField("package", "store").Fatal(err)
	}

	redisClient = rc
}
