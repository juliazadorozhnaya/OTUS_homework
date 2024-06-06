package servergrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/app"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/logger"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/grpc/api"
	memorystorage "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(_ context.Context, _ string) (net.Conn, error) {
	return lis.Dial()
}

func startTestGRPCServer(t *testing.T, logger server.Logger, application *app.Calendar) *grpc.Server {
	t.Helper()
	grpcServer := grpc.NewServer()
	lis = bufconn.Listen(bufSize)

	eventServer := api.NewEventServer(logger, application)
	api.RegisterEventServiceServer(grpcServer, eventServer)

	userServer := api.NewUserServer(logger, application)
	api.RegisterUserServiceServer(grpcServer, userServer)

	errChan := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Server exited with error: %v", err)
		}
	default:
	}

	return grpcServer
}

func TestGRPCServer(t *testing.T) {
	logConfig := config.LoggerConfig{
		Level: "info",
	}

	log := logger.New(&logConfig)
	memoryStorage := memorystorage.New()
	application := app.New(memoryStorage)

	grpcServer := startTestGRPCServer(t, log, application)

	ctx := context.Background()

	//nolint: staticcheck
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := api.NewEventServiceClient(conn)
	userClient := api.NewUserServiceClient(conn)

	t.Run("UserCase", func(t *testing.T) {
		userCase(ctx, t, userClient)
	})

	t.Run("EventCase", func(t *testing.T) {
		eventCase(ctx, t, client, userClient)
	})

	grpcServer.GracefulStop()
}

func userCase(ctx context.Context, t *testing.T, client api.UserServiceClient) {
	t.Helper()
	// Create user
	user := &api.User{
		FirstName: "testuser",
		LastName:  "last",
		Email:     "test@test.com",
		Age:       30,
	}

	_, err := client.CreateUser(ctx, user)
	require.NoError(t, err)

	// Select users
	stream, err := client.SelectUsers(ctx, &api.Void{})
	require.NoError(t, err)

	users := make([]*api.User, 0)
	for {
		user, err := stream.Recv()
		if err != nil {
			break
		}
		users = append(users, user)
	}

	require.NotEmpty(t, users)

	// Delete user
	userID := &api.User{ID: users[0].ID}
	_, err = client.DeleteUser(ctx, userID)
	require.NoError(t, err)
}

func eventCase(ctx context.Context, t *testing.T, client api.EventServiceClient, userClient api.UserServiceClient) {
	t.Helper()
	// Create user for event
	user := &api.User{
		FirstName: "testuser",
		LastName:  "last",
		Email:     "test@test.com",
		Age:       30,
	}

	_, err := userClient.CreateUser(ctx, user)
	require.NoError(t, err)

	stream, err := userClient.SelectUsers(ctx, &api.Void{})
	require.NoError(t, err)

	users := make([]*api.User, 0)
	for {
		user, err := stream.Recv()
		if err != nil {
			break
		}
		users = append(users, user)
	}

	require.NotEmpty(t, users)
	userID := users[0].ID

	// Create event
	event := &api.Event{
		Title:         "testevent",
		Description:   "desc",
		BeginningT:    timestamppb.New(time.Now()),
		FinishT:       timestamppb.New(time.Now().Add(time.Hour)),
		NotificationT: timestamppb.New(time.Now().Add(30 * time.Minute)),
		UserID:        userID,
	}

	_, err = client.CreateEvent(ctx, event)
	require.NoError(t, err)

	// Select events
	streamEvents, err := client.SelectEvents(ctx, &api.Void{})
	require.NoError(t, err)

	events := make([]*api.Event, 0)
	for {
		event, err := streamEvents.Recv()
		if err != nil {
			break
		}
		events = append(events, event)
	}

	require.NotEmpty(t, events)

	// Delete event
	eventID := &api.Event{ID: events[0].ID}
	_, err = client.DeleteEvent(ctx, eventID)
	require.NoError(t, err)
}
