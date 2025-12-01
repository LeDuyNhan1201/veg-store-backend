package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/iface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"

	"github.com/iancoleman/strcase"
)

type Repository[TEntity iface.IEntity, TId model.AllowedId] struct {
	*core.Core
}

func NewRepository[TEntity iface.IEntity, TId model.AllowedId](core *core.Core) *Repository[TEntity, TId] {
	return &Repository[TEntity, TId]{
		Core: core,
	}
}

func (r *Repository[TEntity, TId]) Create(db *data.PostgresDB, ctx context.Context, entity TEntity) error {
	entity.Created()
	return db.WithContext(ctx).Create(entity).Error
}

func (r *Repository[TEntity, TId]) FindById(
	db *data.PostgresDB, ctx context.Context,
	id TId,
	opt ...dto.FindByIDOption,
) (TEntity, error) {
	var entity TEntity
	queryBuilder := db.WithContext(ctx)

	// ---------------------------------
	// APPLY PRELOADS
	// ---------------------------------
	if len(opt) > 0 {
		for _, preload := range opt[0].Preloads {
			queryBuilder = db.Preload(preload)
		}
	}

	// ---------------------------------
	// Query
	// ---------------------------------
	query, args := buildWhereFromID(id)
	err := queryBuilder.Where(query, args...).First(&entity).Error

	return entity, err
}

func (r *Repository[TEntity, TId]) FindAll(db *data.PostgresDB, ctx context.Context) ([]TEntity, error) {
	var result []TEntity
	err := db.WithContext(ctx).Where("is_deleted = false").Find(&result).Error
	return result, err
}

func (r *Repository[TEntity, TId]) OffsetPage(
	db *data.PostgresDB, ctx context.Context,
	opt dto.OffsetPageOption,
) (dto.OffsetPageResult[TEntity], error) {
	var result dto.OffsetPageResult[TEntity]
	var items []TEntity
	var total int64

	queryBuilder := db.WithContext(ctx)

	// ------------------------
	// WHERE conditions
	// ------------------------
	for _, whereCondition := range opt.Where {
		field := strcase.ToSnake(whereCondition.Field)
		queryBuilder = queryBuilder.Where(fmt.Sprintf("%s %s ?", field, whereCondition.Operator.String()), whereCondition.Value)
	}

	// ------------------------
	// SORT conditions
	// ------------------------
	for _, sortCondition := range opt.Sort {
		field := strcase.ToSnake(sortCondition.Field)
		direction := sortCondition.Direction
		if direction.IsValid() {
			direction = dto.Asc
		}
		queryBuilder = queryBuilder.Order(fmt.Sprintf("%s %s", field, direction))
	}

	// ------------------------
	// PRELOAD relations
	// ------------------------
	for _, relation := range opt.Preload {
		queryBuilder = queryBuilder.Preload(relation)
	}

	// ------------------------
	// COUNT total
	// ------------------------
	var entity TEntity
	if err := queryBuilder.Model(&entity).Count(&total).Error; err != nil {
		return result, err
	}

	// ------------------------
	// Pagination
	// ------------------------
	if opt.Page < 1 {
		opt.Page = 1
	}
	if opt.Size < 1 {
		opt.Size = 20
	}

	offset := (opt.Page - 1) * opt.Size

	// ------------------------
	// Query page items
	// ------------------------
	if err := queryBuilder.
		Limit(int(opt.Size)).
		Offset(int(offset)).
		Find(&items).
		Error; err != nil {
		return result, err
	}

	result = dto.OffsetPageResult[TEntity]{
		Items: items,
		Page:  opt.Page,
		Size:  opt.Size,
		Total: total,
	}

	return result, nil
}

func (r *Repository[TEntity, TId]) Update(db *data.PostgresDB, ctx context.Context, entity TEntity) error {
	entity.Updated()
	return db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[TEntity, TId]) SoftDelete(db *data.PostgresDB, ctx context.Context, id TId) error {
	entity, err := r.FindById(db, ctx, id)
	if err != nil {
		return err
	}

	entity.Deleted()
	return db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[TEntity, TId]) HardDelete(db *data.PostgresDB, ctx context.Context, id TId) error {
	query, args := buildWhereFromID(id)
	return db.WithContext(ctx).Where(query, args...).Delete(new(TEntity), id).Error
}

func buildWhereFromID[TID any](id TID) (string, []any) {
	value := reflect.ValueOf(id)
	typeOf := reflect.TypeOf(id)

	if typeOf.Kind() != reflect.Struct {
		return "id = ?", []any{id}
	}

	var parts []string
	var args []any

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := strcase.ToSnake(field.Name)
		parts = append(parts, fmt.Sprintf("%s = ?", column))
		args = append(args, value.Field(i).Interface())
	}

	return strings.Join(parts, " AND "), args
}
