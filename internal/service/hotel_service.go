package service

import (
	"github.com/xDuoooo/go-elasticsearch-demo/internal/controller/req"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/controller/resp"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/dao"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/model"
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
	ListHotels(req *req.HotelListReq) (*resp.PageResult, error)
}

// hotelService 酒店服务实现
type hotelService struct {
	mysqlDAO dao.MysqlDAO
	esDAO    dao.EsDAO
}

// NewHotelService 创建酒店服务实例
func NewHotelService(mysqlDAO dao.MysqlDAO, esDAO dao.EsDAO) HotelService {
	return &hotelService{
		mysqlDAO: mysqlDAO,
		esDAO:    esDAO,
	}
}

// GetHotelByID 根据ID获取酒店
func (s *hotelService) GetHotelByID(id int64) (*model.HotelDoc, error) {
	hotel, err := s.mysqlDAO.FindByID(id)
	if err != nil {
		return nil, err
	}
	return model.NewHotelDoc(hotel), nil
}

// GetAllHotels 获取所有酒店
func (s *hotelService) GetAllHotels() ([]*model.HotelDoc, error) {
	hotels, err := s.mysqlDAO.FindAll()
	if err != nil {
		return nil, err
	}
	return s.convertToHotelDocs(hotels), nil
}

// GetHotelsByPage 分页获取酒店
func (s *hotelService) GetHotelsByPage(page, pageSize int) ([]*model.HotelDoc, int64, error) {
	hotels, total, err := s.mysqlDAO.FindByPage(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return s.convertToHotelDocs(hotels), total, nil
}

// CreateHotel 创建酒店
func (s *hotelService) CreateHotel(hotel *model.Hotel) error {
	return s.mysqlDAO.Create(hotel)
}

// UpdateHotel 更新酒店
func (s *hotelService) UpdateHotel(hotel *model.Hotel) error {
	return s.mysqlDAO.Update(hotel)
}

// DeleteHotel 删除酒店
func (s *hotelService) DeleteHotel(id int64) error {
	return s.mysqlDAO.Delete(id)
}

// GetHotelsByCity 根据城市获取酒店
func (s *hotelService) GetHotelsByCity(city string) ([]*model.HotelDoc, error) {
	hotels, err := s.mysqlDAO.FindByCity(city)
	if err != nil {
		return nil, err
	}
	return s.convertToHotelDocs(hotels), nil
}

// GetHotelsByBrand 根据品牌获取酒店
func (s *hotelService) GetHotelsByBrand(brand string) ([]*model.HotelDoc, error) {
	hotels, err := s.mysqlDAO.FindByBrand(brand)
	if err != nil {
		return nil, err
	}
	return s.convertToHotelDocs(hotels), nil
}

// ListHotels 搜索酒店列表(使用 Elasticsearch)
func (s *hotelService) ListHotels(searchReq *req.HotelListReq) (*resp.PageResult, error) {
	// 调用 ES DAO 进行搜索
	hotels, total, err := s.esDAO.SearchHotels(searchReq)
	if err != nil {
		return nil, err
	}

	// 构建返回结果
	pageResult := &resp.PageResult{
		Total:  total,
		Hotels: hotels,
	}

	return pageResult, nil
}

// convertToHotelDocs 将 Hotel 列表转换为 HotelDoc 列表
func (s *hotelService) convertToHotelDocs(hotels []*model.Hotel) []*model.HotelDoc {
	docs := make([]*model.HotelDoc, len(hotels))
	for i, hotel := range hotels {
		docs[i] = model.NewHotelDoc(hotel)
	}
	return docs
}
