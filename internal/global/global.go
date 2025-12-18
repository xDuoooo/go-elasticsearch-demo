package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/itcast/hotel-demo/internal/service"
	"gorm.io/gorm"
)

var (
	Db           *gorm.DB
	HotelService service.HotelService
	EsClient     *elasticsearch.TypedClient
)
