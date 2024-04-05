package postgres

const defaultLimit = 1000

type Pagination struct {
	limit  uint64
	offset uint64
}

func (p *Pagination) Limit() uint64 {
	return p.limit
}

func (p *Pagination) Offset() uint64 {
	return p.offset
}

func NewPagination(page, perPage uint64) Pagination {
	p := Pagination{}

	p.limit = perPage
	if perPage == 0 || perPage > defaultLimit {
		p.limit = defaultLimit
	}

	if page < 1 {
		page = 1
	}
	p.offset = (page - 1) * p.limit

	return p
}
