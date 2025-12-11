package toolscommon

import (
	"net/url"
)

type Transformable interface {
	Transform() Transformable
}

type PaginatedData[T Transformable] struct {
	Data []T `json:"data" jsonschema_description:"Array of items returned in the current page."`

	Paging struct {
		Last     *string `json:"last,omitempty" jsonschema_description:"Last page of results."`
		Next     *string `json:"next,omitempty" jsonschema_description:"Next page of results. If not set, then this is the last page."`
		Previous *string `json:"previous,omitempty" jsonschema_description:"Previous page of results."`
		Self     *string `json:"self,omitempty" jsonschema_description:"Current page of results."`
		PerPage  *string `json:"per_page,omitempty" jsonschema_description:"How many items returned per page."`
	}

	TotalCount *string `json:"totalCount,omitempty" jsonschema_description:"Total number of items available."`
}

func (p PaginatedData[T]) Transform() Transformable {
	for i := range p.Data {
		p.Data[i] = p.Data[i].Transform().(T)
	}

	extractPagePerPage := func(s string) (string, string) {
		u, err := url.Parse(s)
		if err != nil {
			return "", ""
		}
		q := u.Query()
		return q.Get("page"), q.Get("perpage")
	}

	if p.Paging.Self != nil {
		page, perPage := extractPagePerPage(*p.Paging.Self)
		p.Paging.Self = &page
		p.Paging.PerPage = &perPage
	}

	if p.Paging.Last != nil {
		page, _ := extractPagePerPage(*p.Paging.Last)
		p.Paging.Last = &page
	}

	if p.Paging.Next != nil {
		page, _ := extractPagePerPage(*p.Paging.Next)
		p.Paging.Next = &page
	}

	if p.Paging.Previous != nil {
		page, _ := extractPagePerPage(*p.Paging.Previous)
		p.Paging.Previous = &page
	}

	return p
}
