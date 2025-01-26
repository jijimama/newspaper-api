package gateway_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"newspaper-api/adapter/gateway"
	"newspaper-api/entity"
	"newspaper-api/pkg/tester"
)

type NewspaperRepositorySuite struct {
	tester.DBSQLiteSuite // SQLiteに接続するための構造体(suite.Suitewを埋め込む)
	//repository フィールドがインターフェース型であるため、テスト時に本番用とモック用の実装を柔軟に切り替えられる
	repository gateway.NewspaperRepository // リポジトリ（インターフェース型）を格納するフィールド	
}
// NewspaperRepositorySuite を初期化して実行
func TestNewspaperRepositorySuite(t *testing.T) {
	// NewspaperRepositorySuite 内の各テストメソッドを自動的に呼び出す
	suite.Run(t, new(NewspaperRepositorySuite))
}
// SQLite データベースを初期化し、リポジトリのインスタンスを設定(実際のデータベース)
func (suite *NewspaperRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewNewspaperRepository(suite.DB)
}
// MockDB はテスト用にモックされたデータベースを作成し、リポジトリに設定
func (suite *NewspaperRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewNewspaperRepository(mockGormDB)
	return mock
}
// AfterTest は各テスト後にリポジトリをデフォルトの状態に戻す
func (suite *NewspaperRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewNewspaperRepository(suite.DB)
}
// NewspaperRepository の基本的な CRUD 操作をテスト
func (suite *NewspaperRepositorySuite) TestNewspaperRepositoryCRUD() {
	// 新しい新聞エンティティを作成
	newspaper := &entity.Newspaper{
		Title:      "test",
		ColumnName: "sports",
	}
	newspaper, err := suite.repository.Create(newspaper)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(newspaper.ID)
	suite.Assert().Equal("test", newspaper.Title)
	suite.Assert().Equal("sports", newspaper.ColumnName)

	getNewspaper, err := suite.repository.Get(newspaper.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", getNewspaper.Title)
	suite.Assert().Equal("sports", getNewspaper.ColumnName)

	getNewspaper.Title = "updated"
	updatedNewspaper, err := suite.repository.Save(getNewspaper)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedNewspaper.Title)
	suite.Assert().Equal("sports", updatedNewspaper.ColumnName)

	err = suite.repository.Delete(updatedNewspaper.ID)
	suite.Assert().Nil(err)
	deletedNewspaper, err := suite.repository.Get(updatedNewspaper.ID)
	suite.Assert().Nil(deletedNewspaper)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *NewspaperRepositorySuite) TestNewspaperCreateFailure() {
	mockDB := suite.MockDB()
    mockDB.ExpectBegin()
    mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `newspapers` (`title`,`column_name`) VALUES (?,?)")).WithArgs("test", "sports").WillReturnError(errors.New("create error"))
    mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	newspaper := &entity.Newspaper{
		Title:      "test",
		ColumnName: "sports",
	}

	createdNewspaper, err := suite.repository.Create(newspaper)
	suite.Assert().Nil(createdNewspaper)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *NewspaperRepositorySuite) TestNewspaperGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `newspapers` WHERE `newspapers`.`id` = ? ORDER BY `newspapers`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	newspaper, err := suite.repository.Get(1)
	suite.Assert().Nil(newspaper)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *NewspaperRepositorySuite) TestNewspaperDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `newspapers` WHERE id = ? AND `newspapers`.`id` = ?")).WithArgs(1, 1).WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	err := suite.repository.Delete(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}

func (suite *NewspaperRepositorySuite) TestNewspaperSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `newspapers` WHERE `newspapers`.`id` = ? ORDER BY `newspapers`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("save error"))

	newspaper := &entity.Newspaper{
		ID:         1,
		Title:      "test",
		ColumnName: "sports",
	}

	newspaper, err := suite.repository.Save(newspaper)
	suite.Assert().Nil(newspaper)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("save error", err.Error())
}
