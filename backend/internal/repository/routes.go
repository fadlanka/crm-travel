package repository

import (
	"context"

	"travelcrm/internal/model"
)

// ListRoutes returns all routes ordered by origin then destination.
func (r *Repository) ListRoutes(ctx context.Context) ([]model.Route, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, origin, destination, created_at, updated_at
		FROM routes
		ORDER BY origin, destination`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []model.Route
	for rows.Next() {
		var rt model.Route
		if err := rows.Scan(&rt.ID, &rt.Origin, &rt.Destination, &rt.CreatedAt, &rt.UpdatedAt); err != nil {
			return nil, err
		}
		routes = append(routes, rt)
	}
	return routes, rows.Err()
}
