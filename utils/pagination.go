package utils

type PageMetadata struct {
	TotalCount int32 `json:"total_count"`
	TotalPages int32 `json:"total_pages"`
	PageSize   int32 `json:"page_size"`
	Page       int32 `json:"page"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
	OrderBy    string
	Order      string
	Search     string
}

func (p *PageMetadata) SetTotalCount(count int32) {
	if count <= 0 || p.PageSize <= 0 {
		p.TotalCount = 0
		p.TotalPages = 0
	} else {
		p.TotalCount = count
		p.TotalPages = count / p.PageSize
		if count%p.PageSize != 0 {
			p.TotalPages++
		}
	}
	p.updatePageNavigation()
}

func (p *PageMetadata) SetPage(page int32) {
	if page <= 0 {
		p.Page = 1
	} else if page > p.TotalPages {
		p.Page = p.TotalPages
	} else {
		p.Page = page
	}
	p.updatePageNavigation()
}

func (p *PageMetadata) updatePageNavigation() {
	p.HasPrev = p.Page > 1
	p.HasNext = p.Page < p.TotalPages
}

func (p *PageMetadata) Offset() int32 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p *PageMetadata) ConvertToOrder() string {
	if len(p.OrderBy) <= 0 || len(p.Order) <= 0 {
		return "created_at desc"
	}
	return p.OrderBy + " " + p.Order
}
