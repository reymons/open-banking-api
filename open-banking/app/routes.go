package app 

import (
    "net/http"
    "database/sql"
    "banking/open-banking/config"
    h "banking/open-banking/handlers"
)

func addRoutes(
    mux *http.ServeMux,
    conf *config.GlobalConfig,
    db *sql.DB,
) {
    _ = db
    hCustomer := h.NewCustomerHandler()

    mux.HandleFunc("GET /api/v1/customers/:id", hCustomer.GetCustomer)
}

