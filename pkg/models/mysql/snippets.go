package mysql

import (
	"database/sql"
	"errors"

	"dailynotes.com/snippetbox/pkg/models"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB *sql.DB
}
//Метод получения заметки по ее id
func (m * SnippetModel) Get(id int) (*models.Snippet, error){
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, models.ErrNoRecord
		}else{
			return nil, err
		}
	}
	return s, nil
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
 
	return int(id), nil
}
// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		// Создаем указатель на новую структуру Snippet
		s := &models.Snippet{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Опять же, аргументы предоставленные в row.Scan()
		// должны быть указателями на место, куда требуется скопировать данные и
		// количество аргументов должно быть точно таким же, как количество
		// столбцов из таблицы базы данных, возвращаемых вашим SQL запросом.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		snippets = append(snippets, s)
	}
 
	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}
 
	// Если все в порядке, возвращаем срез с данными.
	return snippets, nil
}