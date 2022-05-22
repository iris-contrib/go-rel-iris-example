package todos

import (
	"context"

	"github.com/iris-contrib/go-rel-iris-example/scores"

	"github.com/go-rel/rel"
)

type create struct {
	repository rel.Repository
	scores     scores.Service
}

func (c create) Create(ctx context.Context, todo *Todo) error {
	if err := todo.Validate(); err != nil {
		return err
	}

	// if completed, then earn a point.
	if todo.Completed {
		return c.repository.Transaction(ctx, func(ctx context.Context) error {
			c.repository.MustInsert(ctx, todo)
			return c.scores.Earn(ctx, "todo completed", 1)
		})
	}

	c.repository.MustInsert(ctx, todo)
	return nil
}
