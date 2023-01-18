package grpcservice

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	api "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/api/stub"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func helperAPIEvent(id int64, userid int64, currTime time.Time) *api.Event {
	aEvent := api.Event{}
	aEvent.ID = &id
	aEvent.UserID = &userid
	aEvent.Title = func(s string) *string { return &s }(fmt.Sprintf("TitleN%v", userid))
	aEvent.Description = func(s string) *string { return &s }(fmt.Sprintf("DescriptionN%v", userid))
	aEvent.OnTime = timestamppb.New(currTime)
	aEvent.OffTime = timestamppb.New(currTime.AddDate(0, 0, 7))
	aEvent.NotifyTime = timestamppb.New(time.Time{})
	return &aEvent
}

func TestGrpcService(t *testing.T) {
	dialer := func() func(context.Context, string) (net.Conn, error) {
		listener := bufconn.Listen(1024 * 1024)

		db := memorystorage.New()
		log, err := logger.New("DEBUG", os.Stdout)
		require.NoError(t, err)
		calendar := app.New(log, db)
		conf := Conf{}

		server := grpc.NewServer(grpc.UnaryInterceptor(UnaryLoggerEnricherInterceptor))
		grpcsrv := New(log, calendar, conf, server)
		api.RegisterCalendarServer(server, grpcsrv)

		go func() {
			if err := server.Serve(listener); err != nil {
				require.NoError(t, err)
			}
		}()

		return func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := api.NewCalendarClient(conn)

	currTime := time.Now()

	tests := []struct {
		name  string
		call  func(context.Context, *api.RequestV1, ...grpc.CallOption) (*api.ReplyV1, error)
		err   error
		reply []*api.Event
		event api.Event
	}{
		{
			"case_insert",
			client.InsertEventV1,
			nil,
			[]*api.Event{
				helperAPIEvent(1, 1, currTime),
			},
			*helperAPIEvent(0, 1, currTime),
		},
		{
			"case_update",
			client.UpdateEventV1,
			nil,
			[]*api.Event{
				helperAPIEvent(1, 2, currTime),
			},
			*helperAPIEvent(1, 2, currTime),
		},
		{
			"case_lookup",
			client.LookupEventV1,
			nil,
			[]*api.Event{
				helperAPIEvent(1, 2, currTime),
			},
			*helperAPIEvent(1, 2, currTime),
		},
		{
			"case_listevents",
			client.ListEventsV1,
			nil,
			[]*api.Event{
				helperAPIEvent(1, 2, currTime),
			},
			*helperAPIEvent(1, 2, currTime),
		},
	}

	for i := 0; i < len(tests); i++ {
		tt := &tests[i]
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			req := &api.RequestV1{
				Event: &tt.event,
			}
			rep, err := tt.call(ctx, req)
			require.NoError(t, err)
			require.Equal(t, len(tt.reply), len(rep.Event))
			for i := 0; i < len(tt.reply); i++ {
				require.EqualValues(t, *tt.reply[i].ID, *rep.Event[i].ID)
				require.EqualValues(t, *tt.reply[i].UserID, *rep.Event[i].UserID)
				require.EqualValues(t, *tt.reply[i].Title, *rep.Event[i].Title)
				require.EqualValues(t, *tt.reply[i].Description, *rep.Event[i].Description)
				require.EqualValues(t, tt.reply[i].OnTime.GetSeconds(), rep.Event[i].OnTime.GetSeconds())
				require.EqualValues(t, tt.reply[i].OffTime.GetSeconds(), rep.Event[i].OffTime.GetSeconds())
				require.EqualValues(t, tt.reply[i].NotifyTime.GetSeconds(), rep.Event[i].NotifyTime.GetSeconds())
			}
		})
	}
}
