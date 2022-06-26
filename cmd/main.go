package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/db_forum/internal/api"
)

func main() {
	// -------------------- Set up viper -------------------- //
	log := logrus.New()

	//formatter := logrus.JSONFormatter{
	//	TimestampFormat: time.RFC3339,
	//}
	log.SetLevel(logrus.InfoLevel)
	//log.SetFormatter(&formatter)

	// -------------------- Set up database -------------------- //

	dbPool, err := pgxpool.Connect(context.Background(), "postgres://root:password@localhost:5432/docker")
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer dbPool.Close()

	// -------------------- Set up service -------------------- //

	svc, err := api.NewAPIService(logrus.NewEntry(log), dbPool)
	if err != nil {
		log.Fatalf("error creating service instance: %s", err)
	}

	go svc.Serve()
	// -------------------- Listen for INT signal -------------------- //

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(5),
	)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
