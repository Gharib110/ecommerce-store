package main

import "net/http"

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	app.infoLogger.Println("Hit the teminal ...")
}
