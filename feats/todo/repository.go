package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/voice0726/todo-app-api/infra"
	"gorm.io/gorm"
)

var ErrNotFound = gorm.ErrRecordNotFound

type Repository interface {
	Create(ctx context.Context, todo *Todo) (*Todo, error)
	FindByID(ctx context.Context, id int) (*Todo, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Todo, error)
	Update(ctx context.Context, todo *Todo) (*Todo, error)
	Delete(ctx context.Context, id int) error
}

type RepositoryImpl struct {
	db *infra.DataBase
}

func NewRepositoryImpl(db *infra.DataBase) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	if err := r.db.WithContext(ctx).Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *RepositoryImpl) FindByID(ctx context.Context, id int) (*Todo, error) {
	var todo Todo
	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *RepositoryImpl) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Todo, error) {
	var todos []*Todo
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *RepositoryImpl) Update(ctx context.Context, todo *Todo) (*Todo, error) {
	if err := r.db.WithContext(ctx).Save(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *RepositoryImpl) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&Todo{}, id).Error; err != nil {
		return err
	}
	return nil
}
