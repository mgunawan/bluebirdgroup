package paging

import "math"

//PageRequest search request payload
type PageRequest struct {
	Page        int
	ItemPerPage int
}

//PageResponse page response payload
type PageResponse struct {
	Page         int         `json:"page"`
	ItemPerPage  int         `json:"item_per_page"`
	NumberOfPage int         `json:"number_of_page"`
	Total        int       `json:"total"`
	Data         interface{} `json:"data"`
}

//SetDefaultValueIfZero set defaut value if all value is zero
func (p PageRequest) SetDefaultValueIfZero() {
	if p.ItemPerPage == 0 {
		p.ItemPerPage = 10
	}
	if p.Page == 0 {
		p.Page = 1
	}
}

//CalculateOffset calculate offset of pagination
func (p PageRequest) CalculateOffset() int {
	return p.ItemPerPage * (p.Page - 1)
}

//CalculateNumberOfPage calculate number of page
func (p PageRequest) CalculateNumberOfPage(total int) float64 {
	return math.Ceil(float64(total) / float64(p.ItemPerPage))
}
