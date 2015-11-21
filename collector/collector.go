package collector

import (
	"github.com/mcmillan/socialite/metadata"
	"github.com/mcmillan/socialite/store"
	"github.com/mcmillan/socialite/twitter"

	log "github.com/Sirupsen/logrus"
)

func Run() {
	sinceID, err := store.SinceID()

	if err != nil {
		log.Panic(err)
	}

	store.Expire()

	statuses, err := twitter.GetListStatuses(sinceID)

	if err != nil {
		log.Panic(err)
	}

	for index, status := range statuses {
		if index == 0 {
			store.SetSinceID(status.ID)
		}

		for _, url := range status.Entities.URLs {
			link, err := metadata.ParseURL(url.ExpandedURL)

			if err != nil {
				log.Error(err)
				continue
			}

			store.Add(link)
		}
	}
}
