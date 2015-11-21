package main

import (
	"github.com/mcmillan/socialite/collector"
	"github.com/mcmillan/socialite/store"

	log "github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
)

func main() {
	defer store.Close()

	godotenv.Load()
	collector.Run()

	p, err := store.Popular()

	if err != nil {
		log.Panic(err)
	}

	for _, l := range p {
		log.Info(l)
	}
}
