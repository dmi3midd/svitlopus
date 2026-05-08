package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"svitlopus/internal/models"

	"github.com/jmoiron/sqlx"
)

var (
	ErrWorkerNotFound = errors.New("worker not found")
)

type WorkerRepository interface {
	// GetWorkerById returns a worker by its ID.
	// It returns the ErrWorkerNotFound error if the worker does not exist.
	GetWorkerById(context context.Context, id string) (*models.Worker, error)
	// GetAllWorkers returns a list of all workers.
	// It returns an empty slice if no workers are found.
	GetAllWorkers(context context.Context) ([]models.Worker, error)
	// CreateWorker adds a new worker to the database.
	CreateWorker(context context.Context, worker *models.Worker) (*models.Worker, error)
	// DeleteWorker removes a worker from the database by its ID.
	DeleteWorker(context context.Context, id string) error
}

type workerRepository struct {
	db *sqlx.DB
}

func NewWorkerRepository(db *sqlx.DB) WorkerRepository {
	return &workerRepository{db: db}
}

func (w *workerRepository) GetWorkerById(ctx context.Context, id string) (*models.Worker, error) {
	op := "WorkerRepository.GetById"

	var worker models.Worker
	query := `
	SELECT id, username, bot_token, is_active, chat_id, created_at
	FROM workers
	WHERE id = $1
	`
	err := w.db.GetContext(ctx, &worker, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrWorkerNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &worker, nil
}

func (w *workerRepository) GetAllWorkers(ctx context.Context) ([]models.Worker, error) {
	op := "WorkerRepository.GetAllWorkers"

	var workers []models.Worker
	query := `
	SELECT id, username, bot_token, is_active, chat_id, created_at
	FROM workers
	`
	err := w.db.SelectContext(ctx, &workers, query)
	if err != nil {
		return []models.Worker{}, fmt.Errorf("%s: %w", op, err)
	}
	return workers, nil
}

func (w *workerRepository) CreateWorker(ctx context.Context, worker *models.Worker) (*models.Worker, error) {
	op := "WorkerRepository.CreateWorker"
	query := `
	INSERT INTO workers (id, username, bot_token, is_active, chat_id, created_at) 
	VALUES (:id, :username, :bot_token, :is_active, :chat_id, :created_at) 
	`
	if _, err := w.db.NamedExecContext(ctx, query, worker); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return worker, nil
}

func (w *workerRepository) DeleteWorker(ctx context.Context, id string) error {
	op := "WorkerRepository.DeleteWorker"
	query := `
	DELETE FROM workers
	WHERE id = $1
	`
	if _, err := w.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
