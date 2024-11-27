package main

import (
	"database/sql"
	"flag"
	"log"      // пакет лля работы с журналом
	"net/http" // пакет для работы с http
	"os"

	"dailynotes.com/snippetbox/pkg/models"
	"dailynotes.com/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)
type templateData struct{
	Snippet  *models.Snippet
    Snippets []*models.Snippet
}
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}
func main(){
	addr := flag.String("addr", ":8080", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "web:0000@/snippetbox?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}