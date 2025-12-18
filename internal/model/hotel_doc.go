package model

import "fmt"

// HotelDoc 酒店文档对象，用于 API 响应或 Elasticsearch 文档
type HotelDoc struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Price    int    `json:"price"`
	Score    int    `json:"score"`
	Brand    string `json:"brand"`
	City     string `json:"city"`
	StarName string `json:"starName"`
	Business string `json:"business"`
	Location string `json:"location"` // 组合的经纬度
	Pic      string `json:"pic"`
}

// NewHotelDoc 从 Hotel 实体创建 HotelDoc
func NewHotelDoc(hotel *Hotel) *HotelDoc {
	return &HotelDoc{
		ID:       hotel.ID,
		Name:     hotel.Name,
		Address:  hotel.Address,
		Price:    hotel.Price,
		Score:    hotel.Score,
		Brand:    hotel.Brand,
		City:     hotel.City,
		StarName: hotel.StarName,
		Business: hotel.Business,
		Location: fmt.Sprintf("%s, %s", hotel.Latitude, hotel.Longitude),
		Pic:      hotel.Pic,
	}
}
