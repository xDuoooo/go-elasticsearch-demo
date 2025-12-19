package dao

import (
	"github.com/xDuoooo/go-elasticsearch-demo/internal/model"
	"gorm.io/gorm"
)

// MysqlDAO MySQL 数据访问层接口
// 负责所有对 MySQL 数据库的操作
type MysqlDAO interface {
	// FindByID 根据ID查询酒店
	FindByID(id int64) (*model.Hotel, error)

	// FindAll 查询所有酒店
	FindAll() ([]*model.Hotel, error)

	// FindByPage 分页查询酒店
	FindByPage(page, pageSize int) ([]*model.Hotel, int64, error)

	// Create 创建酒店
	Create(hotel *model.Hotel) error

	// Update 更新酒店
	Update(hotel *model.Hotel) error

	// Delete 删除酒店
	Delete(id int64) error

	// FindByCity 根据城市查询酒店
	FindByCity(city string) ([]*model.Hotel, error)

	// FindByBrand 根据品牌查询酒店
	FindByBrand(brand string) ([]*model.Hotel, error)
}

// mysqlDAOImpl MySQL DAO 实现
type mysqlDAOImpl struct {
	db *gorm.DB
}

// NewMysqlDAO 创建 MySQL DAO 实例
func NewMysqlDAO(db *gorm.DB) MysqlDAO {
	return &mysqlDAOImpl{db: db}
}

// FindByID 根据ID查询酒店
func (d *mysqlDAOImpl) FindByID(id int64) (*model.Hotel, error) {
	var hotel model.Hotel
	err := d.db.First(&hotel, id).Error
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

// FindAll 查询所有酒店
func (d *mysqlDAOImpl) FindAll() ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	err := d.db.Find(&hotels).Error
	return hotels, err
}

// FindByPage 分页查询酒店
func (d *mysqlDAOImpl) FindByPage(page, pageSize int) ([]*model.Hotel, int64, error) {
	var hotels []*model.Hotel
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	if err := d.db.Model(&model.Hotel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := d.db.Offset(offset).Limit(pageSize).Find(&hotels).Error
	return hotels, total, err
}

// Create 创建酒店
func (d *mysqlDAOImpl) Create(hotel *model.Hotel) error {
	return d.db.Create(hotel).Error
}

// Update 更新酒店
func (d *mysqlDAOImpl) Update(hotel *model.Hotel) error {
	return d.db.Save(hotel).Error
}

// Delete 删除酒店
func (d *mysqlDAOImpl) Delete(id int64) error {
	return d.db.Delete(&model.Hotel{}, id).Error
}

// FindByCity 根据城市查询酒店
func (d *mysqlDAOImpl) FindByCity(city string) ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	err := d.db.Where("city = ?", city).Find(&hotels).Error
	return hotels, err
}

// FindByBrand 根据品牌查询酒店
func (d *mysqlDAOImpl) FindByBrand(brand string) ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	err := d.db.Where("brand = ?", brand).Find(&hotels).Error
	return hotels, err
}
