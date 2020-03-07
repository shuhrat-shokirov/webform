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

func fromFLagOrEnv(flag *string, envName string) (server string, ok bool){
	if *flag != ""{
		return *flag, true
	}
	return os.LookupEnv(envName)
}

func main() {
	flag.Parse()
	hostf, ok := fromFLagOrEnv(host, envHost)
	if !ok {
		hostf = *host
	}
	portf, ok := fromFLagOrEnv(port, envPort)
	if !ok {
		portf = *port
	}
	dsnf, ok := fromFLagOrEnv(dsn, envDSN)
	if !ok {
		dsnf = *dsn
	}

	addr := net.JoinHostPort(hostf, portf)
	start(addr, dsnf)
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
