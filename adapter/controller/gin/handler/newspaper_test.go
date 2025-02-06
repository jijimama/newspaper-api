package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"newspaper-api/adapter/controller/gin/presenter"
	"newspaper-api/entity"
)

type MockNewspaperUseCase struct {
	mock.Mock
}

func NewMockNewspaperUseCase() *MockNewspaperUseCase {
	return &MockNewspaperUseCase{}
}

func (m *MockNewspaperUseCase) Create(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	args := m.Called(newspaper)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Newspaper), args.Error(1)
}

func (m *MockNewspaperUseCase) Get(ID int) (*entity.Newspaper, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Newspaper), args.Error(1)
}
func (m *MockNewspaperUseCase) Save(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	args := m.Called(newspaper)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Newspaper), args.Error(1)
}

func (m *MockNewspaperUseCase) Delete(ID int) error {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

type NewspaperHandlersSuite struct {
	suite.Suite
	newspaperHandler *NewspaperHandler
}

func TestNewspaperHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(NewspaperHandlersSuite))
}

func (suite *NewspaperHandlersSuite) TestCreate() {
	mockUseCase := NewMockNewspaperUseCase()
	newspaper := &entity.Newspaper{
		Title:      "newspaper",
		ColumnName: "sports",
	}

	mockUseCase.On("Create", newspaper).Return(&entity.Newspaper{
		ID:         1,
		Title:      "newspaper",
		ColumnName: "sports",
	}, nil)
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewCreateNewspaperRequest("/api/v1", presenter.CreateNewspaperJSONRequestBody{
		Title:       "newspaper",
		ColumnName:  "sports",
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	// ハンドラーを呼び出して、レスポンスを確認
	suite.newspaperHandler.CreateNewspaper(ginContext)

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse presenter.NewspaperResponse
	err := json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, w.Code)
	suite.Assert().Equal("newspaper", newspaperGetResponse.Title)
	suite.Assert().Equal("sports", newspaperGetResponse.ColumnName)
}

func (suite *NewspaperHandlersSuite) TestCreateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/api/v1/newspaper", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.newspaperHandler.CreateNewspaper(ginContext)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *NewspaperHandlersSuite) TestCreateFailure() {
	mockUseCase := NewMockNewspaperUseCase()
	newspaper := &entity.Newspaper{
		Title:      "newspaper",
		ColumnName: "sports",
	}

	mockUseCase.On("Create", newspaper).Return(nil,
		errors.New("invalid"))
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewCreateNewspaperRequest("/api/v1", presenter.CreateNewspaperJSONRequestBody{
		Title:      "newspaper",
		ColumnName: "sports",
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.CreateNewspaper(ginContext)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message":"invalid"}`, w.Body.String())
}

func (suite *NewspaperHandlersSuite) TestGet() {
	mockUseCase := NewMockNewspaperUseCase()
	mockUseCase.On("Get", 1).Return(&entity.Newspaper{
		ID:         1,
		Title:      "newspaper",
		ColumnName: "sports",
	}, nil)
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewGetNewspaperByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.GetNewspaperById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse presenter.NewspaperResponse
	err := json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("newspaper", newspaperGetResponse.Title)
	suite.Assert().Equal("sports", newspaperGetResponse.ColumnName)
}

func (suite *NewspaperHandlersSuite) TestGetNoNewspaperFailure() {
	mockUseCase := NewMockNewspaperUseCase()
	mockUseCase.On("Get", 1111).Return(nil,
		errors.New("invalid"))
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewGetNewspaperByIdRequest("/api/v1", 1111)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.GetNewspaperById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}

func (suite *NewspaperHandlersSuite) TestUpdate() {
	mockUseCase := NewMockNewspaperUseCase()
	title := "updated"
	columnName := "food"

	newspaper := &entity.Newspaper{
		ID:         1,
		Title:      title,
		ColumnName: columnName,
	}

	mockUseCase.On("Save", newspaper).Return(&entity.Newspaper{
		ID:         1,
		Title:      title,
		ColumnName: columnName,
	}, nil)

	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewUpdateNewspaperByIdRequest("/api/v1", 1,
		presenter.UpdateNewspaperByIdJSONRequestBody{
			Title:      &title,
			ColumnName: &columnName,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.UpdateNewspaperById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse presenter.NewspaperResponse
	err := json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("updated", newspaperGetResponse.Title)
	suite.Assert().Equal("food", newspaperGetResponse.ColumnName)
}

func (suite *NewspaperHandlersSuite) TestUpdateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PATCH", "/api/v1/newspaper", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.newspaperHandler.CreateNewspaper(ginContext)

	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *NewspaperHandlersSuite) TestUpdateNewspaperFailure() {
	mockUseCase := NewMockNewspaperUseCase()
	title := "updated"
	columnName := "food"
	newspaper := &entity.Newspaper{
		ID:         1111,
		Title:      title,
		ColumnName: columnName,
	}

	mockUseCase.On("Save", newspaper).Return(nil,
		errors.New("invalid"))
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewUpdateNewspaperByIdRequest("/api/v1", 1111,
		presenter.UpdateNewspaperByIdJSONRequestBody{
			Title:      &title,
			ColumnName: &columnName,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.UpdateNewspaperById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}

func (suite *NewspaperHandlersSuite) TestDelete() {
	mockUseCase := NewMockNewspaperUseCase()
	mockUseCase.On("Delete", 1).Return(nil, nil)
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewDeleteNewspaperByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.DeleteNewspaperById(ginContext, 1)

	suite.Assert().Equal(http.StatusNoContent, w.Code)
}

func (suite *NewspaperHandlersSuite) TestDeleteNewspaperFailure() {
	mockUseCase := NewMockNewspaperUseCase()
	mockUseCase.On("Delete", 1111).Return(nil, errors.New("invalid"))
	suite.newspaperHandler = NewNewspaperHandler(mockUseCase)

	request, _ := presenter.NewDeleteNewspaperByIdRequest("/api/v1", 1111)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.DeleteNewspaperById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}
