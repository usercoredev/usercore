package pagination

type Metadata struct {
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

func (p *Metadata) SetTotalCount(count int32) {
	if count < 0 {
		count = 0
	}
	p.TotalCount = count
	if p.PageSize <= 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = count / p.PageSize
		if count%p.PageSize != 0 {
			p.TotalPages++
		}
	}
	p.updateNavigation()
}

func (p *Metadata) SetPage(page int32) {
	if page <= 0 || p.TotalPages == 0 {
		p.Page = 1
	} else if page > p.TotalPages {
		p.Page = p.TotalPages
	} else {
		p.Page = page
	}
	p.updateNavigation()
}

func (p *Metadata) updateNavigation() {
	p.HasPrev = p.Page > 1
	p.HasNext = p.Page < p.TotalPages
}

func (p *Metadata) Offset() int32 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}