package serverhttp

import (
	"bytes"
	"context"
	"encoding/json"
	memorystorage "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/storage/memory"
	"net"
	"net/http"
	"sync"
	"testing"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/app"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
)

func userCase(ctx context.Context, t *testing.T, mutex *sync.Mutex, address string) {
	t.Helper()
	mutex.Lock()
	defer mutex.Unlock()

	// Create user
	userData := `{"name": "testuser"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address+"/create/user", bytes.NewBuffer([]byte(userData)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Select users
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, address+"/select/users", nil)
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var users []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	resp.Body.Close()
	require.Nil(t, err)
	require.NotEmpty(t, users)

	// Delete user
	userID := users[0]["id"].(string)
	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, address+"/delete/user/"+userID, nil)
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func eventCase(ctx context.Context, t *testing.T, mutex *sync.Mutex, address string) {
	t.Helper()
	mutex.Lock()
	defer mutex.Unlock()

	// Create event
	eventData := `{"title": "testevent", "date": "2023-01-01"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address+"/create/event",
		bytes.NewBuffer([]byte(eventData)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Select events
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, address+"/select/events", nil)
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var events []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&events)
	resp.Body.Close()
	require.Nil(t, err)
	require.NotEmpty(t, events)

	// Delete event
	eventID := events[0]["id"].(string)
	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, address+"/delete/event/"+eventID, nil)
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestServer(t *testing.T) {
	logConfig := config.LoggerConfig{
		Level: "info",
	}

	log := logger.New(&logConfig)
	memoryStorage := memorystorage.New()
	application := app.New(memoryStorage)

	host := "localhost"
	port := "8080"
	servConfig := config.ServerConfig{
		Host: host,
		Port: port,
	}

	serv := NewServer(log, application, &servConfig)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_ = serv.Start()
	}()

	mutex := sync.Mutex{}
	address := "http://" + net.JoinHostPort(host, port)
	ctx := context.Background()

	userCase(ctx, t, &mutex, address)
	eventCase(ctx, t, &mutex, address)

	err := serv.Stop(ctx)
	require.Nil(t, err)

	wg.Wait()
}
