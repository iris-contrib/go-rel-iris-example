package handler

import (
	"errors"
	"fmt"

	"github.com/iris-contrib/go-rel-iris-example/todos"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/kataras/iris/v12"
)

type ctx int

const (
	loadKey string = "todosLoadKey"
)

// Todos for todos endpoints.
type Todos struct {
	repository rel.Repository
	todos      todos.Service
}

// Index handle GET /.
func (t Todos) Index(c iris.Context) {
	var (
		result []todos.Todo
		filter = todos.Filter{
			Keyword: c.URLParam("keyword"),
		}
	)

	if str := c.URLParam("completed"); str != "" {
		completed := str == "true"
		filter.Completed = &completed
	}

	t.todos.Search(c, &result, filter)
	render(c, result, 200)
}

// Create handle POST /
func (t Todos) Create(c iris.Context) {
	var todo todos.Todo

	if err := c.ReadJSON(&todo); err != nil {
		render(c, ErrBadRequest, 400)
		return
	}

	if err := t.todos.Create(c, &todo); err != nil {
		render(c, err, 422)
		return
	}

	c.Header("Location", fmt.Sprint(c.Request().RequestURI, "/", todo.ID))
	render(c, todo, 201)
}

// Show handle GET /{ID}
func (t Todos) Show(c iris.Context) {
	todo := c.Values().Get(loadKey).(todos.Todo)

	render(c, todo, 200)
}

// Update handle PATCH /{ID}
func (t Todos) Update(c iris.Context) {
	var (
		todo    = c.Values().Get(loadKey).(todos.Todo)
		changes = rel.NewChangeset(&todo)
	)

	if err := c.ReadJSON(&todo); err != nil {
		render(c, ErrBadRequest, 400)
		return
	}

	if err := t.todos.Update(c, &todo, changes); err != nil {
		render(c, err, 422)
		return
	}

	render(c, todo, 200)
}

// Destroy handle DELETE /{ID}
func (t Todos) Destroy(c iris.Context) {
	todo := c.Values().Get(loadKey).(todos.Todo)

	t.todos.Delete(c, &todo)
	render(c, nil, 204)
}

// Clear handle DELETE /
func (t Todos) Clear(c iris.Context) {
	t.todos.Clear(c)
	render(c, nil, 204)
}

// Load is middleware that loads todos to context.
func (t Todos) Load(c iris.Context) {
	var (
		id, _ = c.Params().GetInt("ID")
		todo  todos.Todo
	)

	if err := t.repository.Find(c, &todo, where.Eq("id", id)); err != nil {
		if errors.Is(err, rel.ErrNotFound) {
			render(c, err, 404)
			c.StopExecution()
			return
		}
		panic(err)
	}

	c.Values().Set(loadKey, todo)
	c.Next()
}

// Mount handlers to router group.
func (t Todos) Mount(router iris.Party) {
	router.Get("/", t.Index)
	router.Post("/", t.Create)
	router.Get("/{ID:int}", t.Load, t.Show)
	router.Patch("/{ID:int}", t.Load, t.Update)
	router.Delete("/{ID:int}", t.Load, t.Destroy)
	router.Delete("/", t.Clear)
}

// NewTodos handler.
func NewTodos(repository rel.Repository, todos todos.Service) Todos {
	return Todos{
		repository: repository,
		todos:      todos,
	}
}
