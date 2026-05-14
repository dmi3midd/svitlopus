package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRepoWorkerNotFound = errors.New("repository: worker not found")
)

type WorkerRepository interface {
	// GetById returns a worker by its ID.
	// It returns the ErrRepoWorkerNotFound error if the worker does not exist.
	GetById(ctx context.Context, id string) (*Worker, error)
	// GetAll returns a list of all workers.
	// It returns an empty slice if no workers are found.
	GetAll(ctx context.Context) ([]Worker, error)
	// Create adds a new worker to the database.
	Create(ctx context.Context, worker *Worker) (*Worker, error)
	// Update updates is_active field.
	// It returns the ErrRepoWorkerNotFound error if the worker does not exist.
	Update(ctx context.Context, id string, isActive bool) error
	// Delete removes a worker from the database by its ID.
	Delete(ctx context.Context, id string) error
}

type workerRepository struct {
	db *sqlx.DB
}

func NewWorkerRepository(db *sqlx.DB) WorkerRepository {
	return &workerRepository{
		db: db,
	}
}

func (r *workerRepository) GetById(ctx context.Context, id string) (*Worker, error) {
	op := "WorkerRepository.GetById"

	var worker Worker
	query := `
	SELECT id, username, bot_token, is_active, chat_id, created_at
	FROM workers
	WHERE id = $1
	`
	err := r.db.GetContext(ctx, &worker, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrRepoWorkerNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &worker, nil
}

func (r *workerRepository) GetAll(ctx context.Context) ([]Worker, error) {
	op := "WorkerRepository.GetAll"

	var workers []Worker
	query := `
	SELECT id, username, bot_token, is_active, chat_id, created_at
	FROM workers
	`
	err := r.db.SelectContext(ctx, &workers, query)
	if err != nil {
		return []Worker{}, fmt.Errorf("%s: %w", op, err)
	}
	return workers, nil
}

func (r *workerRepository) Create(ctx context.Context, worker *Worker) (*Worker, error) {
	op := "WorkerRepository.Create"
	query := `
	INSERT INTO workers (id, username, bot_token, is_active, chat_id, created_at) 
	VALUES (:id, :username, :bot_token, :is_active, :chat_id, :created_at) 
	`
	if _, err := r.db.NamedExecContext(ctx, query, worker); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return worker, nil
}

func (r *workerRepository) Update(ctx context.Context, id string, isActive bool) error {
	op := "WorkerRepository.Update"
	query := `
	UPDATE workers
	SET is_active = $1
	WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, isActive, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrRepoWorkerNotFound)
	}
	return nil
}

func (r *workerRepository) Delete(ctx context.Context, id string) error {
	op := "WorkerRepository.Delete"
	query := `
	DELETE FROM workers
	WHERE id = $1
	`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
