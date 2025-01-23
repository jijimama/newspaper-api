package tester

import (
	"fmt"
	"newspaper-api/entity"
	"newspaper-api/infrastructure/database"
	"os"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)
// SQLiteに接続するための構造体
type DBSQLiteSuite struct {
	suite.Suite     // testifyのSuite機能を埋め込む
	DB     *gorm.DB // GORMのDBインスタンス
	DBName string   // テスト用のSQLiteデータベースファイル名
}

// テスト前に自動で実行されるメソッド
func (suite *DBSQLiteSuite) SetupSuite() {
	// テスト用のSQLiteデータベースファイル名を設定
	suite.DBName = fmt.Sprintf("%s.unittest.sqlite", suite.T().Name())
	// 環境変数にデータベース名を設定
	os.Setenv("DB_NAME", suite.DBName)
	// SQLiteデータベース接続を作成
	db, err := database.NewDatabaseSQLFactory(database.InstanceSQLite)
	suite.Assert().Nil(err)
	suite.DB = db
	// モデルのマイグレーションを実行
	for _, model := range entity.NewDomains() {
		err := suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

// テスト後に実行されるメソッド
func (suite *DBSQLiteSuite) TearDownSuite() {
	// テスト用のSQLiteデータベースファイルを削除
	err := os.Remove(suite.DBName)
	suite.Assert().Nil(err)
	// 環境変数からデータベース名を削除
	os.Unsetenv(suite.DBName)
}
