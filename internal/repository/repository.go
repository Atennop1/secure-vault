package repository

type Repository struct {
	// some redis here later

	storage map[string]string
}

func New() *Repository {
	return &Repository {
		storage: make(map[string]string),
	}
}

func (r *Repository) Store(key, value string) error {
	r.storage[key] = value
	return nil
}

func (r *Repository) Load(key string) (string, bool) {
	value, ok := r.storage[key]
	return value, ok
}
