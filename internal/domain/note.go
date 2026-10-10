package domain

import "errors"

var (
	ErrNotFound   = errors.New("заметка не найдена")
	ErrEmptyTitle = errors.New("заголовок не может быть пустым")
)

type Note struct {
	ID    string
	Title string
	Text  string
}

func (n Note) Validate() error {
	if n.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}
