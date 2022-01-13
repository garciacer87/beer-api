package contract

//Beer used to represent beer objects
type Beer struct {
	ID       int     `json:"id" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	Brewery  string  `json:"brewery" validate:"required"`
	Country  string  `json:"country" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}
