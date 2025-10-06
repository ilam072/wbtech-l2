package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/config"
	eventpostgres "github.com/ilam072/wbtech-l2/18-calendar/internal/event/repo/postgres"
	eventrest "github.com/ilam072/wbtech-l2/18-calendar/internal/event/rest"
	eventsrv "github.com/ilam072/wbtech-l2/18-calendar/internal/event/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/middleware"
	usrpostgres "github.com/ilam072/wbtech-l2/18-calendar/internal/user/repo/postgres"
	userrest "github.com/ilam072/wbtech-l2/18-calendar/internal/user/rest"
	usersrv "github.com/ilam072/wbtech-l2/18-calendar/internal/user/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/validator"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/db"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/jwt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fiber.New()
	cfg := config.New()
	l := newLogger()

	DB, err := db.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	app.Use(middleware.Logger(l))

	userRepo := usrpostgres.NewUserRepo(DB)
	manager := jwt.NewManager([]byte(cfg.JWTConfig.SecretKey))
	v := validator.New()
	user := usersrv.NewUser(userRepo, manager, cfg.JWTConfig.TokenTTL)

	userHandler := userrest.NewUserHandler(l, user, v)
	auth := app.Group("/api/auth")
	auth.Post("/sign-up", userHandler.SignUp)
	auth.Post("/sign-in", userHandler.SignIn)

	eventRepo := eventpostgres.NewEventRepo(DB)
	event := eventsrv.NewEvent(eventRepo)
	eventHandler := eventrest.NewEventHandler(l, event, v)

	events := app.Group("/api", middleware.Auth(manager))
	events.Post("/update_event", eventHandler.UpdateEventHandler)
	events.Post("/create_event", eventHandler.CreateEventHandler)
	events.Delete("/delete_event/:id", eventHandler.DeleteEventHandler)
	events.Get("/events_for_day", eventHandler.GetEventsForDayHandler)
	events.Get("/events_for_week", eventHandler.GetEventsForWeekHandler)
	events.Get("/events_for_month", eventHandler.GetEventsForMonthHandler)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(cfg.ServerConfig.Address()); err != nil {
			l.Error("failed to listen server", "err", err)
		}
	}()

	<-sigs
	l.Info("shutting down...")

	if err := app.Shutdown(); err != nil {
		l.Error("failed to shutdown server", "err", err)
	}
}

func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return l
}
