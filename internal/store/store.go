package store

type Store interface {
	Save(l Link) error
	Get(code string) (Link, error)
	IncrementClick(code string) error
}
