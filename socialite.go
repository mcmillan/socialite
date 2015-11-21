package main

import (
	"flag"

	"github.com/mcmillan/socialite/collector"
	"github.com/mcmillan/socialite/store"
	"github.com/mcmillan/socialite/web"

	"github.com/joho/godotenv"
)

func main() {
	defer store.Close()

	godotenv.Load()

	modePointer := flag.String("mode", "", "Either `web` or `collect`")
	flag.Parse()

	mode := string(*modePointer)

	if mode == "web" {
		web.Serve()
	} else if mode == "collect" {
		collector.Run()
	} else {
		flag.Usage()
	}
}
