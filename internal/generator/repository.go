package generator

import "slices"

type Repository struct {
	storage []string
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Store(value string) {
	if r.Contains(value) {
		return
	}

	r.storage = append(r.storage, value)
}

func (r *Repository) Contains(value string) bool {
	return slices.Contains(r.storage, value)
}
