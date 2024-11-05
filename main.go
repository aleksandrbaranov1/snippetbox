package main

import (
	"fmt"
	"log"      // пакет лля работы с журналом
	"net/http" // пакет для работы с http
	"strconv"
)

// обрабатывает http запрос (домашняя страница)
func home(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Привет"))
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request){
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


func main(){
	// создаем маршрутизатор
	mux := http.NewServeMux()
	// привязываем home к корневому url
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/snippet/show", showSnippet)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}