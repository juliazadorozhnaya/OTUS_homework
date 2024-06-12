package serverhttp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/model"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/server"
)

const (
	selectEventsForDayMsg   = "selectEventsForDay: "
	selectEventsForWeekMsg  = "selectEventsForWeek: "
	selectEventsForMonthMsg = "selectEventsForMonth: "
)

type handler struct {
	logger server.Logger
	app    server.Application
}

// newHandler создает новый HTTP хендлер с логгером и приложением.
func newHandler(logger server.Logger, app server.Application) *handler {
	return &handler{
		logger: logger,
		app:    app,
	}
}

// createUser обрабатывает запрос на создание нового пользователя.
func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := readUserFromBody(r)
	if err != nil {
		h.logger.Error("createUser: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("Attempting to create user: " + user.Email)
	if err := h.app.CreateUser(ctx, user); err != nil {
		h.logger.Error("createUser: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("User created: " + user.Email)
	w.WriteHeader(http.StatusOK)
}

// selectUsers обрабатывает запрос на получение списка всех пользователей.
func (h *handler) selectUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug("Selecting users")
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectUsers(ctx)
	})
	if err != nil {
		h.logger.Error("selectUsers: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error("selectUsers: " + err.Error())
	}
	h.logger.Info("Users selected")
}

// deleteUser обрабатывает запрос на удаление пользователя по его ID.
func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getIDFromPath(r.URL.Path)
	if userID == "" {
		h.logger.Error("deleteUser: missing user ID in path")
		http.Error(w, "missing user ID", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Attempting to delete user: " + userID)
	if err := h.app.DeleteUser(ctx, userID); err != nil {
		h.logger.Error("deleteUser: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("User deleted: " + userID)
	w.WriteHeader(http.StatusOK)
}

// createEvent обрабатывает запрос на создание нового события.
func (h *handler) createEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error("createEvent: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("Attempting to create event: " + event.Title)
	if err := h.app.CreateEvent(ctx, event); err != nil {
		h.logger.Error("createEvent: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Event created: " + event.Title)
	w.WriteHeader(http.StatusOK)
}

// selectEvents обрабатывает запрос на получение списка всех событий.
func (h *handler) selectEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug("Selecting events")
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectEvents(ctx)
	})
	if err != nil {
		h.logger.Error("selectEvents: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error("selectEvents: " + err.Error())
	}
	h.logger.Info("Events selected")
}

// updateEvent обрабатывает запрос на обновление существующего события.
func (h *handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error("updateEvent: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("Attempting to update event: " + event.ID)
	if err := h.app.UpdateEvent(ctx, event); err != nil {
		h.logger.Error("updateEvent: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Event updated: " + event.ID)
	w.WriteHeader(http.StatusOK)
}

// deleteEvent обрабатывает запрос на удаление события по его ID.
func (h *handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventID := getIDFromPath(r.URL.Path)
	if eventID == "" {
		h.logger.Error("deleteEvent: missing event ID in path")
		http.Error(w, "missing event ID", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Attempting to delete event: " + eventID)
	if err := h.app.DeleteEvent(ctx, eventID); err != nil {
		h.logger.Error("deleteEvent: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Event deleted: " + eventID)
	w.WriteHeader(http.StatusOK)
}

// selectEventsForDay обрабатывает запрос на получение списка событий на указанный день.
func (h *handler) selectEventsForDay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	date, err := parseDateFromQuery(r, "date")
	if err != nil {
		h.logger.Error(selectEventsForDayMsg + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("Selecting events for day: " + date.String())
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectEventsForDay(ctx, date)
	})
	if err != nil {
		h.logger.Error(selectEventsForDayMsg + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error(selectEventsForDayMsg + err.Error())
	}
	h.logger.Info("Events for day selected")
}

// selectEventsForWeek обрабатывает запрос на получение списка событий на указанную неделю.
func (h *handler) selectEventsForWeek(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDate, err := parseDateFromQuery(r, "startDate")
	if err != nil {
		h.logger.Error(selectEventsForWeekMsg + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("Selecting events for week starting: " + startDate.String())
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectEventsForWeek(ctx, startDate)
	})
	if err != nil {
		h.logger.Error(selectEventsForWeekMsg + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error(selectEventsForWeekMsg + err.Error())
	}
	h.logger.Info("Events for week selected")
}

// handleRoute обрабатывает запросы к /route.
func (h *handler) handleRoute(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{"message": "Route handler reached"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleHealth обрабатывает запросы к /health для проверки состояния сервиса.
func (h *handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

// selectEventsForMonth обрабатывает запрос на получение списка событий на указанный месяц.
func (h *handler) selectEventsForMonth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDate, err := parseDateFromQuery(r, "startDate")
	if err != nil {
		h.logger.Error(selectEventsForMonthMsg + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("Selecting events for month starting: " + startDate.String())
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectEventsForMonth(ctx, startDate)
	})
	if err != nil {
		h.logger.Error(selectEventsForMonthMsg + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error(selectEventsForMonthMsg + err.Error())
	}
	h.logger.Info("Events for month selected")
}

// readUserFromBody читает и разбирает тело запроса в структуру User.
func readUserFromBody(r *http.Request) (*model.User, error) {
	defer r.Body.Close()
	user := new(model.User)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, user); err != nil {
		return nil, err
	}
	return user, nil
}

// readEventFromBody читает и разбирает тело запроса в структуру Event.
func readEventFromBody(r *http.Request) (*model.Event, error) {
	defer r.Body.Close()
	event := new(model.Event)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, event); err != nil {
		return nil, err
	}
	return event, nil
}

// selectAsJSON выполняет функцию селектора и сериализует результат в JSON.
func selectAsJSON(ctx context.Context, sel func(context.Context) (interface{}, error)) ([]byte, error) {
	data, err := sel(ctx)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

// sendData отправляет данные в формате JSON в HTTP-ответ.
func sendData(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

// getIDFromPath извлекает ID из пути запроса.
func getIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[len(parts)-1]
}

// parseDateFromQuery извлекает и парсит дату из параметров запроса.
func parseDateFromQuery(r *http.Request, key string) (time.Time, error) {
	dateStr := r.URL.Query().Get(key)
	if dateStr == "" {
		return time.Time{}, errors.New("missing " + key + " in query")
	}
	return time.Parse("2006-01-02", dateStr)
}
