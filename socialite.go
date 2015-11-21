package main

import (
	"github.com/mcmillan/socialite/collector"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	collector.Run()
}
