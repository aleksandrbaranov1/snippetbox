package main

import (
	"fmt" // пакет лля работы с журналом
	"html/template"
	"log"
	"net/http" // пакет для работы с http
	"strconv"
)

// обрабатывает http запрос (домашняя страница)
func home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files ...)

	if err != nil{
		log.Printf(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil{
		log.Printf(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Отображение заметки"))
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Отображение выбранной заметки с ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request){
	//w.Write([]byte("Форма для создания заметки"))

	if r.Method != http.MethodPost{
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "Метод запрещен", 405)
		return
	}
	w.Write([]byte("Создание новой заметки..."))
}

