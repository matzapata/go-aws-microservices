package services

import (
	"github.com/matzapata/go-aws-microservices/services/names/repositories"
)

type NamesService struct {
	Repository repositories.NamesRepository
}

func NewNamesService(repo repositories.NamesRepository) *NamesService {
	return &NamesService{
		Repository: repo,
	}
}

func (s *NamesService) CreateName(name string) (string, error) {
	id, err := s.Repository.CreateName(name)
	return id, err
}
