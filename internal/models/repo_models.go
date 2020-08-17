package models

type RelationRequest struct {
	ProductID  string `json:"product_id"`
	CategoryID string `json:"category_id"`
}

type ListRequest struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Sort       *SortedBy   `json:"sort,omitempty"`
	Assoc      bool
}

type Pagination struct {
	From string
	To   uint16
}

type SortedBy struct {
	FieldName string
	Desc      bool
}
