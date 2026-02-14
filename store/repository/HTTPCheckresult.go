package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mohamedbeat/pulse/store"
)

type HTTPCheckRepository struct {
	db *sqlx.DB
	sb sq.StatementBuilderType
}

func NewHTTPCheckRepository(db *sqlx.DB) *HTTPCheckRepository {
	return &HTTPCheckRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Save a new HTTP check result
func (r *HTTPCheckRepository) Save(ctx context.Context, check *store.HTTPCheck) error {
	dbCheck := check.ToDB()

	// Use squirrel for INSERT with named parameters
	query, args, err := r.sb.Insert("http_check_results").
		Columns(
			"name", "url", "method", "interval_seconds", "timeout_seconds", "expected_status",
			"headers", "body", "must_match_status", "body_contains", "body_regex",
			"status", "status_code", "content_length", "response_headers", "response_body",
			"error_message", "duration_ms", "checked_at",
		).
		Values(
			dbCheck.Name, dbCheck.URL, dbCheck.Method, dbCheck.IntervalSeconds, dbCheck.TimeoutSeconds, dbCheck.ExpectedStatus,
			dbCheck.Headers, dbCheck.Body, dbCheck.MustMatchStatus, dbCheck.BodyContains, dbCheck.BodyRegex,
			dbCheck.Status, dbCheck.StatusCode, dbCheck.ContentLength, dbCheck.ResponseHeaders, dbCheck.ResponseBody,
			dbCheck.ErrorMessage, dbCheck.DurationMS, dbCheck.CheckedAt,
		).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&check.ID, &check.CreatedAt, &check.UpdatedAt)
	}

	return rows.Err()
}

// Get latest check for an endpoint
func (r *HTTPCheckRepository) GetLatest(ctx context.Context, name string) (*store.HTTPCheck, error) {
	query, args, err := r.sb.Select("*").
		From("http_check_results").
		Where(sq.Eq{"name": name}).
		OrderBy("checked_at DESC").
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbCheck store.DBHTTPCheck
	err = r.db.GetContext(ctx, &dbCheck, query, args...)
	if err != nil {
		return nil, err
	}

	return dbCheck.ToDomain(), nil
}

// Get check history for an endpoint
func (r *HTTPCheckRepository) GetHistory(ctx context.Context, name string, hours int, limit int) ([]*store.HTTPCheck, error) {
	query, args, err := r.sb.Select("*").
		From("http_check_results").
		Where(sq.Eq{"name": name}).
		Where(sq.GtOrEq{"checked_at": time.Now().Add(-time.Duration(hours) * time.Hour)}).
		OrderBy("checked_at DESC").
		Limit(uint64(limit)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbChecks []store.DBHTTPCheck
	err = r.db.SelectContext(ctx, &dbChecks, query, args...)
	if err != nil {
		return nil, err
	}

	checks := make([]*store.HTTPCheck, len(dbChecks))
	for i, dbCheck := range dbChecks {
		checks[i] = dbCheck.ToDomain()
	}

	return checks, nil
}

// Get all active endpoints (distinct names with recent checks)
func (r *HTTPCheckRepository) GetActiveEndpoints(ctx context.Context, since time.Duration) ([]string, error) {
	query, args, err := r.sb.Select("DISTINCT name").
		From("http_check_results").
		Where(sq.GtOrEq{"checked_at": time.Now().Add(-since)}).
		OrderBy("name").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var names []string
	err = r.db.SelectContext(ctx, &names, query, args...)
	return names, err
}

// Get endpoint summary (latest status for all endpoints)
func (r *HTTPCheckRepository) GetSummary(ctx context.Context) ([]*EndpointSummary, error) {
	// DISTINCT ON requires raw SQL or specific handling
	// query := `
	//        SELECT DISTINCT ON (name)
	//            name,
	//            url,
	//            status,
	//            status_code,
	//            duration_ms,
	//            checked_at,
	//            error_message
	//        FROM http_check_results
	//        ORDER BY name, checked_at DESC
	//    `
	//
	// var summaries []*EndpointSummary
	// err := r.db.SelectContext(ctx, &summaries, query)
	// return summaries, err
	query, args, err := r.sb.Select("DISTINCT name").
		From("http_check_results").
		OrderBy("name, checked_at DESC").ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var summaries []*EndpointSummary
	err = r.db.SelectContext(ctx, &summaries, query, args...)
	return summaries, err

}

// Get endpoint summary with pagination using Squirrel
func (r *HTTPCheckRepository) GetSummaryPaginated(ctx context.Context, limit, offset int) ([]*EndpointSummary, error) {
	// Get distinct endpoint names first
	namesQuery, namesArgs, err := r.sb.Select("DISTINCT name").
		From("http_check_results").
		OrderBy("name").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build names query: %w", err)
	}

	var names []string
	err = r.db.SelectContext(ctx, &names, namesQuery, namesArgs...)
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return []*EndpointSummary{}, nil
	}

	// Get latest check for each name
	query, args, err := r.sb.Select(
		"name",
		"url",
		"status",
		"status_code",
		"duration_ms",
		"checked_at",
		"error_message",
	).
		From("http_check_results h1").
		Where(sq.Eq{"name": names}).
		Where(sq.Expr("checked_at = (SELECT MAX(checked_at) FROM http_check_results h2 WHERE h2.name = h1.name)")).
		OrderBy("name").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build summary query: %w", err)
	}

	var summaries []*EndpointSummary
	err = r.db.SelectContext(ctx, &summaries, query, args...)
	return summaries, err
}

