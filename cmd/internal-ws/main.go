package main

import (
	"banking/config"
	"banking/middleware"
	"banking/middleware/auth"
	"banking/util"
	"fmt"
	ws "golang.org/x/net/websocket"
	"log"
	"net/http"
)

func mainHandler(conn *ws.Conn) {
	defer conn.Close()

	user, err := auth.GetJwtUserFromReq(conn.Request())
	if err != nil {
		if err := ws.Message.Send(conn, "Unauthorized"); err != nil {
			log.Fatal(err)
		}
		return
	}

	_ = user

	if err := ws.Message.Send(conn, "Hello from WS!"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Init app config
	appCfg := config.NewInternalWsConfig()

	// Prepare utils
	util.SetJwtSecret(appCfg.GetJwtSecret())

	// Init main handler and middlewares
	wsHandler := ws.Handler(mainHandler)
	authHandler := auth.Middleware(func(w http.ResponseWriter, req *http.Request) {
		wsHandler.ServeHTTP(w, req)
	})
	h := http.Handler(http.HandlerFunc(authHandler))
	h = middleware.Logger(h)
	h = middleware.CORS(h, middleware.CORSConfig{
		Origins:     []string{"*"},
		Credentials: true,
	})

	// Init server
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", appCfg.GetServerHost(), appCfg.GetServerPort()),
		Handler: h,
	}

	log.Printf("Running WS server on %s\n", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
