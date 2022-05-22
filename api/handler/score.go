package handler

import (
	"github.com/go-rel/rel"
	"github.com/iris-contrib/go-rel-iris-example/scores"
	"github.com/kataras/iris/v12"
)

// Score for score endpoints.
type Score struct {
	repository rel.Repository
}

// Index handle GET /
func (s Score) Index(c iris.Context) {
	var result scores.Score

	s.repository.Find(c, &result)
	render(c, result, 200)
}

// Points handle Get /points
func (s Score) Points(c iris.Context) {
	var result []scores.Point

	s.repository.FindAll(c, &result)
	render(c, result, 200)
}

// Mount handlers to router group.
func (s Score) Mount(router iris.Party) {
	router.Get("/", s.Index)
	router.Get("/points", s.Points)
}

// NewScore handler.
func NewScore(repository rel.Repository) Score {
	return Score{
		repository: repository,
	}
}
