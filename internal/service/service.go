package service

import (
	"fmt"
	"notes/internal/domain"
	"strconv"
	"sync/atomic"
)

type Repository interface {
	Save(n domain.Note) error
	Get(id string) (domain.Note, error)
	List() ([]domain.Note, error)
	Delete(id string) error
}

type NoteService struct {
	repo   Repository
	lastID atomic.Int64
}

func New(repo Repository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(title, text string) (string, error) {
	n := domain.Note{Title: title, Text: text}

	if err := n.Validate(); err != nil {
		return "", err
	}

	n.ID = strconv.FormatInt(s.lastID.Add(1), 10)

	if err := s.repo.Save(n); err != nil {
		return "", nil
	}
	return n.ID, nil
}

func (s *NoteService) Get(id string) (domain.Note, error) {
	n, err := s.repo.Get(id)
	if err != nil {
		return domain.Note{}, fmt.Errorf("get note %q, %w", id, err)
	}
	return n, nil
}

func (s *NoteService) List() ([]domain.Note, error) {
	notes, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	return notes, nil
}

func (s *NoteService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete note %q: %w", id, err)
	}
	return nil
}
