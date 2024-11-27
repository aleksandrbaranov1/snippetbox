package main

import (
	"errors"
	"fmt" // пакет лля работы с журналом
	"html/template"

	//"html/template"
	//"log"
	"net/http" // пакет для работы с http
	"strconv"

	"dailynotes.com/snippetbox/pkg/models"
)

// обрабатывает http запрос (домашняя страница)
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }
 
    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }
 
    // Создаем экземпляр структуры templateData,
    // содержащий срез с заметками.
    data := &templateData{Snippets: s}
 
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }
 
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }
 
    // Передаем структуру templateData в шаблонизатор.
    // Теперь она будет доступна внутри файлов шаблона через точку.
    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Страница не найдена.
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	

	fiels := []string{
		"ui/html/show.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(fiels...)
	if err != nil{
		app.serverError(w, err);
	}

	err = ts.Execute(w, s)
	if err != nil{
		app.serverError(w, err)
	}
	// Отображаем весь вывод на странице.
	//fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }
 
	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"
 
	// Передаем данные в метод SnippetModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		//app.serverError(w, err)
		return
	}
 
	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

