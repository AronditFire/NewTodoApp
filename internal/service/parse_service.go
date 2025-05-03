package service

import (
	"errors"
	"strings"

	"github.com/AronditFire/todo-app/entity"
	"github.com/AronditFire/todo-app/internal/repository"
)

type ParseService struct {
	repo repository.ParsingJSON
}

func NewParseService(repo repository.ParsingJSON) *ParseService {
	return &ParseService{repo: repo}
}

func (s *ParseService) ParseJSON(bindfile entity.BindFile) error {
	if !strings.HasSuffix(strings.ToLower(bindfile.File.Filename), ".json") { // check .json
		return errors.New("file must be .json format")
	}

	return s.repo.ParseJSON(bindfile)
}

func (s *ParseService) GetJsonTable() ([]map[string]any, error) {
	return nil, nil
}
