package pkg

import (
	"os"
)

// 環境変数が設定されていない場合にデフォルト値を提供
func GetEnvDefault(key, defVal string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defVal
	}
	return val
}
