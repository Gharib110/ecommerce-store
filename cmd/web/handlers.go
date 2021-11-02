package main

import (
	"net/http"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "terminal", nil); err != nil {
		app.errLogger.Println(err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errLogger.Println(err)
		w.Header().Set("Content-Type", "application/text")
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	cardHolder := r.PostForm.Get("cardholder_name")
	email := r.PostForm.Get("email")
	paymentIntent := r.PostForm.Get("payment_intent")
	paymentMethod := r.PostForm.Get("payment_method")
	paymentAmount := r.PostForm.Get("payment_amount")
	paymentCurrency := r.PostForm.Get("payment_currency")

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency

	if err := app.renderTemplate(w, r, "succeeded", &TemplateData{
		Data: data,
	}); err != nil {
		app.errLogger.Println(err)
		w.Header().Set("Content-Type", "application/text")
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	return
}
