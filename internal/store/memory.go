package store

import "errors"

type MemoryStore struct {
	db map[string]Link
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		db: make(map[string]Link),
	}
}

func (s *MemoryStore) Save(l Link) {
	s.db[l.ShortCode] = l
}

func (s *MemoryStore) Get(code string) (Link, error) {
	link, ok := s.db[code]
	if !ok {
		return Link{}, errors.New("Link not found")
	}

	return link, nil
}
