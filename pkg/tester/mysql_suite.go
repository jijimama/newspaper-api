package tester

import (
	"context"
	"fmt"
	"newspaper-api/entity"
	"newspaper-api/infrastructure/database"
	"newspaper-api/pkg"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// Mysqlに接続するための構造体
type DBMySQLSuite struct {
	suite.Suite                             // testifyのSuite機能を埋め込む
	mySQLContainer testcontainers.Container // MySQLコンテナのインスタンス
	ctx            context.Context          // コンテナ操作用のコンテキスト
	DB             *gorm.DB				    // GORMのDBインスタンス
}

// SetupTestContainersは、MySQLコンテナをセットアップする関数
func (suite *DBMySQLSuite) SetupTestContainers() (err error) {
	configs := database.NewConfigMySQL()
	// 指定されたポートが使用可能になるまで待機
	pkg.WaitForPort(configs.Database, configs.Port, 10*time.Second)
	suite.ctx = context.Background() // コンテキストを初期化
	req := testcontainers.ContainerRequest{ // コンテナに対する設定のリクエストを作成
		Image: "mysql:8",
		Env: map[string]string{
			"MYSQL_DATABASE":             configs.Database,
			"MYSQL_USER":                 configs.User,
			"MYSQL_PASSWORD":             configs.Password,
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{fmt.Sprintf("%s:3306/tcp", configs.Port)}, // ポートマッピング
		WaitingFor:   wait.ForLog("port: 3306  MySQL Community Server"),  // MySQLの準備完了をログで確認
	}
	// リクエストをもとにコンテナを作成
	suite.mySQLContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true, // コンテナを開始する
	})

	suite.Assert().Nil(err)
	return nil
}

// SetupSuiteはテストスイート全体の初期設定を行う関数
func (suite *DBMySQLSuite) SetupSuite() {
	// MySQLコンテナのセットアップを実行
	err := suite.SetupTestContainers()
	suite.Assert().Nil(err)
	// データベース接続を作成
	db, err := database.NewDatabaseSQLFactory(database.InstanceMySQL)
	suite.Assert().Nil(err)
	suite.DB = db
	 // モデルのマイグレーションを実行
	for _, model := range entity.NewDomains() {
		err = suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

// TearDownSuiteはテストスイート全体のクリーンアップを行う関数
func (suite *DBMySQLSuite) TearDownSuite() {
	// MySQLコンテナが存在しない場合は何もしない
	if suite.mySQLContainer == nil {
		return
	}
	// MySQLコンテナを終了（停止と削除）する
	err := suite.mySQLContainer.Terminate(suite.ctx)
	suite.Assert().Nil(err)
}
