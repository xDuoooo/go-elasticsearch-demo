package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xDuoooo/go-elasticsearch-demo/internal/controller"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/dao"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/global"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/model"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/service"
)

func main() {
	// 加载配置
	loadConfig()

	// 初始化数据库
	initDatabase()

	// 初始化 Elasticsearch
	initElasticsearch()

	// 初始化 DAO 层
	initDAO()

	// 初始化 Service 层
	initServices()
	// 初始化路由
	router := initRouter()

	// 启动服务器
	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8089
	}
	log.Printf("服务器启动在端口 %d", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// loadConfig 加载配置文件
func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("读取配置文件失败: %v, 使用默认配置", err)
	}
}

// initDatabase 初始化数据库连接
func initDatabase() {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	// 默认值
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 3306
	}
	if username == "" {
		username = "root"
	}
	if password == "" {
		password = "123"
	}
	if dbname == "" {
		dbname = "heima"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功")

}

// initElasticsearch 初始化 Elasticsearch 客户端
func initElasticsearch() {
	token := viper.GetString("es.token")
	if token == "" {
		log.Fatal("ES token 未配置，请检查配置文件")
	}

	typedClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		APIKey:    token,
		Addresses: []string{viper.GetString("es.address")},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("连接 Elasticsearch 失败: %v", err)
	}

	global.EsClient = typedClient
	log.Println("Elasticsearch 连接成功")
}

// initDAO 初始化 DAO 层
func initDAO() {
	global.MysqlDAO = dao.NewMysqlDAO(global.DB)
	global.EsDAO = dao.NewEsDAO(global.EsClient)
	log.Println("DAO 层初始化成功")
}

// initServices 初始化 Service 层
func initServices() {
	global.HotelService = service.NewHotelService(global.MysqlDAO, global.EsDAO)
	log.Println("Service 层初始化成功")
}

// initRouter 初始化路由
func initRouter() *gin.Engine {
	router := gin.Default()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API 路由组
	api := router.Group("/api")
	{
		hotels := api.Group("/hotel")
		{
			hotels.GET("", listHotels)
			hotels.GET("/:id", getHotel)
			hotels.POST("", createHotel)
			hotels.PUT("/:id", updateHotel)
			hotels.DELETE("/:id", deleteHotel)
			hotels.GET("/city/:city", getHotelsByCity)
			hotels.GET("/brand/:brand", getHotelsByBrand)
			hotels.POST("/list", controller.List)
		}
	}

	return router
}

// listHotels 获取酒店列表（支持分页）
func listHotels(c *gin.Context) {
	page := 1
	pageSize := 10

	if p, ok := c.GetQuery("page"); ok {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps, ok := c.GetQuery("pageSize"); ok {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	hotels, total, err := global.HotelService.GetHotelsByPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     hotels,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// getHotel 根据ID获取酒店
func getHotel(c *gin.Context) {
	var id int64
	fmt.Sscanf(c.Param("id"), "%d", &id)

	hotel, err := global.HotelService.GetHotelByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "酒店不存在"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// createHotel 创建酒店
func createHotel(c *gin.Context) {
	var hotel model.Hotel
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.HotelService.CreateHotel(&hotel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hotel)
}

// updateHotel 更新酒店
func updateHotel(c *gin.Context) {
	var id int64
	fmt.Sscanf(c.Param("id"), "%d", &id)

	var hotel model.Hotel
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotel.ID = id
	if err := global.HotelService.UpdateHotel(&hotel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// deleteHotel 删除酒店
func deleteHotel(c *gin.Context) {
	var id int64
	fmt.Sscanf(c.Param("id"), "%d", &id)

	if err := global.HotelService.DeleteHotel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// getHotelsByCity 根据城市获取酒店
func getHotelsByCity(c *gin.Context) {
	city := c.Param("city")

	hotels, err := global.HotelService.GetHotelsByCity(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

// getHotelsByBrand 根据品牌获取酒店
func getHotelsByBrand(c *gin.Context) {
	brand := c.Param("brand")

	hotels, err := global.HotelService.GetHotelsByBrand(brand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}
