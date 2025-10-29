package config

import (
	"banking/util"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type InternalConfig struct {
	appName        string
	srvHost        string
	srvPort        string
	jwtSecret      string
	allowedOrigins []string
	dbUrl          string
	mainWebsiteUrl string
	emailNoreply   string
}

func NewInternalConfig() *InternalConfig {
	cfg := &InternalConfig{}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		fatalEmpty("App name", "APP_NAME")
	}
	cfg.appName = appName

	srvHost := os.Getenv("SERVER_HOST")
	if srvHost == "" {
		fatalEmpty("Server host", "SERVER_HOST")
	}
	cfg.srvPort = srvHost

	srvPort := os.Getenv("SERVER_PORT")
	if !util.IsValidPortStr(srvPort) {
		fatalInvalidPort("Server port", "SERVER_PORT")
	}
	cfg.srvPort = srvPort

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fatalEmpty("Jwt secret", "JWT_SECRET")
	}
	cfg.jwtSecret = jwtSecret

	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		fatalEmpty("Allowed origins", "ALLOWED_ORIGINS")
	}
	cfg.allowedOrigins = strings.Split(allowedOriginsStr, ",")

	mainWebsiteUrl := os.Getenv("MAIN_WEBSITE_URL")
	if mainWebsiteUrl == "" {
		fatalEmpty("Main website url", "MAIN_WEBSITE_URL")
	}
	cfg.mainWebsiteUrl = mainWebsiteUrl

	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbPswd := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	if dbUser == "" {
		fatalEmpty("Database user", "DB_USER")
	}
	if dbPswd == "" {
		fatalEmpty("Database password", "DB_PASSWORD")
	}
	if dbHost == "" {
		fatalEmpty("Database host", "DB_HOST")
	}
	if dbName == "" {
		fatalEmpty("Database name", "DB_NAME")
	}
	if !util.IsValidPortStr(dbPort) {
		fatalInvalidPort("Database port", "DB_PORT")
	}
	cfg.dbUrl = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPswd, dbHost, dbPort, dbName,
	)

	emailNoreply := os.Getenv("EMAIL_NOREPLY")
	if !util.IsValidEmail(emailNoreply) {
		fatalInvalidEmail("Noreply", "EMAIL_NOREPLY")
	}
	cfg.emailNoreply = emailNoreply

	return cfg
}

func (cfg *InternalConfig) GetAppName() string {
	return cfg.appName
}

func (cfg *InternalConfig) GetServerHost() string {
	return cfg.srvHost
}

func (cfg *InternalConfig) GetServerPort() string {
	return cfg.srvPort
}

func (cfg *InternalConfig) GetJwtSecret() string {
	return cfg.jwtSecret
}

func (cfg *InternalConfig) GetAllowedOrigins() []string {
	return cfg.allowedOrigins
}

func (cfg *InternalConfig) GetDbUrl() string {
	return cfg.dbUrl
}

func (cfg *InternalConfig) GetMainWebsiteUrl() string {
	return cfg.mainWebsiteUrl
}

func (cfg *InternalConfig) GetEmailNoreply() string {
	return cfg.emailNoreply
}
