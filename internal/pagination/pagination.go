package pagination

type Metadata struct {
	TotalCount int32 `json:"total_count"`
	TotalPages int32 `json:"total_pages"`
	PageSize   int32 `json:"page_size"`
	Page       int32 `json:"page"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func (p *Metadata) SetTotalCount(count int32) {
	p.TotalCount = max(count, 0)
	p.TotalPages = max(p.TotalCount/p.PageSize+min(1, p.TotalCount%p.PageSize), 0)
	p.updateNavigation()
}

func (p *Metadata) SetPage(page int32) {
	p.Page = clamp(page, 1, p.TotalPages)
	p.updateNavigation()
}

func (p *Metadata) updateNavigation() {
	p.HasPrev = p.Page > 1
	p.HasNext = p.Page < p.TotalPages
}

func (p *Metadata) Offset() int32 {
	return max((p.Page-1)*p.PageSize, 0)
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func clamp(val, min, max int32) int32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
