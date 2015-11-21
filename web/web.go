package web

import (
	"github.com/mcmillan/socialite/store"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

func Serve() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/", popularLinks)
	m.Run()
}

func popularLinks(r render.Render) {
	links, err := store.Popular()

	if err != nil {
		log.WithField("package", "web").Error(err)

		r.JSON(500, map[string]string{
			"error": err.Error(),
		})

		return
	}

	r.JSON(200, links)
}
