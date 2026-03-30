package generator

import "slices"

// TODO: these slugs aren't in sync with storage ones, will fix when add Redis

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
