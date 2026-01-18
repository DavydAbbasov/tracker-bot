package service

import (
	"tracker-bot/internal/repo"
)

type EntryService struct {
	repository repo.Repo
}

func New(repository repo.Repo) *EntryService {
	return &EntryService{
		repository: repository,
	}
}
