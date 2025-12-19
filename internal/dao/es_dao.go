package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/controller/req"
	"github.com/xDuoooo/go-elasticsearch-demo/internal/model"
)

// EsDAO Elasticsearch 数据访问层接口
// 负责所有对 Elasticsearch 的操作
type EsDAO interface {
	// 索引管理
	CreateIndex(indexName string) error
	DeleteIndex(indexName string) error
	ExistsIndex(indexName string) (bool, error)

	// 文档操作
	CreateDocument(indexName string, id string, doc interface{}) error
	UpdateDocument(indexName string, id string, doc map[string]interface{}) error
	DeleteDocument(indexName string, id string) error
	GetDocument(indexName string, id string) (*model.HotelDoc, error)

	// 查询操作
	SearchHotels(searchReq *req.HotelListReq) ([]*model.HotelDoc, int64, error)

	// 批量操作
	BulkIndex(indexName string, hotels []*model.HotelDoc) error
}

// esDAOImpl Elasticsearch DAO 实现
type esDAOImpl struct {
	client *elasticsearch.TypedClient
}

// NewEsDAO 创建 Elasticsearch DAO 实例
func NewEsDAO(client *elasticsearch.TypedClient) EsDAO {
	return &esDAOImpl{client: client}
}

// CreateIndex 创建酒店索引
func (d *esDAOImpl) CreateIndex(indexName string) error {
	mapping := &types.TypeMapping{
		Properties: map[string]types.Property{
			"id": types.KeywordProperty{},
			"name": types.TextProperty{
				Analyzer: ptr("ik_max_word"),
				CopyTo:   []string{"all"},
			},
			"address": types.KeywordProperty{
				Index: ptr(false),
			},
			"price":    types.IntegerNumberProperty{},
			"score":    types.IntegerNumberProperty{},
			"brand":    types.KeywordProperty{CopyTo: []string{"all"}},
			"city":     types.KeywordProperty{},
			"starName": types.KeywordProperty{},
			"business": types.KeywordProperty{CopyTo: []string{"all"}},
			"pic":      types.KeywordProperty{Index: ptr(false)},
			"location": types.GeoPointProperty{},
			"all":      types.TextProperty{Analyzer: ptr("ik_max_word")},
		},
	}

	_, err := d.client.Indices.Create(indexName).Mappings(mapping).Do(context.Background())
	return err
}

// DeleteIndex 删除索引
func (d *esDAOImpl) DeleteIndex(indexName string) error {
	_, err := d.client.Indices.Delete(indexName).Do(context.Background())
	return err
}

// ExistsIndex 判断索引是否存在
func (d *esDAOImpl) ExistsIndex(indexName string) (bool, error) {
	exists, err := d.client.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateDocument 创建文档
func (d *esDAOImpl) CreateDocument(indexName string, id string, doc interface{}) error {
	_, err := d.client.Create(indexName, id).Document(doc).Do(context.Background())
	return err
}

// UpdateDocument 更新文档
func (d *esDAOImpl) UpdateDocument(indexName string, id string, doc map[string]interface{}) error {
	_, err := d.client.Update(indexName, id).Doc(doc).Do(context.Background())
	return err
}

// DeleteDocument 删除文档
func (d *esDAOImpl) DeleteDocument(indexName string, id string) error {
	_, err := d.client.Delete(indexName, id).Do(context.Background())
	return err
}

// GetDocument 根据ID获取文档
func (d *esDAOImpl) GetDocument(indexName string, id string) (*model.HotelDoc, error) {
	res, err := d.client.Get(indexName, id).Do(context.Background())
	if err != nil {
		return nil, err
	}

	if res.Found {
		var hotelDoc model.HotelDoc
		if err := json.Unmarshal(res.Source_, &hotelDoc); err != nil {
			return nil, err
		}
		return &hotelDoc, nil
	}

	return nil, fmt.Errorf("文档不存在, ID: %s", id)
}

// SearchHotels 搜索酒店
func (d *esDAOImpl) SearchHotels(searchReq *req.HotelListReq) ([]*model.HotelDoc, int64, error) {
	// 构建查询
	searchBuilder := d.client.Search().Index("hotel")

	// 构建查询条件
	if searchReq.Key != "" {
		searchBuilder = searchBuilder.Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"all": {
					Query: searchReq.Key,
				},
			},
		})
	} else {
		searchBuilder = searchBuilder.Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		})
	}

	// 分页
	page := searchReq.Page
	size := searchReq.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	from := (page - 1) * size
	searchBuilder = searchBuilder.From(from).Size(size)

	// 排序
	// sortBy: 0-默认排序, 1-价格升序, 2-价格降序, 3-评分升序, 4-评分降序
	switch searchReq.SortBy {
	case 1:
		searchBuilder = searchBuilder.Sort(types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"price": {Order: &sortorder.Asc},
			},
		})
	case 2:
		searchBuilder = searchBuilder.Sort(types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"price": {Order: &sortorder.Desc},
			},
		})
	case 3:
		searchBuilder = searchBuilder.Sort(types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"score": {Order: &sortorder.Asc},
			},
		})
	case 4:
		searchBuilder = searchBuilder.Sort(types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"score": {Order: &sortorder.Desc},
			},
		})
	}

	// 高亮
	if searchReq.Key != "" {
		searchBuilder = searchBuilder.Highlight(&types.Highlight{
			Fields: map[string]types.HighlightField{
				"name": {
					RequireFieldMatch: ptr(false),
				},
			},
			PreTags:  []string{"<em>"},
			PostTags: []string{"</em>"},
		})
	}

	// 执行查询
	searchResult, err := searchBuilder.Do(context.Background())
	if err != nil {
		return nil, 0, err
	}

	// 解析结果
	total := searchResult.Hits.Total.Value
	hotels := make([]*model.HotelDoc, 0, len(searchResult.Hits.Hits))

	for _, hit := range searchResult.Hits.Hits {
		var hotelDoc model.HotelDoc
		if err := json.Unmarshal(hit.Source_, &hotelDoc); err != nil {
			continue
		}

		// 处理高亮
		if nameHighlight, ok := hit.Highlight["name"]; ok && len(nameHighlight) > 0 {
			hotelDoc.Name = nameHighlight[0]
		}

		hotels = append(hotels, &hotelDoc)
	}

	return hotels, total, nil
}

// BulkIndex 批量索引文档
func (d *esDAOImpl) BulkIndex(indexName string, hotels []*model.HotelDoc) error {
	bulkAPI := d.client.Bulk().Index(indexName)

	for _, hotel := range hotels {
		operation := types.IndexOperation{
			Index_: ptr(indexName),
			Id_:    ptr(strconv.FormatInt(hotel.ID, 10)),
		}
		if err := bulkAPI.IndexOp(operation, hotel); err != nil {
			return err
		}
	}

	res, err := bulkAPI.Do(context.Background())
	if err != nil {
		return err
	}

	// 检查是否有错误
	if res.Errors {
		return errors.New("批量索引存在错误")
	}

	return nil
}

// ptr 辅助函数，返回指针
func ptr[T any](v T) *T {
	return &v
}
