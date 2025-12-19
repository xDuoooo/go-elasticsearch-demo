package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/dao"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/service"
	"gorm.io/gorm"
)

// 全局变量定义
// 类似 Java Spring 的 Bean 容器，统一管理应用级别的单例对象

var (
	// DB MySQL 数据库连接
	DB *gorm.DB

	// EsClient Elasticsearch 客户端
	EsClient *elasticsearch.TypedClient

	// DAO 层实例
	MysqlDAO dao.MysqlDAO
	EsDAO    dao.EsDAO

	// Service 层实例
	HotelService service.HotelService
)
