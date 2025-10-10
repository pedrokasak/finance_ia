package database

import (
	"finance-ia/internal/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindAll() ([]*user.User, error) {
	var users []*user.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return users, nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err == nil {
        return &user.User{
            ID:        u.ID,
            Email:     u.Email,
            Password:  u.Password,
            CreatedAt: u.CreatedAt,
            UpdatedAt: u.UpdatedAt,
        }, nil
    }
	return &u, nil
}

func (r *UserRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) Delete(u *user.User) error {
	return r.db.Delete(u).Error
}

var _ user.Repository = (*UserRepository)(nil)
