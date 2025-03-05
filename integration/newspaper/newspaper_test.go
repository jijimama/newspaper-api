package integration

import (
	"context"
	"net/http"
	"testing"
	"github.com/stretchr/testify/suite"

	"newspaper-api/adapter/controller/gin/presenter"
	"newspaper-api/pkg"
)

type NewspaperTestSuite struct {
	suite.Suite
}

func TestNewspaperSuite(t *testing.T) {
	suite.Run(t, new(NewspaperTestSuite))
}

func (suite *NewspaperTestSuite) TestNewspaperCreateGetDelete() {
	// Create
	baseEndpoint := pkg.GetEndpoint("/api/v1")
	apiClient, _ := presenter.NewClientWithResponses(baseEndpoint)
	createResponse, err := apiClient.CreateNewspaperWithResponse(context.Background(), presenter.CreateNewspaperJSONRequestBody{
		Title:      "test",
		ColumnName: "sports",
	})
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, createResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().NotNil(createResponse.JSON201.Id)
	suite.Assert().Equal("test", createResponse.JSON201.Title)
	suite.Assert().Equal("sports", createResponse.JSON201.ColumnName)

	// Get
	getResponse, err := apiClient.GetNewspaperByIdWithResponse(context.Background(), createResponse.JSON201.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, getResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().Equal(createResponse.JSON201.Id, getResponse.JSON200.Id)
	suite.Assert().Equal("test", getResponse.JSON200.Title)
	suite.Assert().Equal("sports", getResponse.JSON200.ColumnName)

	// Update
	title := "updated"
	columnName := "food"
	updateResponse, err := apiClient.UpdateNewspaperByIdWithResponse(context.Background(), getResponse.JSON200.Id, presenter.UpdateNewspaperByIdJSONRequestBody{
		Title:      &title,
		ColumnName: &columnName,
	})
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, updateResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updateResponse.JSON200.Title)
	suite.Assert().Equal("food", updateResponse.JSON200.ColumnName)

	// Delete
	deleteResponse, err := apiClient.DeleteNewspaperByIdWithResponse(context.Background(), updateResponse.JSON200.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusNoContent, deleteResponse.StatusCode())
}
