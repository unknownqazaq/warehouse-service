package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"warehouse-service/internal/domain/store"
	"warehouse-service/pkg/storage"
)

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (s *StoreRepository) Select(ctx context.Context) (dest []store.Entity, err error) {
	query := `
        SELECT id, merchant_id, city_id, name, address, location, rating, is_active
        FROM stores`

	err = s.db.SelectContext(ctx, &dest, query)

	return
}

func (s *StoreRepository) Create(ctx context.Context, data store.Entity) (id string, err error) {
	query := `
        INSERT INTO stores (merchant_id, city_id, name, address, location)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	args := []interface{}{data.MerchantID, data.CityID, data.Name, data.Address, data.Location}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *StoreRepository) Get(ctx context.Context, id string) (dest store.Entity, err error) {
	query := `
        SELECT id, merchant_id, city_id, name, address, location, rating, is_active
        FROM stores
        WHERE id=$1`

	args := []interface{}{id}

	if err = s.db.GetContext(ctx, &dest, query, args...); err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = storage.ErrorNotFound
	}

	return
}

func (s *StoreRepository) Update(ctx context.Context, id string, data store.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE stores SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = s.db.ExecContext(ctx, query, args...)
		if err != nil && err != sql.ErrNoRows {
			return
		}

		if err == sql.ErrNoRows {
			err = storage.ErrorNotFound
		}
	}

	return
}

func (s *StoreRepository) prepareArgs(data store.Entity) (sets []string, args []any) {

	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Address != nil {
		args = append(args, data.Address)
		sets = append(sets, fmt.Sprintf("address=$%d", len(args)))
	}

	if data.Location != nil {
		args = append(args, data.Location)
		sets = append(sets, fmt.Sprintf("location=$%d", len(args)))
	}

	if data.Rating != nil {
		args = append(args, data.Rating)
		sets = append(sets, fmt.Sprintf("rating=$%d", len(args)))
	}

	return
}

func (s *StoreRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
        DELETE 
        FROM stores
        WHERE id=$1`

	args := []interface{}{id}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = storage.ErrorNotFound
	}

	return
}
