package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/controller/req"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/global"
)

// List 查询酒店列表
func List(c *gin.Context) {
	var searchReq req.HotelListReq
	if err := c.ShouldBindJSON(&searchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := global.HotelService.ListHotels(&searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
