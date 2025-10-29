// This command builds an HTTP REST server
// Its API is aimed for the internal banking usage ONLY
package main

import (
	"banking/config"
	"banking/db/pg"
	"banking/handler/rest"
	"banking/middleware"
	"banking/middleware/auth"
	"banking/service"
	"banking/store"
	"banking/util"
	"context"
	"fmt"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// Seed random
	rand.Seed(time.Now().UnixNano())

	// Set up app config
	appCfg := config.NewInternalConfig()

	// Initialize AWS
	awsCfgUSEast1, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion("us-east-1"), // SES works only in this region
	)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	pgcli, err := pg.NewClient(appCfg.GetDbUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer pgcli.Close()

	// Prepare utils
	util.SetJwtSecret(appCfg.GetJwtSecret())

	// Initialize dependencies
	// AWS
	sesClient := ses.NewFromConfig(awsCfgUSEast1)
	// Stores
	clientStore := store.NewClient(pgcli)
	accountStore := store.NewAccount(pgcli)
	verifStore := store.NewVerification()
	emailVerifStore := store.NewEmailVerification(pgcli, verifStore)
	resetPswdReqStore := store.NewResetPasswordReq(pgcli)
	// Services
	permService := service.NewPerm()
	emailService := service.NewEmailService(sesClient, appCfg.GetEmailNoreply(), appCfg.GetAppName())
	authService := service.NewAuth(clientStore, emailVerifStore, emailService)
	accountService := service.NewAccount(permService, accountStore)
	passwordService := service.NewPassword(pgcli, emailService, resetPswdReqStore, clientStore, appCfg)
	// Handlers
	userHandler := rest.NewUserHandler(clientStore)
	authHandler := rest.NewAuthHandler(authService)
	accountHandler := rest.NewAccountHandler(accountService)
	healthcheckHandler := rest.NewHealthcheckHandler(pgcli)
	passwordHandler := rest.NewPasswordHandler(passwordService)

	// Initialize server
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("GET /api/healthcheck", healthcheckHandler.Run)
	mux.HandleFunc("POST /api/v1/auth/sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("DELETE /api/v1/auth/sign-out", authHandler.SignOut)
	mux.HandleFunc("POST /api/v1/auth/verification", authHandler.SubmitVerification)
	mux.HandleFunc("POST /api/v1/auth/verification/code", authHandler.SendVerificationCode)
	mux.HandleFunc("GET /api/v1/accounts", auth.Middleware(accountHandler.GetAll))
	mux.HandleFunc("POST /api/v1/accounts", auth.Middleware(accountHandler.Request))
	mux.HandleFunc("GET /api/v1/users/me", auth.Middleware(userHandler.GetMe))
	mux.HandleFunc("POST /api/v1/passwords/reset-requests", passwordHandler.RequestPasswordReset)
	mux.HandleFunc("POST /api/v1/passwords/reset-requests/submit", passwordHandler.ResetPassword)

	// Add middlewares
	h := http.Handler(mux)
	h = middleware.Logger(h)
	h = middleware.CORS(h, middleware.CORSConfig{
		Credentials: true,
		Origins:     appCfg.GetAllowedOrigins(),
		MaxAge:      300, // 5 min
		Headers:     []string{"Content-Type"},
		Methods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
		},
	})

	// Run server
	serv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", appCfg.GetServerHost(), appCfg.GetServerPort()),
		Handler: h,
	}

	log.Println("Listening on " + serv.Addr)
	if err := serv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
