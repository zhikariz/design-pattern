package main

import (
	"context"
	"design-pattern/configs"
	"design-pattern/pkg/database"
	"design-pattern/pkg/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

type User struct {
	ID   int64
	Name string
}

func main() {
	cfg, err := configs.NewConfig(".env")
	checkError(err)

	_, err = database.InitDatabase(cfg.PostgresConfig)
	checkError(err)

	// init builder dan butuh db
	// init routes

	srv := server.NewServer(cfg)
	runServer(srv, cfg.PORT)
	waitForShutdown(srv)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func runServer(srv *server.Server, port string) {
	go func() {
		err := srv.Start(fmt.Sprintf(":%s", port))
		log.Fatal(err)
	}()
}

func waitForShutdown(srv *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			srv.Logger.Fatal(err)
		}
	}()
}
