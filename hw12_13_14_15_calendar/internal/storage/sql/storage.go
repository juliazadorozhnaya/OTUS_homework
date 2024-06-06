package sqlstorage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Storage{
		pool: pool,
	}, nil
}

// SelectUsers - возвращает всех пользователей из базы данных.
func (s *Storage) SelectUsers(ctx context.Context) (users []model.User, err error) {
	users = make([]model.User, 0)
	sql := `SELECT id, firstname, lastname, email, age FROM calendar.users;`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return users, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows, err := tx.Query(ctx, sql)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// CreateUser - вставляет нового пользователя в базу данных.
func (s *Storage) CreateUser(ctx context.Context, user model.User) error {
	sql := `INSERT INTO calendar.users (firstname, lastname, email, age) VALUES ($1, $2, $3, $4);`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, sql, user.FirstName, user.LastName, user.Email, user.Age)
	return err
}

// DeleteUser - удаляет пользователя по его идентификатору.
func (s *Storage) DeleteUser(ctx context.Context, userID string) error {
	sql := `DELETE FROM calendar.users WHERE id = $1;`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, sql, userID)
	return err
}

// CreateEvent - вставляет новое событие в базу данных.
func (s *Storage) CreateEvent(ctx context.Context, event model.Event) error {
	sql := `INSERT INTO calendar.events (title, description, beginning, finish, notification, userid) 
			VALUES ($1, $2, $3, $4, $5, $6);`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, sql, event.Title, event.Description, event.Beginning, event.Finish,
		event.Notification, event.UserID)
	return err
}

// DeleteEvent - удаляет событие по его идентификатору.
func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	sql := `DELETE FROM calendar.events WHERE id = $1;`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, sql, eventID)
	return err
}

// UpdateEvent - обновляет существующее событие в базе данных.
func (s *Storage) UpdateEvent(ctx context.Context, event model.Event) error {
	sql := `UPDATE calendar.events
			SET title = $2, description = $3, beginning = $4, finish = $5, notification = $6, userid = $7
			WHERE id = $1;`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, sql, event.ID, event.Title, event.Description, event.Beginning, event.Finish,
		event.Notification, event.UserID)
	return err
}

// SelectEvents - возвращает все события из базы данных.
func (s *Storage) SelectEvents(ctx context.Context) (events []model.Event, err error) {
	events = make([]model.Event, 0)
	sql := `SELECT id, title, description, beginning, finish, notification, userid FROM calendar.events;`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return events, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows, err := tx.Query(ctx, sql)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var event model.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.Beginning, &event.Finish,
			&event.Notification, &event.UserID)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}
