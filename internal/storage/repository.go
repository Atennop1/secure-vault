package storage

type Repository struct {
	storage map[string][]byte

	// some redis here later
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[string][]byte),
	}
}

func (r *Repository) Store(key string, value []byte) {
	r.storage[key] = value
}

func (r *Repository) Load(key string) ([]byte, bool) {
	value, ok := r.storage[key]
	return value, ok
}
