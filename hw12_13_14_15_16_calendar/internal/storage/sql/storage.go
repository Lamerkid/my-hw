package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	// Use pgx driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
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

func (s *Storage) Write(ctx context.Context, event domain.Event) error {
	query := `INSERT INTO events(id, title, description, start_time, end_time, user_id, notify_before, notified) 
	values($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.db.ExecContext(ctx,
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UserID,
		event.NotifyBefore,
		event.Notified)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, event domain.Event) error {
	query := `UPDATE events
	SET (title, description, start_time, end_time, user_id, notify_before, notified)
	values($2, $3, $4, $5, $6, $7, $8)
	WHERE events.id = $1`

	_, err := s.db.ExecContext(ctx,
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UserID,
		event.NotifyBefore,
		event.Notified)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events
	WHERE events.id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error) {
	query := `SELECT id, title, description, start_time, end_time, user_id, notify_before, notified
	FROM events
	WHERE events.id = $1`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return domain.Event{}, err
	}
	defer rows.Close()

	var event domain.Event

	err = rows.Scan(&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.UserID,
		&event.NotifyBefore,
		&event.Notified)
	if err != nil {
		return domain.Event{}, err
	}

	if err := rows.Err(); err != nil {
		return domain.Event{}, err
	}

	return event, nil
}

func (s *Storage) EventsByDay(ctx context.Context, date string) ([]domain.Event, error) {
	var events []domain.Event

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_before, notified
	FROM events
	WHERE start_time::date = $1`

	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.UserID,
			&event.NotifyBefore,
			&event.Notified)
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

func (s *Storage) EventsByWeek(ctx context.Context, date string) ([]domain.Event, error) {
	var events []domain.Event

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_before, notified
	FROM events
	WHERE start_time BETWEEN $1::date AND $1::date + interval '7 day'`

	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.UserID,
			&event.NotifyBefore,
			&event.Notified)
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

func (s *Storage) EventsByMonth(ctx context.Context, date string) ([]domain.Event, error) {
	var events []domain.Event

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_before, notified
	FROM events
	WHERE start_time BETWEEN $1::date AND $1::date + interval '1 month'`

	rows, err := s.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.UserID,
			&event.NotifyBefore,
			&event.Notified)
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

func (s *Storage) EventsForNotification(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	now := time.Now()

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_before, notified
	FROM events
	WHERE start_time > $1
	AND start_time - (notify_before / 1000 * interval '1 microsecond') <= $1
	AND notified IS FALSE`

	rows, err := s.db.QueryContext(ctx, query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.UserID,
			&event.NotifyBefore,
			&event.Notified)
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

func (s *Storage) MarkNotified(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE events
	SET notified = true
	WHERE events.id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
