// This command builds an HTTP REST server
// Its API is aimed for the internal banking usage ONLY
package main

import (
	"banking/db/pg"
	"banking/handler/rest"
	"banking/middleware"
	"banking/middleware/auth"
	"banking/service"
	"banking/store"
	"banking/util"
	"fmt"
	"log"
	"net/http"
	"os"
)

type body struct {
	ID int `json:id"`
}

func (body) Valid() map[string]string {
	return map[string]string{}
}

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
	accountStore := store.NewAccount(pgcli)
	// Services
	permService := service.NewPerm()
	authService := service.NewAuth(clientStore)
	accountService := service.NewAccount(permService, accountStore)
	// Handlers
	authHandler := rest.NewAuthHandler(authService)
	clientHandler := rest.NewClientHandler()
	accountHandler := rest.NewAccountHandler(accountService)

	// Initialize server
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("POST /api/v1/auth/sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("GET /api/v1/clients", auth.Middleware(clientHandler.GetAll))
	mux.HandleFunc("GET /api/v1/accounts", auth.Middleware(accountHandler.GetAll))
	mux.HandleFunc("POST /api/v1/accounts", auth.Middleware(accountHandler.Request))

	// Add middlewares
	h := http.Handler(mux)
	h = middleware.Logger(h)
	h = middleware.CORS(h, middleware.CORSConfig{
		Credentials: true,
		Origins: []string{
			"http://localhost:6969",
			"http://localhost:7171",
		},
		MaxAge: 300, // 5 min
	})

	// Run server
	serv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", servHost, servPort),
		Handler: h,
	}

	log.Println("Listening on " + serv.Addr)
	if err := serv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
