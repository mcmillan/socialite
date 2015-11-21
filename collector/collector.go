package collector

import (
	"github.com/mcmillan/socialite/metadata"
	"github.com/mcmillan/socialite/twitter"

	log "github.com/Sirupsen/logrus"
)

func Run() {
	var sinceID string
	statuses, err := twitter.GetListStatuses(sinceID)

	if err != nil {
		log.Panic(err)
	}

	for _, status := range statuses {
		for _, url := range status.Entities.URLs {
			md, err := metadata.ParseURL(url.ExpandedURL)

			if err != nil {
				log.Error(err)
				continue
			}

			log.Info(md)
		}
	}
}
