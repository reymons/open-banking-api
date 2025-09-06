// This command builds an HTTP REST server
// Its API is aimed for the internal banking usage ONLY
package main

import (
	"banking/db/pg"
	"banking/handler/rest"
	"banking/middleware/auth"
	"banking/service"
	"banking/store"
	"banking/util"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		dbUrl     = os.Getenv("DB_URL")
		servHost  = os.Getenv("SERVER_HOST")
		servPort  = os.Getenv("SERVER_PORT")
		jwtSecret = os.Getenv("JWT_SERCERT")
	)

	// Initialize database
	pgcli, err := pg.NewClient(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pgcli.Close()

	// Prepare utils
	util.SetJwtSecret(jwtSecret)

	// Initialize dependencies
	// Stores
	clientStore := store.NewClient(pgcli)
	// Services
	authService := service.NewAuth(clientStore)
	// Handlers
	authHandler := rest.NewAuthHandler(authService)
	clientHandler := rest.NewClientHandler()

	// Start server
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("POST /api/v1/auth/sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("GET /api/v1/clients", auth.Middleware(clientHandler.GetAll))

	// Run server
	serv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", servHost, servPort),
		Handler: mux,
	}

	log.Println("Listening on " + serv.Addr)
	if err := serv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
