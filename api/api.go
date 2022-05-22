package api

import (
	"github.com/iris-contrib/go-rel-iris-example/api/handler"
	"github.com/iris-contrib/go-rel-iris-example/scores"
	"github.com/iris-contrib/go-rel-iris-example/todos"

	"github.com/go-rel/rel"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
)

// New api.
func New(repository rel.Repository) *iris.Application {
	var (
		router         = iris.New()
		scores         = scores.New(repository)
		todos          = todos.New(repository, scores)
		healthzHandler = handler.NewHealthz()
		todosHandler   = handler.NewTodos(repository, todos)
		scoreHandler   = handler.NewScore(repository)
	)

	healthzHandler.Add("database", repository)

	router.UseRouter(recover.New())

	router.Use(cors.New().Handler())
	router.Use(withAccessLogger("requests.json"))
	router.Use(requestid.New())

	healthzHandler.Mount(router.Party("/healthz"))
	todosHandler.Mount(router.Party("/todos"))
	scoreHandler.Mount(router.Party("/score"))

	if err := router.Build(); err != nil {
		panic(err)
	}

	return router
}

func withAccessLogger(filename string) iris.Handler {
	// Initialize a new request access log middleware,
	// note that we use unbuffered data so we can have the results as fast as possible,
	// this has its cost use it only on debug.
	ac := accesslog.FileUnbuffered(filename)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = false
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})

	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac.Handler
}
