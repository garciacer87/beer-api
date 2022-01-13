package contract

//BoxPriceResp response object used in GetBoxPrice handler
type BoxPriceResp struct {
	Total float64 `json:"totalPrice"`
}
