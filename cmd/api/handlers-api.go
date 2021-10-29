package main

import (
	"encoding/json"
	"net/http"
)

type stripePaylod struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Content string `json:"content"`
	ID      string `json:"id"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	rsp := jsonResponse{
		OK:      true,
		Message: "",
		Content: "",
		ID:      "",
	}
	out, err := json.MarshalIndent(rsp, "", "	")
	if err != nil {
		app.errLogger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	return
}
