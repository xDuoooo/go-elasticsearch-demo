package model

import "gorm.io/gorm"

// Hotel 酒店实体，对应数据库表 tb_hotel
type Hotel struct {
	ID        int64  `gorm:"column:id;primaryKey" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	Address   string `gorm:"column:address" json:"address"`
	Price     int    `gorm:"column:price" json:"price"`
	Score     int    `gorm:"column:score" json:"score"`
	Brand     string `gorm:"column:brand" json:"brand"`
	City      string `gorm:"column:city" json:"city"`
	StarName  string `gorm:"column:star_name" json:"starName"`
	Business  string `gorm:"column:business" json:"business"`
	Longitude string `gorm:"column:longitude" json:"longitude"`
	Latitude  string `gorm:"column:latitude" json:"latitude"`
	Pic       string `gorm:"column:pic" json:"pic"`
}

// TableName 指定表名
func (Hotel) TableName() string {
	return "tb_hotel"
}

// BeforeCreate GORM 钩子
func (h *Hotel) BeforeCreate(tx *gorm.DB) error {
	return nil
}
