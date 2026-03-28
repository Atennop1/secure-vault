package storage

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Store(key string, value []byte) {
	s.repo.Store(key, value)
}

func (s *Service) Load(key string) ([]byte, bool) {
	return s.repo.Load(key)
}
