package rest

import (
	"banking/db/pg"
	"fmt"
	"net/http"
)

type HealthcheckHandler struct {
	pgcli *pg.Client
}

func NewHealthcheckHandler(pc *pg.Client) *HealthcheckHandler {
	return &HealthcheckHandler{pc}
}

func (h *HealthcheckHandler) Run(w http.ResponseWriter, req *http.Request) {
	if err := h.pgcli.DB().Ping(); err != nil {
		http.Error(w, fmt.Sprintf("db ping: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}