// Get statistics for dashboard
func (r *HTTPCheckRepository) GetStats(ctx context.Context, hours int) (*HTTPCheckStats, error) {
	// Complex aggregations are better with raw SQL, but we can use Squirrel for conditions
	baseQuery := r.sb.Select(
		"COUNT(DISTINCT name) as total_endpoints",
		"COUNT(*) as total_checks",
		"SUM(CASE WHEN status = 'up' THEN 1 ELSE 0 END) as successful_checks",
		"SUM(CASE WHEN status = 'down' THEN 1 ELSE 0 END) as failed_checks",
		"SUM(CASE WHEN status = 'degraded' THEN 1 ELSE 0 END) as degraded_checks",
		"ROUND(AVG(duration_ms) FILTER (WHERE status = 'up'), 2) as avg_latency_ms",
		"MAX(checked_at) as last_checked_at",
	).From("http_check_results")

	if hours > 0 {
		baseQuery = baseQuery.Where(sq.GtOrEq{
			"checked_at": time.Now().Add(-time.Duration(hours) * time.Hour),
		})
	}

	query, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build stats query: %w", err)
	}

	var stats HTTPCheckStats
	err = r.db.GetContext(ctx, &stats, query, args...)
	return &stats, err
}

// Get failing checks for alerts
func (r *HTTPCheckRepository) GetFailingChecks(ctx context.Context, minutes int) ([]*store.HTTPCheck, error) {
	query, args, err := r.sb.Select("*").
		From("http_check_results").
		Where(sq.Or{
			sq.Eq{"status": "down"},
			sq.Eq{"status": "degraded"},
		}).
		Where(sq.GtOrEq{"checked_at": time.Now().Add(-time.Duration(minutes) * time.Minute)}).
		OrderBy("checked_at DESC").
		Limit(100).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbChecks []store.DBHTTPCheck
	err = r.db.SelectContext(ctx, &dbChecks, query, args...)
	if err != nil {
		return nil, err
	}

	checks := make([]*store.HTTPCheck, len(dbChecks))
	for i, dbCheck := range dbChecks {
		checks[i] = dbCheck.ToDomain()
	}

	return checks, nil
}

