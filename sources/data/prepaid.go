package data

var prepaidLists = map[string]Product{
	"PLN-TOKEN-50": {
		Id:     "PLN-TOKEN-50",
		Amount: 50_000,
		Fee:    1_000,
	},
}

func GetPrepaidProductById(id string) (Product, bool) {
	product, ok := prepaidLists[id]
	return product, ok
}
