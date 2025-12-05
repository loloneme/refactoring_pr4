package specification

import sq "github.com/Masterminds/squirrel"

type GetLastNByIdDescSpecification struct {
	idField string
	limit   int
	fields  []string
}

func NewGetLastNByIdDescSpec(idField string, limit int, fields []string) *GetLastNByIdDescSpecification {
	return &GetLastNByIdDescSpecification{
		idField: idField,
		limit:   limit,
		fields:  fields,
	}
}

func (s *GetLastNByIdDescSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.
		OrderBy(s.idField + " DESC").
		Limit(uint64(s.limit))
}

func (s *GetLastNByIdDescSpecification) GetFields() []string {
	return s.fields
}
