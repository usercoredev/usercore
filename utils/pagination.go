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

func (l *PageMetadata) SetTotalCount(count int32) {
	if count <= 0 || l.PageSize <= 0 {
		l.TotalCount = 0
		l.TotalPages = 0
		return
	}

	l.TotalCount = count
	l.TotalPages = count / l.PageSize
	if count%l.PageSize != 0 {
		l.TotalPages++
	}
}

func (l *PageMetadata) SetPage(page int32) {
	if page <= 0 {
		page = 1
	}

	l.Page = page
	l.HasPrev = page > 1
	l.HasNext = page < l.TotalPages
}

func (l *PageMetadata) Offset() int32 {
	return (l.Page - 1) * l.PageSize
}

func (l *PageMetadata) ConvertToOrder() string {
	if len(l.OrderBy) <= 0 || len(l.Order) <= 0 {
		return "created_at desc"
	}
	return l.OrderBy + " " + l.Order
}
