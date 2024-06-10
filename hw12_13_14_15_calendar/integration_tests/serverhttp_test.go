package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/logger"
	serverhttp "github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func userCase(ctx context.Context, t *testing.T, mutex *sync.Mutex, address string) {
	t.Helper()
	mutex.Lock()
	defer mutex.Unlock()

	// Create user
	userData := `{"name": "testuser"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address+"/create/user",
		bytes.NewBuffer([]byte(userData)))
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

	eventData := `{"title": "testevent", "date": "2023-01-01"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address+"/create/event",
		bytes.NewBuffer([]byte(eventData)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

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
	application := app.New(memoryStorage, *log)

	host := "localhost"
	port := "8080"
	servConfig := config.ServerConfig{
		Host: host,
		Port: port,
	}

	serv := serverhttp.NewServer(log, application, &servConfig)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_ = serv.Start()
	}()

	mutex := sync.Mutex{}
	address := "http://" + net.JoinHostPort(host, port)
	ctx := context.Background()

	time.Sleep(1 * time.Second)
	userCase(ctx, t, &mutex, address)
	time.Sleep(1 * time.Second)
	eventCase(ctx, t, &mutex, address)

	err := serv.Stop(ctx)
	require.Nil(t, err)

	wg.Wait()
}

func eventForPeriodCase(ctx context.Context, t *testing.T, mutex *sync.Mutex, address, endpoint string,
	queryParam string,
) {
	t.Helper()
	mutex.Lock()
	defer mutex.Unlock()

	eventData := `{"title": "testevent", "description": "this is a test event", 
"beginning": "2024-06-09T12:00:00Z", "finish": "2024-06-09T14:00:00Z", "notification": "2024-06-09T11:00:00Z", 
"userid": "1"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address+"/create/event",
		bytes.NewBuffer([]byte(eventData)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, address+endpoint+"?"+queryParam, nil)
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var events []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&events)
	resp.Body.Close()
	require.Nil(t, err)
	require.NotEmpty(t, events)

	eventID := events[0]["id"].(string)
	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, address+"/delete/event/"+eventID, nil)
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestServerTwo(t *testing.T) {
	logConfig := config.LoggerConfig{
		Level: "info",
	}

	log := logger.New(&logConfig)
	memoryStorage := memorystorage.New()
	application := app.New(memoryStorage, *log)

	host := "localhost"
	port := "8080"
	servConfig := config.ServerConfig{
		Host: host,
		Port: port,
	}

	serv := serverhttp.NewServer(log, application, &servConfig)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_ = serv.Start()
	}()

	mutex := sync.Mutex{}
	address := "http://" + net.JoinHostPort(host, port)
	ctx := context.Background()

	time.Sleep(1 * time.Second)
	userCase(ctx, t, &mutex, address)
	time.Sleep(1 * time.Second)
	eventCase(ctx, t, &mutex, address)

	time.Sleep(1 * time.Second)
	eventForPeriodCase(ctx, t, &mutex, address, "/select/events/day", "date=2024-06-09")
	time.Sleep(1 * time.Second)
	eventForPeriodCase(ctx, t, &mutex, address, "/select/events/week", "startDate=2024-06-09")
	time.Sleep(1 * time.Second)
	eventForPeriodCase(ctx, t, &mutex, address, "/select/events/month", "startDate=2024-06-01")

	err := serv.Stop(ctx)
	require.Nil(t, err)

	wg.Wait()
}
