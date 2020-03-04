package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"net/http"
	"path/filepath"
	"webform/cmd/crud/app"
	"webform/pkg/crud/services/burgers"
)

var (
	host = flag.String("host", "0.0.0.0", "Server host")
	port = flag.String("port", "9999", "Server port")
	dsn  = flag.String("dsn", "postgres://user:pass@localhost:5432/app", "Postgres DSN")
)

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)
	start(addr, *dsn)
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
