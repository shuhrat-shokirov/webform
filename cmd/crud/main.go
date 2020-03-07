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
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)
const envHost = "HOST"
const envPort = "PORT"
const envDSN  = "DATABASE_URL"

func m(hostF string) (server string){
	server, ok := os.LookupEnv(hostF)
	if !ok {
		server = *host
	}
	return
}

func main() {
	flag.Parse()
	addr := net.JoinHostPort(m(envHost), m(envPort))
	start(addr, m(envDSN))
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
