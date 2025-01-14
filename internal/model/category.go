package model

type Category struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	CategoryName string `gorm:"type:varchar(255);not null"`
	ParentId     uint   `gorm:"type:int;default:0"`
}
