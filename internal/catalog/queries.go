package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) List(ctx context.Context) ([]Canteen, error) {
	canteens := make([]Canteen, 0)
	canteenIndex := make(map[string]int)

	rows, err := s.pool.Query(ctx, "SELECT id::text, name FROM canteens ORDER BY sort_order, name")
	if err != nil {
		return nil, fmt.Errorf("query canteens: %w", err)
	}
	for rows.Next() {
		var canteen Canteen
		canteen.Floors = make([]Floor, 0)
		if err := rows.Scan(&canteen.ID, &canteen.Name); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan canteen: %w", err)
		}
		canteenIndex[canteen.ID] = len(canteens)
		canteens = append(canteens, canteen)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate canteens: %w", err)
	}
	rows.Close()

	type floorPosition struct{ canteen, floor int }
	floorIndex := make(map[string]floorPosition)
	rows, err = s.pool.Query(ctx, "SELECT id::text, canteen_id::text, name FROM floors ORDER BY sort_order, name")
	if err != nil {
		return nil, fmt.Errorf("query floors: %w", err)
	}
	for rows.Next() {
		var floor Floor
		var canteenID string
		floor.Windows = make([]Window, 0)
		if err := rows.Scan(&floor.ID, &canteenID, &floor.Name); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan floor: %w", err)
		}
		ci, ok := canteenIndex[canteenID]
		if !ok {
			continue
		}
		floorIndex[floor.ID] = floorPosition{canteen: ci, floor: len(canteens[ci].Floors)}
		canteens[ci].Floors = append(canteens[ci].Floors, floor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate floors: %w", err)
	}
	rows.Close()

	rows, err = s.pool.Query(ctx, `
		SELECT id::text, floor_id::text, external_code, name, description, business_hours
		FROM food_windows
		WHERE active = true
		ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, fmt.Errorf("query windows: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var window Window
		var floorID string
		if err := rows.Scan(
			&window.ID, &floorID, &window.ExternalCode, &window.Name,
			&window.Description, &window.BusinessHours,
		); err != nil {
			return nil, fmt.Errorf("scan window: %w", err)
		}
		position, ok := floorIndex[floorID]
		if !ok {
			continue
		}
		floor := &canteens[position.canteen].Floors[position.floor]
		floor.Windows = append(floor.Windows, window)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate windows: %w", err)
	}

	return canteens, nil
}
