package main

import (
	"context"
	"design-pattern/configs"
	"design-pattern/internal/builder"
	"design-pattern/pkg/database"
	"design-pattern/pkg/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg, err := configs.NewConfig(".env")
	checkError(err)

	db, err := database.InitDatabase(cfg.PostgresConfig)
	checkError(err)

	publicRoutes := builder.BuildPublicRoutes(cfg, db)
	privateRoutes := builder.BuildPrivateRoutes(cfg, db)

	srv := server.NewServer(cfg, publicRoutes, privateRoutes)
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
