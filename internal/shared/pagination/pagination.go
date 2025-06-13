package pagination

// PaginationParams represents the pagination parameters
type PaginationParams struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

func (p *PaginationParams) Offset() int32 {
	return (p.Page - 1) * p.Limit
}

func (p *PaginationParams) NextPage(total int64) *int32 {
	nextPage := p.Page + 1
	if int64(nextPage*p.Limit) >= total {
		return nil
	}

	return &nextPage
}

// PaginateResult represents a paginated result set
type PaginateResult[T any] struct {
	Data       []T     `json:"data"`
	Limit      int32   `json:"limit"`
	Page       int32   `json:"page"`
	Total      int64   `json:"total"`
	NextPage   *int32  `json:"next_page,omitempty"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

func (p PaginateResult[T]) GetData() []T {
	return p.Data
}

func (p PaginateResult[T]) GetLimit() int32 {
	return p.Limit
}

func (p PaginateResult[T]) GetPage() int32 {
	return p.Page
}

func (p PaginateResult[T]) GetTotal() int64 {
	return p.Total
}

func (p PaginateResult[T]) GetNextPage() *int32 {
	return p.NextPage
}

func (p PaginateResult[T]) GetNextCursor() *string {
	return p.NextCursor
}
