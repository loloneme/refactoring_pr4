package specification

import sq "github.com/Masterminds/squirrel"

type GetListOfOSDRSpecification struct {
	limit int
}

func NewGetLastNByIdDescSpec(limit int) *GetListOfOSDRSpecification {
	return &GetListOfOSDRSpecification{
		limit: limit,
	}
}

func (s *GetListOfOSDRSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.
		OrderBy("inserted_at DESC").
		Limit(uint64(s.limit))
}

func (s *GetListOfOSDRSpecification) GetFields() []string {
	return []string{"*"}
}
