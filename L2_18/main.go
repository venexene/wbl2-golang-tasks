package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/venexene/calendar/handlers"
	"github.com/venexene/calendar/internal"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db := calendar.NewCalendar()

	router := gin.Default()
	log.Printf("Created GIN router")

	router.Use(handlers.CalendarMiddleware(db))
	router.Use(handlers.LoggingMiddleware())

	router.GET("/server_check", func(c *gin.Context) {
		handlers.TestServerHandle(c)
	})

	router.POST("/create_event", func(c *gin.Context) {
		handlers.AddHandle(c)
	})

	router.POST("/update_event", func(c *gin.Context) {
		handlers.UpdateHandle(c)
	})

	router.POST("/delete_event", func(c *gin.Context) {
		handlers.DeleteHandle(c)
	})

	router.GET("/events_for_day", func(c *gin.Context) {
		handlers.DayEventsHandle(c)
	})

	router.GET("/events_for_week", func(c *gin.Context) {
		handlers.WeekEventsHandle(c)
	})

	router.GET("/events_for_month", func(c *gin.Context) {
		handlers.MonthEventsHandle(c)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("Created server")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	log.Printf("Started HTTP server on port %s", port)

	<-ctx.Done()
	stop()
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Shutdown server")
}
