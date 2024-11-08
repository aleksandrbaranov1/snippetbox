package main

import (
	"log"      // пакет лля работы с журналом
	"net/http" // пакет для работы с http
)



func main(){
	// создаем маршрутизатор
	mux := http.NewServeMux()
	// привязываем home к корневому url
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/snippet/show", showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}