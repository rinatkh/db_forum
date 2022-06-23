package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"github.com/rinatkh/db_forum/internal/api"
)

const (
	configPathEnvVar = "CONFIG_PATH"
	defaultAddress   = "0.0.0.0"
	defaultPort      = "8080"
)

func main() {
	// -------------------- Set up viper -------------------- //

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(configPathEnvVar))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("fatal error config file: %s \n", err)
		os.Exit(1)
	}

	viper.SetDefault("service.bind.address", defaultAddress)
	viper.SetDefault("service.bind.port", defaultPort)

	// -------------------- Set up logging -------------------- //

	log := logrus.New()

	formatter := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	var debug bool
	switch viper.GetString("logging.level") {
	case "warning":
		log.SetLevel(logrus.WarnLevel)
	case "notice":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		debug = true
		formatter.PrettyPrint = true
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	log.SetFormatter(&formatter)

	log.Infof("log level: %s", log.Level.String())

	// -------------------- Set up database -------------------- //

	dbPool, err := pgxpool.Connect(context.Background(), viper.GetString("db.connection_string"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer dbPool.Close()

	// -------------------- Set up service -------------------- //

	svc, err := api.NewAPIService(logrus.NewEntry(log), dbPool, debug)
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
		time.Second*time.Duration(viper.GetInt("service.shutdown_timeout")),
	)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
