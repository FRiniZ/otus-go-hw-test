package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestSqlStorage(t *testing.T) {
	event := storage.Event{
		ID:          0,
		UserID:      1,
		Title:       "TitleN1",
		Description: "DescriptionN1",
		OnTime:      time.Now(),
		OffTime:     time.Now().AddDate(0, 0, 7),
		NotifyTime:  time.Time{},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := Storage{dsn: "", db: db}

	t.Run("case_insert", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO events (userid, title, description, ontime, offtime, notifytime)
		                              values ($1, $2, $3, $4, $5, $6) RETURNING id`).
			WithArgs(event.UserID, event.Title, event.Description,
				timeValue(event.OnTime), timeValue(event.OffTime), timeValue(event.NotifyTime)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

		err = storage.InsertEvent(context.Background(), &event)
		require.NoError(t, err)
		require.EqualValues(t, event.ID, int64(1))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case_update", func(t *testing.T) {
		event.UserID = 400
		mock.ExpectExec(`UPDATE events
						 SET userid = $2,
						 	 title = $3,
							 description = $4,
							 ontime = $5,
							 offtime = $6,
							 notifytime = $7
						WHERE id = $1`).
			WithArgs(event.ID, event.UserID, event.Title, event.Description,
				timeValue(event.OnTime), timeValue(event.OffTime), timeValue(event.NotifyTime)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = storage.UpdateEvent(context.Background(), &event)
		require.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case_delete", func(t *testing.T) {
		event.ID = 100
		mock.ExpectExec("DELETE FROM events WHERE id = $1").
			WithArgs(event.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = storage.DeleteEvent(context.Background(), &event)
		require.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case_lookup", func(t *testing.T) {
		eID := int64(100)
		userID := int64(200)
		mock.ExpectQuery(`SELECT id, userid, title, description, ontime, offtime, notifytime
						  FROM events WHERE id = $1`).
			WithArgs(eID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userid", "title", "description", "ontime", "offtime", "notifytime"}).
				AddRow(eID, userID, "TitleN100", "DescriptionN100",
					timeValue(time.Now()), timeValue(time.Now().AddDate(0, 0, 7)), timeValue(time.Time{})))

		eFound, err := storage.LookupEvent(context.Background(), eID)
		require.NoError(t, err)
		require.EqualValues(t, userID, eFound.UserID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case_listevents", func(t *testing.T) {
		eID1 := int64(100)
		eID2 := int64(100)
		userID := int64(200)
		mock.ExpectQuery("SELECT id, userid, title, description, ontime, offtime, notifytime FROM events WHERE userid = $1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userid", "title", "description", "ontime", "offtime", "notifytime"}).
				AddRow(eID1, userID, "TitleN100", "DescriptionN100",
					timeValue(time.Now()), timeValue(time.Now().AddDate(0, 0, 7)), timeValue(time.Time{})).
				AddRow(eID2, userID, "TitleN101", "DescriptionN101",
					timeValue(time.Now()), timeValue(time.Now().AddDate(0, 0, 7)), timeValue(time.Time{})))

		eFound, err := storage.ListEvents(context.Background(), userID)
		require.NoError(t, err)
		require.EqualValues(t, 2, len(eFound))
		require.EqualValues(t, eID1, eFound[0].ID)
		require.EqualValues(t, eID2, eFound[1].ID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
