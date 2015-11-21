package web

import (
	"github.com/mcmillan/socialite/store"

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
		r.JSON(500, err)
		return
	}

	r.JSON(200, links)
}
