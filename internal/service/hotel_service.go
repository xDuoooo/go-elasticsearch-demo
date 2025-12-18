package service

import (
	"github.com/xDuoooo/go-elasticsearch-demo/internal/model"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/repository"
)

// HotelService 酒店服务接口
type HotelService interface {
	GetHotelByID(id int64) (*model.HotelDoc, error)
	GetAllHotels() ([]*model.HotelDoc, error)
	GetHotelsByPage(page, pageSize int) ([]*model.HotelDoc, int64, error)
	CreateHotel(hotel *model.Hotel) error
	UpdateHotel(hotel *model.Hotel) error
	DeleteHotel(id int64) error
	GetHotelsByCity(city string) ([]*model.HotelDoc, error)
	GetHotelsByBrand(brand string) ([]*model.HotelDoc, error)
}

// hotelService 酒店服务实现
type hotelService struct {
	repo repository.HotelRepository
}

// NewHotelService 创建酒店服务实例
func NewHotelService(repo repository.HotelRepository) HotelService {
	return &hotelService{repo: repo}
}

// GetHotelByID 根据ID获取酒店
func (s *hotelService) GetHotelByID(id int64) (*model.HotelDoc, error) {
	hotel, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return model.NewHotelDoc(hotel), nil
}

// GetAllHotels 获取所有酒店
func (s *hotelService) GetAllHotels() ([]*model.HotelDoc, error) {
	hotels, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.convertToHotelDocs(hotels), nil
}

// GetHotelsByPage 分页获取酒店
func (s *hotelService) GetHotelsByPage(page, pageSize int) ([]*model.HotelDoc, int64, error) {
	hotels, total, err := s.repo.FindByPage(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return s.convertToHotelDocs(hotels), total, nil
}

// CreateHotel 创建酒店
func (s *hotelService) CreateHotel(hotel *model.Hotel) error {
	return s.repo.Create(hotel)
}

// UpdateHotel 更新酒店
func (s *hotelService) UpdateHotel(hotel *model.Hotel) error {
	return s.repo.Update(hotel)
}

// DeleteHotel 删除酒店
func (s *hotelService) DeleteHotel(id int64) error {
	return s.repo.Delete(id)
}

// GetHotelsByCity 根据城市获取酒店
func (s *hotelService) GetHotelsByCity(city string) ([]*model.HotelDoc, error) {
	hotels, err := s.repo.FindByCity(city)
	if err != nil {
		return nil, err
	}
	return s.convertToHotelDocs(hotels), nil
}

// GetHotelsByBrand 根据品牌获取酒店
func (s *hotelService) GetHotelsByBrand(brand string) ([]*model.HotelDoc, error) {
	hotels, err := s.repo.FindByBrand(brand)
	if err != nil {
		return nil, err
	}
	return s.convertToHotelDocs(hotels), nil
}

// convertToHotelDocs 将 Hotel 列表转换为 HotelDoc 列表
func (s *hotelService) convertToHotelDocs(hotels []*model.Hotel) []*model.HotelDoc {
	docs := make([]*model.HotelDoc, len(hotels))
	for i, hotel := range hotels {
		docs[i] = model.NewHotelDoc(hotel)
	}
	return docs
}
