package app

import (
    "fmt"
    "log"
    "net/http"
    "banking/db"
    "banking/open-banking/config"
)

func RunServer() {
    conf, err := config.NewGlobalConfig()
    if err != nil {
        log.Fatal(err)
    }

    db, err := dbase.NewDB(conf.DbUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    mux := http.NewServeMux()
    addRoutes(mux, conf, db)
    handler := http.Handler(mux)

    serv := http.Server{
        Addr: fmt.Sprintf("%s:%d", conf.ServerHost, conf.ServerPort),
        Handler: handler,
    }

    log.Printf("Open Banking API server is running on %s\n", serv.Addr)
    if err := serv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

// TODO: add graceful exit

