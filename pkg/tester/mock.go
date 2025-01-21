package tester

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"newspaper-api/pkg/logger"
)
// MockDB は、sqlmock と GORM を使用してモックデータベース接続を作成
func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	// sqlmock.New を使用してモックデータベース接続を作成
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		logger.Fatal(err.Error())
	}
    // GORM を使用してモックデータベース接続を開く
	mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "mock_db",
		DriverName:                "mysql",
		Conn:                      mockDB, // MySQLドライバが実際のデータベース接続の代わりにモック接続を使用
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mock, mockGormDB
}

type mockClock struct {
	t time.Time
}

func NewMockClock(t time.Time) mockClock {
	return mockClock{t}
}

func (m mockClock) Now() time.Time {
	return m.t
}
