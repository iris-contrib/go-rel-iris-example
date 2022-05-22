package handler

import (
	"errors"

	"github.com/kataras/iris/v12"
)

// ErrBadRequest error.
var ErrBadRequest = errors.New("Bad Request")

func render(c iris.Context, body interface{}, status int) {
	c.StatusCode(status)

	switch v := body.(type) {
	case string:
		c.JSON(struct {
			Message string `json:"message"`
		}{
			Message: v,
		})
	case error:
		c.JSON(struct {
			Error string `json:"error"`
		}{
			Error: v.Error(),
		})
	case nil:
		return
	default:
		c.JSON(body)
	}
}
