package resp

import "github.com/xDuoooo/go-elasticsearch-demo/internal/model"

type PageResult struct {
	Total  int64             `json:"total"`
	Hotels []*model.HotelDoc `json:"hotels"`
}
