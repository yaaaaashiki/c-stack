package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/yaaaaashiki/cstack/domain/model"
)

const Zero = 0

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (f *ItemRepository) FindAll() ([]model.Item, error) {
	items := []model.Item{}
	if err := f.db.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (f *ItemRepository) FindByID(id int) (*model.Item, error) {
	items := model.Item{}
	if err := f.db.Find(&items, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &items, nil
}

func (f *ItemRepository) FindByIDOrNil(id uint) (*model.Item, error) {
	item := model.Item{}
	res := f.db.Find(&item, "id=?", id)
	if res.RecordNotFound() {
		return nil, nil
	} else {
		if res.Error != nil {
			return nil, res.Error
		}
	}
	return &item, nil
}

func (f *ItemRepository) FindByEmailOrNil(email string) (*model.Item, error) {
	item := model.Item{}
	res := f.db.Find(&item, "email=?", email)
	if res.RecordNotFound() {
		return nil, nil
	} else {
		if res.Error != nil {
			return nil, res.Error
		}
	}
	return &item, nil
}

func (f *ItemRepository) FindAllByUserIDOrNil(userID string) ([]model.Item, error) {
	items := []model.Item{}
	res := f.db.Find(&items, "user_id=?", userID)
	if res.RecordNotFound() {
		return nil, nil
	} else {
		if res.Error != nil {
			return nil, res.Error
		}
	}
	return items, nil
}

//If first return value is true, input data is duplicate in items table
func (f *ItemRepository) IsExistItem(userID uint, name string) (bool, error) {
	item := model.Item{}
	res := f.db.Raw(`select * from items where user_id = ? and name = ?`, userID, name).Find(&item)
	if res.RecordNotFound() {
		return false, nil
	} else {
		if res.Error != nil {
			return true, res.Error
		}
	}
	return true, nil
}

func (f *ItemRepository) RegisterItem(userID uint, name string, price int, iconImage string, description string) (*model.Item, error) {
	newItem := model.Item{}
	newItem.UserID = userID
	newItem.Name = name
	newItem.Price = price
	newItem.CurrentPaymentPrice = Zero
	newItem.IconImage = iconImage
	newItem.Description = description

	if err := f.db.Create(&newItem).Error; err != nil {
		return nil, err
	}
	return &newItem, nil
}
