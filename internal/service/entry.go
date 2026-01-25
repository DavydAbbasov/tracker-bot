package service

import (
	"tracker-bot/internal/repo"
)

type EntryService interface{}

type entryService struct {
	repo repo.EntryRepository
}

func NewEntryService(repo repo.EntryRepository) EntryService {
	return &entryService{
		repo: repo,
	}
}
