package entity

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) IncrementLimit() Pagination {
	return Pagination{Limit: p.Limit + 1, Offset: p.Offset}
}

type GetRevisionsData struct {
	Claims    Claims
	Status    EntityStatus
	DisputeID *string
}

type GetCorrespondencesData struct {
	Claims     Claims
	RevisionID string
}

type paginated[T any] struct {
	Data        []T
	HasMoreData bool
}

type PaginatedRevisions paginated[Revision]

type GetOrganizationsData struct {
	Claims           Claims
	OrganizationCode *string
	SearchToken      *string
}

type PaginatedOrganizations paginated[Organization]
