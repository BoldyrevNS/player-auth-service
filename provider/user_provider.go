package provider

import (
	"auth-ms/model"
	"gorm.io/gorm"
)

type UserProvider interface {
	Create(user model.User) error
	Delete(userId uint) error
	FindOneByEmail(email string) (model.User, error)
	FindAllExcludeCurrent(currentUserId uint) ([]model.User, error)
	FindOneById(id uint) (model.User, error)
}

type userProviderImpl struct {
	db *gorm.DB
}

func NewUserProvider(dbInstance *gorm.DB) UserProvider {
	return &userProviderImpl{
		db: dbInstance,
	}
}

func (p *userProviderImpl) Create(user model.User) error {
	res := p.db.Create(&user)
	return res.Error
}

func (p *userProviderImpl) Delete(userId uint) error {
	var user model.User
	res := p.db.Where("id = ?", userId).Delete(&user)
	return res.Error
}

func (p *userProviderImpl) FindOneByEmail(email string) (model.User, error) {
	var user model.User
	res := p.db.Where("email = ?", email).First(&user)
	return user, res.Error
}

func (p *userProviderImpl) FindOneById(id uint) (model.User, error) {
	var user model.User
	res := p.db.Where("id = ?", id).First(&user)
	return user, res.Error
}

func (p *userProviderImpl) FindAllExcludeCurrent(currentUserId uint) ([]model.User, error) {
	var users []model.User
	res := p.db.Where("id != ?", currentUserId).Find(&users)
	return users, res.Error
}
