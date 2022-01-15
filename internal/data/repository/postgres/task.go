package repository

import (
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresTaskRepository struct {
	db *sqlx.DB
}

var _ interfaces.TaskRepository = (*PostgresTaskRepository)(nil)

func newPostgresTaskRepository(db *sqlx.DB) PostgresTaskRepository {
	return PostgresTaskRepository{db: db}
}

func (r PostgresTaskRepository) Get(taskId int) (entity.Task, error)
func (r PostgresTaskRepository) GetMany(taskParams interface{}) ([]entity.Task, error)
func (r PostgresTaskRepository) Solve(taskId int) error
func (r PostgresTaskRepository) Save(task entity.Task) (taskId int, err error)
func (r PostgresTaskRepository) Update(oldTask, newTask entity.Task) error
