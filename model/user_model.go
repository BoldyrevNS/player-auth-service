package model

type User struct {
	Id        uint   `gorm:"type:int;primary_key"`
	Firstname string `gorm:"type:varchar(255)"`
	Lastname  string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255);unique"`
	Password  string `gorm:"type:varchar(255)"`
	Role      string `gorm:"type:varchar(255)"`
}