// Get failing checks with more control
func (r *HTTPCheckRepository) GetFailingChecksWithOptions(ctx context.Context, options FailingChecksOptions) ([]*store.HTTPCheck, error) {
	queryBuilder := r.sb.Select("*").
		From("http_check_results").
		Where(sq.Or{
			sq.Eq{"status": "down"},
			sq.Eq{"status": "degraded"},
		}).
		Where(sq.GtOrEq{"checked_at": options.Since}).
		OrderBy("checked_at DESC")

	if options.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(options.Limit))
	}

	if options.EndpointName != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": options.EndpointName})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbChecks []store.DBHTTPCheck
	err = r.db.SelectContext(ctx, &dbChecks, query, args...)
	if err != nil {
		return nil, err
	}

	checks := make([]*store.HTTPCheck, len(dbChecks))
	for i, dbCheck := range dbChecks {
		checks[i] = dbCheck.ToDomain()
	}

	return checks, nil
}

// Cleanup old checks
func (r *HTTPCheckRepository) CleanupOldChecks(ctx context.Context, daysToKeep int) error {
	query, args, err := r.sb.Delete("").
		From("http_check_results").
		Where(sq.Lt{"checked_at": time.Now().Add(-time.Duration(daysToKeep) * 24 * time.Hour)}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// Cleanup with batch size to avoid large transactions
func (r *HTTPCheckRepository) CleanupOldChecksBatch(ctx context.Context, daysToKeep int, batchSize int) error {
	cutoff := time.Now().Add(-time.Duration(daysToKeep) * 24 * time.Hour)

	for {
		query, args, err := r.sb.Delete("").
			From("http_check_results").
			Where(sq.Lt{"checked_at": cutoff}).
			Limit(uint64(batchSize)).
			ToSql()

		if err != nil {
			return fmt.Errorf("failed to build delete query: %w", err)
		}

		result, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			break
		}
	}

	return nil
}

// Get checks with filters
func (r *HTTPCheckRepository) GetChecksWithFilters(ctx context.Context, filters CheckFilters) ([]*store.HTTPCheck, error) {
	queryBuilder := r.sb.Select("*").
		From("http_check_results").
		OrderBy("checked_at DESC")

	if filters.Name != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filters.Name})
	}

	if len(filters.Status) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"status": filters.Status})
	}

	if !filters.StartTime.IsZero() {
		queryBuilder = queryBuilder.Where(sq.GtOrEq{"checked_at": filters.StartTime})
	}

	if !filters.EndTime.IsZero() {
		queryBuilder = queryBuilder.Where(sq.LtOrEq{"checked_at": filters.EndTime})
	}

	if filters.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(filters.Limit))
	}

	if filters.Offset > 0 {
		queryBuilder = queryBuilder.Offset(uint64(filters.Offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbChecks []store.DBHTTPCheck
	err = r.db.SelectContext(ctx, &dbChecks, query, args...)
	if err != nil {
		return nil, err
	}

	checks := make([]*store.HTTPCheck, len(dbChecks))
	for i, dbCheck := range dbChecks {
		checks[i] = dbCheck.ToDomain()
	}

	return checks, nil
}

// Get endpoint names with optional filters
func (r *HTTPCheckRepository) GetEndpointNames(ctx context.Context, search string, limit int) ([]string, error) {
	queryBuilder := r.sb.Select("DISTINCT name").
		From("http_check_results").
		OrderBy("name")

	if search != "" {
		queryBuilder = queryBuilder.Where(sq.Like{"name": "%" + search + "%"})
	}

	if limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(limit))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var names []string
	err = r.db.SelectContext(ctx, &names, query, args...)
	return names, err
}

// Count checks with filters
func (r *HTTPCheckRepository) CountChecks(ctx context.Context, filters CheckFilters) (int, error) {
	queryBuilder := r.sb.Select("COUNT(*)").
		From("http_check_results")

	if filters.Name != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filters.Name})
	}

	if len(filters.Status) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"status": filters.Status})
	}

	if !filters.StartTime.IsZero() {
		queryBuilder = queryBuilder.Where(sq.GtOrEq{"checked_at": filters.StartTime})
	}

	if !filters.EndTime.IsZero() {
		queryBuilder = queryBuilder.Where(sq.LtOrEq{"checked_at": filters.EndTime})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	return count, err
}
