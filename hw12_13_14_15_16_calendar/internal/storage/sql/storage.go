package sqlstorage

import (
	"context"
	"database/sql"

	// Use pgx driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Write(ctx context.Context, event storage.Event) error {
	query := `INSERT INTO events(id, title, start_time, end_time, description, user_id) 
	values($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx,
		query,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Update(ctx context.Context, event storage.Event) error {
	query := `UPDATE events
	SET (title, start_time, end_time, description, user_id) 
	values($1, $2, $3, $4, $5)
	WHERE events.id = $6`
	_, err := s.db.ExecContext(ctx,
		query,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, event storage.Event) error {
	query := `DELETE FROM events
	WHERE events.id = $1`
	_, err := s.db.ExecContext(ctx, query, event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) EventsByDay(ctx context.Context, date string) ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, start_time, end_time, description, user_id
	FROM events
	WHERE start_time::date = $1`
	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
			&event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Storage) EventsByWeek(ctx context.Context, date string) ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, start_time, end_time, description, user_id
	FROM events
	WHERE start_time BETWEEN $1::date and $1::date + interval '7 day'`
	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
			&event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Storage) EventsByMonth(ctx context.Context, date string) ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, start_time, end_time, description, user_id
	FROM events
	WHERE start_time BETWEEN $1::date and $1::date + interval '1 month'`
	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
			&event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
