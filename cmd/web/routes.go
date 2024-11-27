package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// создаем маршрутизатор
	mux := http.NewServeMux()
	// привязываем home к корневому url
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/snippet/show", app.showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}