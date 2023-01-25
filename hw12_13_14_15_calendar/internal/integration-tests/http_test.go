package integrationtests

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

const (
	msgInserted string = "Inserted"
	msgUpdated  string = "Updated"
	msgDeleted  string = "Deleted"
)

type ReplayMsg struct {
	Msg string `json:"msg"`
	Err string `json:"error"`
}

func helperDecode(stream io.Reader, r interface{}) error {
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(r); err != nil {
		return err
	}
	return nil
}

func TestIntegrationHTTPAPi(t *testing.T) { //nolint:nolintlint
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	httpHost := "calendar"
	if host, ok := os.LookupEnv("CALENDAR_HOST"); ok {
		httpHost = host
	}

	httpPort := "8089"

	userID400 := int64(400)
	bodyUserID := fmt.Sprintf(`{"userid": %d}`, userID400)

	body1 := `{
		"id": 1,
		"userid": 200,
		"title" : "Title_N200",
		"description" : "Description_N200",
		"ontime" : "2015-09-18T00:00:00Z",
		"offtime" : "2015-09-19T00:00:00Z",
		"notifytime" : "0001-01-01T00:00:00Z"
	}`

	body2 := fmt.Sprintf(`{
	 	"id": 1,
	 	"userid": %d,
	 	"title" : "Title_N400",
	 	"description" : "Description_N400",
	 	"ontime" : "2015-09-18T00:00:00Z",
	 	"offtime" : "2015-09-20T00:00:00Z",
	 	"notifytime" : "0001-01-01T00:00:00Z"
	 	}`, userID400)

	body3 := fmt.Sprintf(`{
					"id": 2,
					"userid": %d,
					"title" : "Title_N402",
					"description" : "Description_N402",
					"ontime" : "2015-10-18T00:00:00Z",
					"offtime" : "2015-10-20T00:00:00Z",
					"notifytime" : "0001-01-01T00:00:00Z"
				}`, userID400)

	t.Run("case_insert", func(t *testing.T) {
		var rep ReplayMsg
		httpcli := &http.Client{}

		reader := strings.NewReader(body1)
		url := "http://" + httpHost + ":" + httpPort + "/InsertEvent"

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, reader)
		require.NoError(t, err)
		res, err := httpcli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		err = helperDecode(res.Body, &rep)
		require.NoError(t, err)

		require.Empty(t, rep.Err)
		require.EqualValues(t, msgInserted, rep.Msg)

		reader = strings.NewReader(body3)
		req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, url, reader)
		require.NoError(t, err)
		res, err = httpcli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		err = helperDecode(res.Body, &rep)
		require.NoError(t, err)

		require.Empty(t, rep.Err)
		require.EqualValues(t, msgInserted, rep.Msg)
	})

	t.Run("case_update", func(t *testing.T) {
		var rep ReplayMsg
		httpcli := &http.Client{}
		url := "http://" + httpHost + ":" + httpPort + "/UpdateEvent"

		reader := strings.NewReader(body2)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, reader)
		require.NoError(t, err)
		res, err := httpcli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		err = helperDecode(res.Body, &rep)
		require.NoError(t, err)

		require.Empty(t, rep.Err)
		require.EqualValues(t, msgUpdated, rep.Msg)
	})

	t.Run("case_lookup", func(t *testing.T) {
		var rep model.Event
		httpcli := &http.Client{}
		url := "http://" + httpHost + ":" + httpPort + "/LookupEvent"

		reader := strings.NewReader(body1)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, reader)
		require.NoError(t, err)
		res, err := httpcli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		err = helperDecode(res.Body, &rep)
		require.NoError(t, err)
		require.EqualValues(t, userID400, rep.UserID)
	})

	t.Run("case_listevents", func(t *testing.T) {
		var rep []model.Event
		httpcli := &http.Client{}
		url := "http://" + httpHost + ":" + httpPort + "/ListEvents"

		reader := strings.NewReader(bodyUserID)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, reader)
		require.NoError(t, err)
		res, err := httpcli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		err = helperDecode(res.Body, &rep)
		require.NoError(t, err)
		require.EqualValues(t, int(2), len(rep))
		require.EqualValues(t, userID400, rep[0].UserID)
		require.EqualValues(t, userID400, rep[1].UserID)
	})

	t.Run("case_delete", func(t *testing.T) {
		var rep ReplayMsg
		httpcli := &http.Client{}
		url := "http://" + httpHost + ":" + httpPort + "/DeleteEvent"

		reader := strings.NewReader(body2)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, reader)
		require.NoError(t, err)
		res, err := httpcli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		err = helperDecode(res.Body, &rep)
		require.NoError(t, err)

		require.Empty(t, rep.Err)
		require.EqualValues(t, msgDeleted, rep.Msg)
	})
}
