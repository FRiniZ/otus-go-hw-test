package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	api "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/api/stub"
	logger "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	internalgrpc "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/grpcservice"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func helperAPIEvent(id int64, userid int64, onTime, offTime time.Time) *api.Event {
	aEvent := api.Event{}
	aEvent.ID = &id
	aEvent.UserID = &userid
	aEvent.Title = func(s string) *string { return &s }(fmt.Sprintf("TitleN%v", userid))
	aEvent.Description = func(s string) *string { return &s }(fmt.Sprintf("DescriptionN%v", userid))
	aEvent.OnTime = timestamppb.New(onTime)
	aEvent.OffTime = timestamppb.New(offTime)
	aEvent.NotifyTime = timestamppb.New(time.Time{})
	return &aEvent
}

func TestGrpcService(t *testing.T) { //nolint:funlen
	currTime := time.Now()
	attempt := 10
	step := 100 // Must be more attempt

	require.Less(t, attempt, step)

	db := memorystorage.New()
	log, err := logger.New("DEBUG", os.Stdout)
	require.NoError(t, err)
	calendar := &Calendar{log: log, storage: db}

	dialer := func() func(context.Context, string) (net.Conn, error) {
		listener := bufconn.Listen(1024 * 1024)

		unarayLoggerEnricherIntercepter := func(ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (interface{}, error) { //nolint:gofumpt
			var b strings.Builder
			ip, _ := peer.FromContext(ctx)
			md, ok := metadata.FromIncomingContext(ctx)
			userAgent := "unknown"

			if ok {
				userAgent = md["user-agent"][0]
			}

			b.WriteString(ip.Addr.String())
			b.WriteString(" ")
			b.WriteString(time.Now().Format("02/Jan/2006:15:04:05 -0700"))
			b.WriteString(" ")
			b.WriteString(info.FullMethod)
			b.WriteString(" ")
			b.WriteString(userAgent)
			b.WriteString("\"\n")
			log.Infof(b.String())

			return handler(ctx, req)
		}

		server := grpc.NewServer(grpc.UnaryInterceptor(unarayLoggerEnricherIntercepter))
		grpcsrv := internalgrpc.New(log, calendar, internalgrpc.Conf{}, server)
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

	getConn := func(ctx context.Context) *grpc.ClientConn {
		conn, err := grpc.DialContext(ctx, "",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer()))
		if err != nil {
			t.Fatal(err)
		}
		return conn
	}

	t.Run("case_insert", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		client := api.NewCalendarClient(conn)
		for i := 0; i < attempt; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(int64(i+step), int64(i+step), currTime, currTime.AddDate(0, 0, 2)),
			}

			rep, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)
			require.NotZero(t, rep.ID)
		}
	})

	step += step
	t.Run("case_update", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		client := api.NewCalendarClient(conn)
		for i := 0; i < attempt; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, int64(i+step), currTime, currTime.AddDate(0, 0, 2)),
			}

			rep, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)

			event.Event.ID = rep.ID
			_, err = client.UpdateEvent(ctx, &event)
			require.NoError(t, err)
		}
	})

	step += step
	t.Run("case_delete", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		client := api.NewCalendarClient(conn)
		for i := 0; i < attempt; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, int64(i+step), currTime, currTime.AddDate(0, 0, 2)),
			}

			rep, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)

			_, err = client.DeleteEvent(ctx, &api.ReqByID{ID: rep.ID})
			require.NoError(t, err)
		}
	})

	step += step
	t.Run("case_lookup", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		client := api.NewCalendarClient(conn)
		for i := 0; i < attempt; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, int64(i+step), currTime, currTime.AddDate(0, 0, 2)),
			}

			rep, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)

			found, err := client.LookupEvent(ctx, &api.ReqByID{ID: rep.ID})
			require.NoError(t, err)
			require.Len(t, found.GetEvent(), 1)
			require.EqualValues(t, found.GetEvent()[0].GetID(), rep.GetID())
		}
	})

	step += step
	t.Run("case_listevents", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		userID := int64(step)
		client := api.NewCalendarClient(conn)
		for i := 0; i < attempt; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, userID, currTime.AddDate(0, 0, i*2), currTime.AddDate(0, 0, i*2+1)),
			}

			_, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)
		}

		founds, err := client.ListEvents(ctx, &api.ReqByUser{UserID: &userID})
		require.NoError(t, err)
		require.Len(t, founds.GetEvent(), attempt)
	})

	step += step
	t.Run("case_listevents_day", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		userID := int64(step)
		client := api.NewCalendarClient(conn)
		currTime2 := currTime
		for i := 0; i < attempt; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, userID, currTime2, currTime2.AddDate(0, 0, 1)),
			}
			currTime2 = currTime2.AddDate(0, 0, 2)

			_, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)
		}

		founds, err := client.ListEventsDay(ctx, &api.ReqByUserByDate{UserID: &userID, Date: timestamppb.New(currTime)})
		require.NoError(t, err)
		// Only one event with currTime
		require.Len(t, founds.GetEvent(), 1)
	})

	step += step
	t.Run("case_listevents_week", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		userID := int64(step)
		client := api.NewCalendarClient(conn)

		layout := "2006-01-02 15:04:05 -0700 MST"
		currTime, err := time.Parse(layout, "2023-01-02 00:00:01 -0700 MST")
		require.NoError(t, err)
		currTime2 := currTime

		for i := 0; i < attempt+7; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, userID, currTime2, currTime2.Add(86399*time.Second)),
			}
			currTime2 = currTime2.Add(86400 * time.Second)

			_, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)
		}

		founds, err := client.ListEventsWeek(ctx, &api.ReqByUserByDate{UserID: &userID, Date: timestamppb.New(currTime)})
		require.NoError(t, err)
		// Should be 7 events in week
		require.Len(t, founds.GetEvent(), 7)
	})

	step += step
	t.Run("case_listevents_month", func(t *testing.T) {
		step := step
		t.Parallel()
		ctx := context.Background()
		conn := getConn(ctx)
		defer conn.Close()

		userID := int64(step)
		client := api.NewCalendarClient(conn)

		layout := "2006-01-02 15:04:05 -0700 MST"
		currTime, err := time.Parse(layout, "2023-01-01 00:00:01 -0700 MST")
		require.NoError(t, err)
		currTime2 := currTime

		for i := 0; i < attempt+31; i++ {
			event := api.ReqByEvent{
				Event: helperAPIEvent(0, userID, currTime2, currTime2.Add(86399*time.Second)),
			}
			currTime2 = currTime2.Add(86400 * time.Second)

			_, err := client.InsertEvent(ctx, &event)
			require.NoError(t, err)
		}

		founds, err := client.ListEventsMonth(ctx, &api.ReqByUserByDate{UserID: &userID, Date: timestamppb.New(currTime)})
		require.NoError(t, err)
		// Should be 31 events in January 2023
		require.Len(t, founds.GetEvent(), 31)
	})
}