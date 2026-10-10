package memory

import (
	"notes/internal/domain"
	"sort"
	"strconv"
	"sync"
)

type Storage struct {
	mu    sync.RWMutex
	notes map[string]domain.Note
}

func New() *Storage {
	return &Storage{notes: make(map[string]domain.Note)}
}

func (s *Storage) Save(n domain.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[n.ID] = n
	return nil
}

func (s *Storage) Get(id string) (domain.Note, error) {
	s.mu.RLock()

	out := make([]domain.Note, 0, len(s.notes))
	for _, n := range s.notes {
		out = append(out, n)
	}

	sort.Slice(out, func(i, j int) bool {
		a, _ := strconv.Atoi(out[i].ID)
		b, _ := strconv.Atoi(out[j].ID)
		return a < b
	})
	return out, nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.notes, id)
	return nil
}
