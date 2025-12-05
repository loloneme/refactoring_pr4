package osdr

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) SaveOSDR(ctx context.Context, entity OSDRItem) error {
	builder := sq.Insert("osdr_items").
		Columns(append([]string{r.Columns.GetIDField()}, r.Columns.ForInsert()...)...).
		Values(r.GetValuesForEntity(entity)).
		PlaceholderFormat(sq.Dollar)

	if entity.DatasetID != nil {
		onConflict := r.Columns.GetOnConflictStatement()
		builder = builder.Suffix(fmt.Sprintf("ON CONFLICT (dataset_id) DO UPDATE SET %s", onConflict))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	builder := sq.Select("count(*)").From(r.TableName)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	err = r.DB.SelectContext(ctx, &count, sqlStr, args...)
	return count, err
}
