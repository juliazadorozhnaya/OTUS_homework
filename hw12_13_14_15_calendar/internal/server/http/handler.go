package serverhttp

import (
	"context"
	"encoding/json"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"io"
	"net/http"
)

type handler struct {
	logger server.Logger
	app    server.Application
}

func newHandler(logger server.Logger, app server.Application) *handler {
	return &handler{
		logger: logger,
		app:    app,
	}
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := readUserFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.CreateUser(ctx, user); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) selectUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectUsers(ctx)
	})
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := readUserFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteUser(ctx, user.GetID()); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) createEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.CreateEvent(ctx, event); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) selectEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	marshal, err := selectAsJSON(ctx, func(ctx context.Context) (interface{}, error) {
		return h.app.SelectEvents(ctx)
	})
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendData(w, marshal); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.UpdateEvent(ctx, event); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteEvent(ctx, event.GetID()); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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

func selectAsJSON(ctx context.Context, sel func(context.Context) (interface{}, error)) ([]byte, error) {
	events, err := sel(ctx)
	if err != nil {
		return nil, err
	}
	return json.Marshal(events)
}

func sendData(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}
