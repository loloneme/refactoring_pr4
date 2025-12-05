package specification

import sq "github.com/Masterminds/squirrel"

type GetLastNBySourceSpecification struct {
	idField    string
	limit      int
	sourceName string
}

func NewGetLastNBySourceSpec(idField string, limit int, sourceName string) *GetLastNBySourceSpecification {
	return &GetLastNBySourceSpecification{
		idField:    idField,
		limit:      limit,
		sourceName: sourceName,
	}
}

func (s *GetLastNBySourceSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(sq.Eq{"source": s.sourceName}).
		OrderBy(s.idField + " DESC").
		Limit(uint64(s.limit))
}

func (s *GetLastNBySourceSpecification) GetFields() []string {
	return []string{"*"}
}
