package rest

import (
	"banking/core"
	"banking/middleware/auth"
	"banking/service"
	"banking/util"
	"net/http"
)

type AccountHandler struct {
	accountService service.Account
}

func NewAccountHandler(as service.Account) *AccountHandler {
	return &AccountHandler{as}
}

func (h *AccountHandler) GetAll(w http.ResponseWriter, req *http.Request) {
	user, ok := auth.GetJwtUser(w, req)
	if !ok {
		return
	}

	accounts, err := h.accountService.GetAllByUserID(
		req.Context(),
		user.ID,
		user.Role,
	)
	if err != nil {
		sendHttpError(w, req, "get account by user id", err)
		return
	}

	result := make([]accRes, 0)
	for _, acc := range accounts {
		result = append(result, accRes{
			ID:       acc.ID,
			Number:   acc.Number,
			Balance:  acc.Balance,
			Currency: acc.CurrencyCode,
			Status:   acc.Status,
		})
	}

	util.EncodeBody(w, http.StatusOK, &getAllAccsRes{result})
}

func (h *AccountHandler) Request(w http.ResponseWriter, req *http.Request) {
	user, ok := auth.GetJwtUser(w, req)
	if !ok {
		return
	}

	body, ok := util.DecodeBody[requestAccReq](w, req)
	if !ok {
		return
	}

	currencyID, err := core.GetCurrencyID(body.Currency)
	if err != nil {
		sendHttpError(w, req, "get currency id", err)
		return
	}

	err = h.accountService.Request(
		req.Context(),
		currencyID,
		user.ID,
		user.Role,
	)
	if err != nil {
		sendHttpError(w, req, "request account", err)
	}
}
