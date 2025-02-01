package usecase

import (
	"newspaper-api/adapter/gateway"
	"newspaper-api/entity"
)
// ビジネスロジック「新聞に関する操作」の抽象化を提供する NewspaperUseCase インターフェース
type (
	NewspaperUseCase interface {
		Create(newspaper *entity.Newspaper) (*entity.Newspaper, error)
		Get(ID int) (*entity.Newspaper, error)
		Save(*entity.Newspaper) (*entity.Newspaper, error)
		Delete(ID int) error
	}
)
// データベース操作は gateway.NewspaperRepository に完全に委譲される(依存性の注入)
type newspaperUseCase struct {
	newspaperRepository gateway.NewspaperRepository
}
// 依存性注入の実現
func NewNewspaperUseCase(newspaperRepository gateway.NewspaperRepository) *newspaperUseCase {
	return &newspaperUseCase{
		newspaperRepository: newspaperRepository,
	}
}

func (a *newspaperUseCase) Create(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	return a.newspaperRepository.Create(newspaper)
}

func (a *newspaperUseCase) Get(ID int) (*entity.Newspaper, error) {
	return a.newspaperRepository.Get(ID)
}

func (a *newspaperUseCase) Save(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	return a.newspaperRepository.Save(newspaper)
}

func (a *newspaperUseCase) Delete(ID int) error {
	return a.newspaperRepository.Delete(ID)
}
