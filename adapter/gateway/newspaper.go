package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"newspaper-api/entity"
)
// インターフェイス NewspaperRepository は、新聞に関するデータベース操作を抽象化（内側の層からアクセス）
type NewspaperRepository interface {
	Create(newspaper *entity.Newspaper) (*entity.Newspaper, error)
	Get(ID int) (*entity.Newspaper, error)
	Save(*entity.Newspaper) (*entity.Newspaper, error)
	Delete(ID int) error
}

type newspaperRepository struct {
	db *gorm.DB // GORMを使用したデータベース接続
}
// newspaperRepository のインスタンスを作成するファクトリーメソッド
func NewNewspaperRepository(db *gorm.DB) NewspaperRepository {
	return &newspaperRepository{db: db}
}

func (a *newspaperRepository) Create(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	// newspaperは既にポインタ型（*entity.Newspaper）で渡されている
	if err := a.db.Create(newspaper).Error; err != nil {
		return nil, err
	}
	return newspaper, nil
}

func (a *newspaperRepository) Get(ID int) (*entity.Newspaper, error) {
	var newspaper = entity.Newspaper{}
	// First メソッドは、データをポインタ経由で書き込むため、&newspaperを渡す必要がある
	if err := a.db.First(&newspaper, ID).Error; err != nil {
		return nil, err
	}
	return &newspaper, nil
}

func (a *newspaperRepository) Save(newspaper *entity.Newspaper) (*entity.Newspaper, error) {
	selectedNewspaper, err := a.Get(newspaper.ID)
	if err != nil {
		return nil, err
	}

	// copier を使用して、newspaper のフィールドを selectedNewspaper にコピー
	// コピー時に空のフィールドを無視し、深いコピーを実行
	if err := copier.CopyWithOption(selectedNewspaper, newspaper, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := a.db.Save(selectedNewspaper).Error; err != nil {
		return nil, err
	}

	return selectedNewspaper, nil
}

func (a *newspaperRepository) Delete(ID int) error {
	// newspaper.ID は値型（int）で渡されているため、ポインタ型（&newspaper.ID）に変換する必要がある
	newspaper := entity.Newspaper{ID: ID}
	if err := a.db.Where("id = ?", &newspaper.ID).Delete(&newspaper).Error; err != nil {
		return err
	}
	return nil
}
