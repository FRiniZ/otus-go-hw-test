package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" // needs for init
)

type Storage struct {
	dsn string
	db  *sql.DB
}

type EventSQL struct {
	ID          sql.NullInt64
	UserID      sql.NullInt64
	Title       sql.NullString
	Description sql.NullString
	OnTime      sql.NullTime
	OffTime     sql.NullTime
	NotifyTime  sql.NullTime
}

func GetEvent(e EventSQL) (event storage.Event) {
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

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	err = s.db.PingContext(ctx)

	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	return err
}

func (s *Storage) Close(ctx context.Context) error {
	s.db.Close()
	ctx.Done()
	return nil
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

func (s *Storage) InsertEvent(e *storage.Event) error {
	query := `INSERT INTO events (userid, title, description, ontime)
						  values ($1, $2, $3, $4) RETURNING id`

	rows, err := s.db.Query(query, e.UserID, stringValue(e.Title),
		stringValue(e.Description), timeValue(e.OnTime))
	if err != nil {
		return fmt.Errorf("failed insert event: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&e.ID); err != nil {
			return fmt.Errorf("failed rows.Scan: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed rows.Next: %w", err)
	}

	return nil
}

func (s *Storage) UpdateEvent(e *storage.Event) error {
	query := `UPDATE events SET userid = $2,
								title = $3,
								description = $4,
								ontime = $5,
								offtime = $6,
								notifytime = $7
	          WHERE id = $1`

	res, err := s.db.Exec(query, e.ID, e.UserID, e.Title, e.Description,
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

func (s *Storage) DeleteEvent(e *storage.Event) error {
	query := `DELETE FROM events
	          WHERE id = $1`

	if _, err := s.db.Exec(query, e.ID); err != nil {
		return fmt.Errorf("failed delete event: %w", err)
	}

	return nil
}

func (s *Storage) ListEvents(userID int64) (events []storage.Event, err error) {
	var e storage.Event
	var eSQL EventSQL

	query := `SELECT id, userid, title, description, ontime, offtime, notifytime
	          FROM events
			  WHERE userid = $1`

	rows, err := s.db.Query(query, userID)
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

func (s *Storage) LookupEvent(eID int64) (e storage.Event, err error) {
	var eSQL EventSQL
	query := `SELECT id, userid, title, description, ontime, offtime, notifytime
	          FROM events
			  WHERE id = $1`

	rows, err := s.db.Query(query, eID)
	if err != nil {
		return e, fmt.Errorf("failed lookup event: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&eSQL.ID, &eSQL.UserID, &eSQL.Title, &eSQL.Description,
			&eSQL.OnTime, &eSQL.OffTime, &eSQL.NotifyTime); err != nil {
			return e, fmt.Errorf("failed rows.Scan: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return e, fmt.Errorf("failed rows.Next: %w", err)
	}

	e = GetEvent(eSQL)

	return e, err
}
