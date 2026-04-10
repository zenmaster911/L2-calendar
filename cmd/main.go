package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/znmaster911/L2-calendar/internal/config"
	"github.com/znmaster911/L2-calendar/internal/db"
	"github.com/znmaster911/L2-calendar/internal/logger"
	"github.com/znmaster911/L2-calendar/internal/server"
	"github.com/znmaster911/L2-calendar/pkg/handler"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
	"github.com/znmaster911/L2-calendar/pkg/services"
)

var wg sync.WaitGroup

func main() {
	cfg := config.MustLoad()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logg, err := logger.NewSLog(cfg.Log.Path, slog.LevelInfo)
	if err != nil {
		log.Fatalf("failed to start logger duer to %s", err)
	}
	slog.SetDefault(logg.Logger)

	DB, err := db.NewPostgresDb(*cfg.DB)
	if err != nil {
		for i := 0; err != nil && i < cfg.DB.MaxRetries; i++ {
			logg.ErrorCtx(ctx, "failed to start db:", err, "attempt", i)
			DB, err = db.NewPostgresDb(*cfg.DB)
		}
		if err != nil {
			logg.ErrorCtx(ctx, "failed to start db and attempts are over:", err)
			os.Exit(1)
		}
	}
	defer DB.DB.Close()
	DB.DbInit()

	repo := repositories.NewRepo(DB.DB)
	service := services.NewService(repo)
	handler := handler.NewHandler(service)

	srv := new(server.Server)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	shtdwnCtx, shtdwndelay := context.WithTimeout(ctx, 5*time.Second)
	defer shtdwndelay()

	go func() {
		<-sig
		logg.InfoCtx(ctx, "termination signal received")
		defer cancel()
		if err := srv.HttpsServer.Shutdown(shtdwnCtx); err != nil {
			logg.ErrorCtx(shtdwnCtx, "error occured in server shudown", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Run(cfg.App.Port, handler.InitRouter(), logg); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.ErrorCtx(ctx, "failed to start server", err)
			os.Exit(1)
		}
	}()
	logg.Info("server started sucessfully")
	wg.Wait()
	logg.Info("server stopped")

}
