package vo

type Pagination struct {
	Page  uint `json:"page,omitempty"`
	Size  uint `json:"size,omitempty"`
	Total uint `json:"total,omitempty"`
}
