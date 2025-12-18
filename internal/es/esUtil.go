package es

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/itcast/hotel-demo/internal/global"
	"log"
)

func ptr[T any](v T) *T { return &v }

func CreateHotelIndex() {
	var HotelMapping = &types.TypeMapping{
		Properties: map[string]types.Property{
			"id": types.KeywordProperty{},
			"name": types.TextProperty{
				Analyzer: ptr("ik_max_word"),
				CopyTo:   []string{"all"},
			},
			"address": types.KeywordProperty{
				Index: ptr(false),
			},
			"price": types.IntegerNumberProperty{},
			"score": types.IntegerNumberProperty{},
			"brand": types.KeywordProperty{
				CopyTo: []string{"all"},
			},
			"city":     types.KeywordProperty{},
			"starName": types.KeywordProperty{},
			"business": types.KeywordProperty{
				CopyTo: []string{"all"},
			},
			"pic": types.KeywordProperty{
				Index: ptr(false),
			},
			"location": types.GeoPointProperty{},
			"all": types.TextProperty{
				Analyzer: ptr("ik_max_word"),
			},
		},
	}
	res, err := global.EsClient.Indices.Create("hotel").Mappings(HotelMapping).Do(context.Background())
	fmt.Println(res)
	if err != nil {
		log.Fatalln("创建索引库失败", err)
	}

}
func DeleteHotelIndex(indexName string) {
	_, err := global.EsClient.Indices.Delete(indexName).Do(context.Background())
	if err != nil {
		log.Fatalln("删除索引库失败", err)
	}
}
func ExistHotelIndex(indexName string) bool {
	res, err := global.EsClient.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		log.Fatalln("判断是否存在索引库失败", err)
	}
	return res
}

func GetHotel() {
	res, err := global.EsClient.Get("heima", "1").Do(context.TODO())
	if err != nil {
		log.Fatalln("创建索引失败", err)
	}
	fmt.Println(res)
}
func ListAllIndices() {
	// 1. 调用 Cat.Indices 接口
	resp, err := global.EsClient.Cat.Indices().Do(context.Background())
	if err != nil {
		log.Fatalf("获取索引列表失败: %v", err)
	}

	fmt.Println("当前集群中的索引列表：")
	for _, index := range resp {
		// index.Index 包含索引名称
		fmt.Printf("名称: %-20s | 状态: %-10s | 文档数: %s\n",
			*index.Index, *index.Health, *index.DocsCount)
	}
}
