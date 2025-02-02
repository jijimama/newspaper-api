package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"newspaper-api/entity"
)

type mockNewspaperRepository struct {
	mock.Mock // mock.Mockを埋め込むことでモックのメソッドを提供
}
// 新しいモックリポジトリを作成するファクトリーメソッド
func NewMockNewspaperRepository() *mockNewspaperRepository {
	return &mockNewspaperRepository{}
}
// モックリポジトリの Create メソッド
func (m *mockNewspaperRepository) Create(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	args := m.Called(newspaper) // Calledでモックされた挙動を呼び出す
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Newspaper), args.Error(1)
}

func (m *mockNewspaperRepository) Get(ID int) (*entity.Newspaper, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Newspaper), args.Error(1)
}
func (m *mockNewspaperRepository) Save(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	args := m.Called(newspaper)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Newspaper), args.Error(1)
}

func (m *mockNewspaperRepository) Delete(ID int) error {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

type NewspaperUseCaseSuite struct {
	suite.Suite
	newspaperUseCase *newspaperUseCase
}

func TestNewspaperUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(NewspaperUseCaseSuite))
}

func (suite *NewspaperUseCaseSuite) SetupSuite() {
}

func (suite *NewspaperUseCaseSuite) TestCreate() {
	title := "newspaper"
	columnName := "sports"
	mockNewspaperRepository := NewMockNewspaperRepository() // モックリポジトリを作成
	suite.newspaperUseCase = NewNewspaperUseCase(mockNewspaperRepository)// テスト対象の UseCase を作成

	newspaper := &entity.Newspaper{
		Title:      title,
		ColumnName: columnName,
	}
	// モックの挙動を設定
	mockNewspaperRepository.On("Create", newspaper).Return(&entity.Newspaper{
		ID:         1,
		Title:      title,
		ColumnName: columnName,
	}, nil)

	newspaper, err := suite.newspaperUseCase.Create(newspaper)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, newspaper.Title)
	suite.Assert().Equal(columnName, newspaper.ColumnName)
}

func (suite *NewspaperUseCaseSuite) TestGet() {
	title := "newspaper"
	columnName := "sports"
	mockNewspaperRepository := NewMockNewspaperRepository()
	suite.newspaperUseCase = NewNewspaperUseCase(mockNewspaperRepository)
	mockNewspaperRepository.On("Get", 1).Return(&entity.Newspaper{
		ID:         1,
		Title:      title,
		ColumnName: columnName,
	}, nil)

	newspaper, err := suite.newspaperUseCase.Get(1)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, newspaper.Title)
	suite.Assert().Equal(columnName, newspaper.ColumnName)
}

func (suite *NewspaperUseCaseSuite) TestUpdate() {
	title := "newspaper"
	columnName := "sports"
	newspaper := &entity.Newspaper{
		ID:         1,
		Title:      title,
		ColumnName: columnName,
	}

	mockNewspaperRepository := NewMockNewspaperRepository()
	suite.newspaperUseCase = NewNewspaperUseCase(mockNewspaperRepository)
	mockNewspaperRepository.On("Save", newspaper).Return(&entity.Newspaper{
		ID:         1,
		Title:      title,
		ColumnName: columnName,
	}, nil)

	newspaper, err := suite.newspaperUseCase.Save(newspaper)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, newspaper.Title)
	suite.Assert().Equal(columnName, newspaper.ColumnName)
}

func (suite *NewspaperUseCaseSuite) TestDelete() {
	mockNewspaperRepository := NewMockNewspaperRepository()
	suite.newspaperUseCase = NewNewspaperUseCase(mockNewspaperRepository)
	mockNewspaperRepository.On("Delete", 1).Return(nil, nil)
	err := suite.newspaperUseCase.Delete(1)
	suite.Assert().Nil(err)
}
