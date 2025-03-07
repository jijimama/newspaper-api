package action

import (
	"newspaper-api/entity"
	"newspaper-api/pkg/logger"
	"newspaper-api/usecase"
)

var NewspaperName string

type NewspaperAction struct {
	newspaperUseCase usecase.NewspaperUseCase
}

func NewNewspaperAction(newspaperUseCase usecase.NewspaperUseCase) *NewspaperAction {
	return &NewspaperAction{
		newspaperUseCase: newspaperUseCase,
	}
}

func (a *NewspaperAction) CreateNewspaper(title string) (*entity.Newspaper, error) {
	newspaper := &entity.Newspaper{
		Title:       title,
	}

	createdNewspaper, err := a.newspaperUseCase.Create(newspaper)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return createdNewspaper, nil
}
