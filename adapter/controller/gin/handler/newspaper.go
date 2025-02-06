package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"newspaper-api/adapter/controller/gin/presenter"
	"newspaper-api/entity"
	"newspaper-api/pkg/logger"
	"newspaper-api/usecase"
)

type NewspaperHandler struct {
	newspaperUseCase usecase.NewspaperUseCase
}

func NewNewspaperHandler(newspaperUseCase usecase.NewspaperUseCase) *NewspaperHandler {
	return &NewspaperHandler{
		newspaperUseCase: newspaperUseCase,
	}
}
// Newspaperエンティティをレスポンス形式に変換する
func newspaperToResponse(newspaper *entity.Newspaper) *presenter.NewspaperResponse {
	return &presenter.NewspaperResponse{
		Id:         newspaper.ID,
		Title:      newspaper.Title,
		ColumnName: newspaper.ColumnName,
	}
}

func (a *NewspaperHandler) CreateNewspaper(c *gin.Context) {
	var requestBody presenter.CreateNewspaperJSONRequestBody
	// JSONボディをバインドして、リクエストが正しいかチェック
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	newspaper := &entity.Newspaper{
		Title:      requestBody.Title,
		ColumnName: requestBody.ColumnName,
	}

	createdNewspaper, err := a.newspaperUseCase.Create(newspaper)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newspaperToResponse(createdNewspaper))
}

func (a *NewspaperHandler) GetNewspaperById(c *gin.Context, ID int) {
	newspaper, err := a.newspaperUseCase.Get(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, newspaperToResponse(newspaper))
}

func (a *NewspaperHandler) UpdateNewspaperById(c *gin.Context, ID int) {
	var requestBody presenter.UpdateNewspaperByIdJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	newspaper := &entity.Newspaper{ID: ID, Title: *requestBody.Title, ColumnName: *requestBody.ColumnName}

	updatedNewspaper, err := a.newspaperUseCase.Save(newspaper)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, newspaperToResponse(updatedNewspaper))
}

func (a *NewspaperHandler) DeleteNewspaperById(c *gin.Context, ID int) {
	if err := a.newspaperUseCase.Delete(ID); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
