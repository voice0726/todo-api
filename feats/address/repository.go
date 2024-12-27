package address

import (
	"context"

	"github.com/google/uuid"

	"github.com/voice0726/todo-app-api/infra"
	"github.com/voice0726/todo-app-api/models"
)

type Repository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Address, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Address, error)
	Update(ctx context.Context, address *models.Address) (*models.Address, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RepositoryImpl struct {
	db *infra.DataBase
}

func NewRepositoryImpl(db *infra.DataBase) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (r *RepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	var address models.Address
	if err := r.db.WithContext(ctx).First(&address, id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *RepositoryImpl) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Address, error) {
	var addresses []*models.Address
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *RepositoryImpl) Update(ctx context.Context, address *models.Address) (*models.Address, error) {
	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (r *RepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Address{}, id).Error; err != nil {
		return err
	}
	return nil
}
