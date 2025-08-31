package handlers

import (
    "net/http"
)

type customerHandler struct {}

func (*customerHandler) GetCustomer(w http.ResponseWriter, req *http.Request) {}

func NewCustomerHandler() *customerHandler {
    return &customerHandler{}
}

