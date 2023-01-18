package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" // needs for init
)

type Storage struct {
	dsn string
	db  *sql.DB
}

var ErrEventNotFound = errors.New("event not found")

type EventDTO struct {
	ID          sql.NullInt64
	UserID      sql.NullInt64
	Title       sql.NullString
	Description sql.NullString
	OnTime      sql.NullTime
	OffTime     sql.NullTime
	NotifyTime  sql.NullTime
}

func GetEvent(e EventDTO) (event storage.Event) {
	if e.ID.Valid {
		event.ID = e.ID.Int64
	}

	if e.UserID.Valid {
		event.UserID = e.UserID.Int64
	}

	if e.Title.Valid {
		event.Title = e.Title.String
	}

	if e.Description.Valid {
		event.Description = e.Description.String
	}

	if e.OnTime.Valid {
		event.OnTime = e.OnTime.Time
	}

	if e.OffTime.Valid {
		event.OffTime = e.OffTime.Time
	}

	if e.NotifyTime.Valid {
		event.NotifyTime = e.NotifyTime.Time
	}
	return event
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error

	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("failed connect to db: %w", err)
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.db.PingContext(queryCtx)

	if err != nil {
		return fmt.Errorf("failed connect to db: %w", err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	/*
		Сохранено временно. Возможно нужно как то доработать. Пока не знаю как
		ctxClose, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dbCloseChan := make(chan error)
		go func(che chan error) {
			if err := s.db.Close(); err != nil {
				che <- err
				return
			}
			che <- nil
		}(dbCloseChan)

		select {
		case <-ctxClose.Done():
			return ctxClose.Err()
		case err, ok := <-dbCloseChan:
			if ok && err != nil {
				fmt.Println("CloseChan")
				return err
			}
		}
	*/
	return s.db.Close()
}

func timeValue(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func stringValue(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func (s *Storage) InsertEvent(ctx context.Context, e *storage.Event) error {
	query := `INSERT INTO events (userid, title, description, ontime, offtime, notifytime)
						  values ($1, $2, $3, $4, $5, $6) RETURNING id`

	row := s.db.QueryRowContext(ctx, query, e.UserID, stringValue(e.Title),
		stringValue(e.Description), timeValue(e.OnTime), timeValue(e.OffTime),
		timeValue(e.NotifyTime))

	if err := row.Scan(&e.ID); err != nil {
		return fmt.Errorf("failed rows.Scan: %w", err)
	}

	if err := row.Err(); err != nil {
		return fmt.Errorf("failed rows.Next: %w", err)
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, e *storage.Event) error {
	query := `UPDATE events SET userid = $2,
								title = $3,
								description = $4,
								ontime = $5,
								offtime = $6,
								notifytime = $7
	          WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, e.ID, e.UserID, e.Title, e.Description,
		timeValue(e.OnTime),
		timeValue(e.OffTime),
		timeValue(e.NotifyTime))
	if err != nil {
		return fmt.Errorf("failed update event: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed get RowsAffected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("failed rowsAffected: %v", rowsAffected)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, e *storage.Event) error {
	query := `DELETE FROM events
	          WHERE id = $1`

	if _, err := s.db.ExecContext(ctx, query, e.ID); err != nil {
		return fmt.Errorf("failed delete event: %w", err)
	}

	return nil
}

func (s *Storage) ListEvents(ctx context.Context, userID int64) (events []storage.Event, err error) {
	var e storage.Event
	var eSQL EventDTO

	query := `SELECT id, userid, title, description, ontime, offtime, notifytime
	          FROM events
			  WHERE userid = $1`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return events, fmt.Errorf("failed lookup event: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&eSQL.ID, &eSQL.UserID, &eSQL.Title, &eSQL.Description,
			&eSQL.OnTime, &eSQL.OffTime, &eSQL.NotifyTime); err != nil {
			return events, fmt.Errorf("failed rows.Scan: %w", err)
		}
		e = GetEvent(eSQL)
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("failed lookup event: %w", err)
	}

	return events, err
}

func (s *Storage) firstDayOfWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func (s *Storage) firstDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

func (s *Storage) ListEventsWeek(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	return s.ListEventsDay(ctx, userID, s.firstDayOfWeek(date))
}

func (s *Storage) ListEventsMonth(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	return s.ListEventsDay(ctx, userID, s.firstDayOfMonth(date))
}

func (s *Storage) ListEventsDay(ctx context.Context, userID int64, date time.Time) (events []storage.Event, err error) {
	var e storage.Event
	var eSQL EventDTO

	query := `SELECT id, userid, title, description, ontime, offtime, notifytime
	          FROM events
			  WHERE userid = $1 && date($2) BETWEEN date(ontime) AND date(offtime)`

	rows, err := s.db.QueryContext(ctx, query, userID, date)
	if err != nil {
		return events, fmt.Errorf("failed lookup event: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&eSQL.ID, &eSQL.UserID, &eSQL.Title, &eSQL.Description,
			&eSQL.OnTime, &eSQL.OffTime, &eSQL.NotifyTime); err != nil {
			return events, fmt.Errorf("failed rows.Scan: %w", err)
		}
		e = GetEvent(eSQL)
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("failed lookup event: %w", err)
	}

	return events, err
}

func (s *Storage) LookupEvent(ctx context.Context, eID int64) (e storage.Event, err error) {
	var eSQL EventDTO
	query := `SELECT id, userid, title, description, ontime, offtime, notifytime
	          FROM events
			  WHERE id = $1`

	rows := s.db.QueryRowContext(ctx, query, eID)

	if err := rows.Scan(&eSQL.ID, &eSQL.UserID, &eSQL.Title, &eSQL.Description,
		&eSQL.OnTime, &eSQL.OffTime, &eSQL.NotifyTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, ErrEventNotFound
		}
		return e, fmt.Errorf("failed rows.Scan: %w", err)
	}

	if err := rows.Err(); err != nil {
		return e, fmt.Errorf("failed rows.Next: %w", err)
	}

	e = GetEvent(eSQL)

	return e, err
}

func (s *Storage) IsBusyDateTimeRange(ctx context.Context, userID int64, onTime, offTime time.Time) (bool, error) {
	var eSQL EventDTO
	query := `SELECT id
	          FROM events
			  WHERE userid = $1 AND
			  (($2 >= event.ontime && $2 <= event.offtime) OR
			   ($3 >= event.ontime && $3 <= event.offtime))`

	rows := s.db.QueryRowContext(ctx, query, userID, onTime, offTime)

	if err := rows.Scan(&eSQL.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed rows.Scan: %w", err)
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("failed rows.Next: %w", err)
	}

	return true, nil
}
