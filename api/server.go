package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cloud-inventory/api/app"
	"cloud-inventory/utils"

	"github.com/labstack/echo/v4"
)

var ReadHeaderTimeout = 15 * time.Second

type Server struct {
	api  *echo.Echo
	errs chan error
}

func LaunchServer() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	SpecifyRoutes(e)
	server := &Server{api: e}
	// server.Start()
	// Start server
	go func() {
		port := ":5799"
		if err := server.api.Start(port); err != nil {
			e.Logger.Info("Shutting down api server!!!!!")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 4 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

func (s *Server) Start() {
	s.errs = make(chan error)
	go s.start(s.api, "cloud-inventory")
}

func (s *Server) start(e *echo.Echo, name string) {
	port := ":5799"
	fmt.Printf("  http server %s started on %q\n", name, port)
	s.errs <- e.Start(port)
}

func (s *Server) Wait() <-chan error {
	return s.errs
}

func (s *Server) Shutdown(ctx context.Context) error {
	g := utils.NewGroupShutdown(s.api)
	fmt.Print("  shutting down servers...")
	if err := g.Shutdown(ctx); err != nil {
		fmt.Println("failed: ", err.Error())
		return err
	}
	fmt.Println("ok.")
	return nil
}

func SpecifyRoutes(e *echo.Echo) {
	// Dummy ping pong route
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	// OVH
	g := e.Group("/ovh")
	g.GET("/server", app.GetAllServersHandler)
}
