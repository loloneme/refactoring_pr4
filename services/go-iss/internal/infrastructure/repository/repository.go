package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	ErrNotFound = errors.New("entity not found")
)

type dto[ID any, M any] interface {
	Values() []any
	GetID() ID
}

type Repository[ID any, M any, T dto[ID, M]] struct {
	DB        *sqlx.DB
	Columns   *Columns
	TableName string
}

func NewRepository[ID any, M any, T dto[ID, M]](db *sqlx.DB, columns *Columns, tableName string) *Repository[ID, M, T] {
	return &Repository[ID, M, T]{
		DB:        db,
		Columns:   columns,
		TableName: tableName,
	}
}

type FindSpecification interface {
	GetRule(builder sq.SelectBuilder) sq.SelectBuilder
	GetFields() []string
}

func (r *Repository[ID, M, T]) Find(ctx context.Context, spec FindSpecification) ([]M, error) {
	entities := make([]M, 0)

	builder := st.Select(r.Columns.ForSelect(spec.GetFields())...).
		From(r.TableName)

	builder = spec.GetRule(builder)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.DB.SelectContext(ctx, &entities, sqlStr, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return entities, nil
}

func (r *Repository[ID, M, T]) GetIDFieldName(ctx context.Context) string {
	return r.Columns.idField
}

func (r *Repository[ID, M, T]) Save(ctx context.Context, entity T) error {
	return r.SaveMany(ctx, []T{entity})
}

func (r *Repository[ID, M, T]) SaveMany(ctx context.Context, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	updateQuery, updateArgs, err := r.getUpdateQuery(entities)
	if err != nil {
		return fmt.Errorf("generating update query: %w", err)
	}

	_, err = r.DB.ExecContext(ctx, updateQuery, updateArgs...)
	if err != nil {
		return fmt.Errorf("execute update query: %w", err)
	}
	return nil
}

func (r *Repository[ID, M, T]) getUpdateQuery(entities []T) (string, []interface{}, error) {
	queryBuilder := st.
		Insert(r.TableName).
		Columns(append([]string{r.Columns.idField}, r.Columns.ForInsert()...)...)

	for _, entity := range entities {
		values, err := r.GetValuesForEntity(entity)
		if err != nil {
			return "", nil, fmt.Errorf("getting values to update: %w", err)
		}
		queryBuilder = queryBuilder.Values(values...)
	}
	suffix := fmt.Sprintf("ON CONFLICT (id) DO UPDATE SET %s", r.Columns.GetOnConflictStatement())

	return queryBuilder.Suffix(suffix).ToSql()
}

func (r *Repository[ID, M, T]) GetValuesForEntity(entity T) ([]interface{}, error) {
	values := entity.Values()

	entityID := entity.GetID()
	isDef, err := isDefault(entityID)
	if err != nil {
		return nil, err
	}
	if isDef {
		return append([]interface{}{sq.Expr("DEFAULT")}, values...), nil
	}
	return append([]interface{}{entityID}, values...), nil
}

func isDefault(id any) (bool, error) {
	switch value := id.(type) {
	case int64:
		return value == 0, nil
	case uuid.UUID:
		return value == uuid.Nil, nil
	default:
		return false, errors.New("unknown ID type")
	}
}
