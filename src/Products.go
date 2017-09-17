package src

type (
	Product struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		FileName    string  `json:"fileName"`
		WebUrl      string  `json:"url"`
		ImageUrl    string  `json:"imageUrl"`
		Payload     string  `json:"payload"`
		//ignore this field from JSON parsing
		FilePath string `json:"-" `
	}
)

var (
	Products   = map[int]*Product{}
	ProductSeq = 1
)

func NewProduct() *Product {
	p := new(Product)
	p.Payload = "fdsaf"
	p.ID = ProductSeq
	return p
}

func GetProductsList() []Product {
	keys := make([]Product, len(Products))

	i := 0
	for _, v := range Products {
		keys[i] = *v
		i++
	}
	return keys
}
