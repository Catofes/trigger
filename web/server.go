package main

import (
	"context"
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	uuid "github.com/satori/go.uuid"
)

type server struct {
	verbose bool
	db      database
	e       echo.Echo
}

func (s *server) init(path string) {
	s.db.path = path
	s.db.load()
}

func (s *server) listen(port string) {
	e := echo.New()
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	e.PUT("/", s.put)
	e.GET("/:uuid", s.get)
	e.POST("/:uuid", s.post)
	e.GET("/:uuid/wait/:time", s.wait)

	e.Logger.Fatal(e.Start(":" + port))
}

func (s *server) put(ctx echo.Context) error {
	t := trigger{
		LastTriggerTime: time.Now(),
	}
	t.c, t.cancel = context.WithCancel(context.Background())
	u := uuid.NewV4().String()
	s.db.D[u] = t
	s.db.save()
	ctx.String(http.StatusOK, u)
	return nil
}

func (s *server) post(ctx echo.Context) error {
	u := ctx.Param("uuid")
	t, ok := s.db.D[u]
	if !ok {
		return echo.ErrNotFound
	}
	t.LastTriggerTime = time.Now()
	t.cancel()
	t.c, t.cancel = context.WithCancel(context.Background())
	s.db.D[u] = t
	s.db.save()
	ctx.NoContent(http.StatusOK)
	return nil
}

func (s *server) get(ctx echo.Context) error {
	u := ctx.Param("uuid")
	t, ok := s.db.D[u]
	if !ok {
		return echo.ErrNotFound
	}
	ctx.JSON(http.StatusOK, t)
	return nil
}

func (s *server) wait(ctx echo.Context) error {
	u := ctx.Param("uuid")
	t, ok := s.db.D[u]
	if !ok {
		return echo.ErrNotFound
	}
	uploadTime := ctx.Param("time")
	timestamp, err := strconv.ParseInt(uploadTime, 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}
	if timestamp < t.LastTriggerTime.Unix() {
		ctx.NoContent(http.StatusOK)
		return nil
	}
	select {
	case <-time.After(60 * time.Minute):
		ctx.NoContent(http.StatusNoContent)
		return nil
	case <-t.c.Done():
		ctx.NoContent(http.StatusOK)
		return nil
	}
}

func main() {
	path := flag.String("d", "./trigger.db", "Database Path.")
	port := flag.String("p", "1080", "Listen Port.")
	verbose := flag.Bool("verbose", false, "Verbose mode.")
	flag.Parse()
	s := &server{verbose: *verbose}
	s.init(*path)
	s.listen(*port)
}
