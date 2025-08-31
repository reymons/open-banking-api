package config

import (
    "os"
    "fmt"
    "strconv"
    "net/url"
)

type GlobalConfig struct {
    DbUrl string
    ServerHost string
    ServerPort uint16
}

func NewGlobalConfig() (*GlobalConfig, error) {
    var conf GlobalConfig 
    
    var (
        dbUrl = os.Getenv("DB_URL")
        servHost = os.Getenv("SERVER_HOST")
        servPort = os.Getenv("SERVER_PORT")
    )

    if _, err := url.ParseRequestURI(dbUrl); err != nil {
        return nil, fmt.Errorf("parse db url: %w", err)
    }
    conf.DbUrl = dbUrl

    if servHost == "" {
        return nil, fmt.Errorf("invalid server host %s", servHost)
    }
    conf.ServerHost = servHost

    port, err := strconv.Atoi(servPort)
    if err != nil {
        return nil, fmt.Errorf("convert server port to integer: %w", err)
    }
    if port < 1 || port > 65536 {
        return nil, fmt.Errorf("server port should be between 1 and 65536")
    }
    conf.ServerPort = uint16(port)

    return &conf, nil
}

