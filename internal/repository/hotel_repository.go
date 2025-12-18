package repository

import (
	"github.com/xDuoooo/go-elasticsearch-demo/internal/model"
	"gorm.io/gorm"
)

// HotelRepository 酒店数据访问接口
type HotelRepository interface {
	FindByID(id int64) (*model.Hotel, error)
	FindAll() ([]*model.Hotel, error)
	FindByPage(page, pageSize int) ([]*model.Hotel, int64, error)
	Create(hotel *model.Hotel) error
	Update(hotel *model.Hotel) error
	Delete(id int64) error
	FindByCity(city string) ([]*model.Hotel, error)
	FindByBrand(brand string) ([]*model.Hotel, error)
}

// hotelRepository 酒店数据访问实现
type hotelRepository struct {
	db *gorm.DB
}

// NewHotelRepository 创建酒店数据访问实例
func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{db: db}
}

// FindByID 根据ID查询酒店
func (r *hotelRepository) FindByID(id int64) (*model.Hotel, error) {
	var hotel model.Hotel
	err := r.db.First(&hotel, id).Error
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

// FindAll 查询所有酒店
func (r *hotelRepository) FindAll() ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	err := r.db.Find(&hotels).Error
	return hotels, err
}

// FindByPage 分页查询酒店
func (r *hotelRepository) FindByPage(page, pageSize int) ([]*model.Hotel, int64, error) {
	var hotels []*model.Hotel
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	if err := r.db.Model(&model.Hotel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.Offset(offset).Limit(pageSize).Find(&hotels).Error
	return hotels, total, err
}

// Create 创建酒店
func (r *hotelRepository) Create(hotel *model.Hotel) error {
	return r.db.Create(hotel).Error
}

// Update 更新酒店
func (r *hotelRepository) Update(hotel *model.Hotel) error {
	return r.db.Save(hotel).Error
}

// Delete 删除酒店
func (r *hotelRepository) Delete(id int64) error {
	return r.db.Delete(&model.Hotel{}, id).Error
}

// FindByCity 根据城市查询酒店
func (r *hotelRepository) FindByCity(city string) ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	err := r.db.Where("city = ?", city).Find(&hotels).Error
	return hotels, err
}

// FindByBrand 根据品牌查询酒店
func (r *hotelRepository) FindByBrand(brand string) ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	err := r.db.Where("brand = ?", brand).Find(&hotels).Error
	return hotels, err
}
