package data

type Product struct {
	Id     string
	Amount float64
	Fee    float64
}

type Transaction struct {
	Id          string
	ReferenceId string
	WorkflowId  string
	Status      string
}
