package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"webform/cmd/crud/app"
	"webform/pkg/crud/services/burgers"
)

var (
	hostF = flag.String("host", "0.0.0.0", "Server host")
	portF = flag.String("port", "9999", "Server port")
	dsnF  = flag.String("dsn", "postgres://jnsrsnhikdbaud:e82a7cf38ca97604587acd57f90bea13299f0a45b6a1a117e899c113b240da2e@ec2-35-172-85-250.compute-1.amazonaws.com:5432/d3dpla4h2nijpd", "Postgres DSN")
)

func main() {
	flag.Parse()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = *portF
	}
	addr := net.JoinHostPort(*hostF, port)
	start(addr, *dsnF)
}

func start(addr string, dsn string) {
	router := app.NewExactMux()
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	burgersSvc := burgers.NewBurgersSvc(pool)
	server := app.NewServer(
		router,
		pool,
		burgersSvc,
		filepath.Join("web", "templates"),
		filepath.Join("web", "assets"),
	)
	server.InitRoutes()
	panic(http.ListenAndServe(addr, server))
}
