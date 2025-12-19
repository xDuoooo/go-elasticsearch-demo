package req

type HotelListReq struct {
	Key    string `json:"key"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	SortBy int    `json:"sortBy"`
}
