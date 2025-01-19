package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"newspaper-api/entity"
)

// Newspaper構造体のテスト
func TestNewspaper(t *testing.T) {
	newspaper := entity.Newspaper{
		ID:          1,
		Title:       "newspaper",
		ColumnName:  "columnName",
	}
	assert.Equal(t, 1, newspaper.ID)
	assert.Equal(t, "newspaper", newspaper.Title)
	assert.Equal(t, "columnName", newspaper.ColumnName)
}
