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
	host = flag.String("host", "0.0.0.0", "Server host")
	port = flag.String("port", "9999", "Server port")
	dsn  = flag.String("dsn", "postgres://app:pass@localhost:5432/app", "Postgres DSN")
)
const envHost = "HOST"
const envPort = "PORT"
const envDSN  = "DATABASE_URL"


func main() {
	flag.Parse()
	serverHost, ok := os.LookupEnv(envHost)
	if !ok {
		serverHost = *host
	}
	serverPort, ok := os.LookupEnv(envPort)
	if !ok {
		serverPort = *port
	}
	serverDsn, ok := os.LookupEnv(envDSN)
	if !ok {
		serverDsn = *dsn
	}
	addr := net.JoinHostPort(serverHost, serverPort)
	start(addr, serverDsn)
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
