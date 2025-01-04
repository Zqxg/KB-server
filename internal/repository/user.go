package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "projectName/api/v1"
	"projectName/internal/model"
	"time"
)

type UserRepository interface {
	// 表：sys_users
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByUserId(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	DeleteByUserId(ctx context.Context, userId string) error
	// 表：sys_user_auths
	CreateUserAuth(ctx context.Context, userAuth *model.UserAuth) error
	GetUserAuthByUserId(ctx context.Context, userId string) (*model.UserAuth, error)
	UpdateUserAuth(ctx context.Context, userAuth *model.UserAuth) error
	// redis
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

func NewUserRepository(
	r *Repository,
) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Table("sys_users").Create(user).Error; err != nil {
		r.logger.WithContext(ctx).Error("userRepository.Create error", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Table("sys_users").Save(user).Error; err != nil {
		r.logger.WithContext(ctx).Error("userRepository.Update error", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) GetByUserId(ctx context.Context, userId string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Table("sys_users").Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		r.logger.WithContext(ctx).Error("userRepository.GetByUserId error", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Table("sys_users").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.WithContext(ctx).Error("userRepository.GetByEmail error", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Table("sys_users").Where("is_deleted = 0 and phone = ? ", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.WithContext(ctx).Error("userRepository.GetByPhone error", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) DeleteByUserId(ctx context.Context, userId string) error {
	// 获取当前时间
	now := time.Now()
	// 更新 is_deleted 字段为 1，并设置 deleted_at 为当前时间
	if err := r.DB(ctx).Table("sys_users").
		Where("user_id = ?", userId).
		Update("is_deleted", 1).
		Update("deleted_at", now).Error; err != nil {
		r.logger.WithContext(ctx).Error("userRepository.DeleteByUserId error", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.rdb.Set(ctx, key, value, expiration).Err(); err != nil {
		r.logger.WithContext(ctx).Error("userRepository.Set error", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) Get(ctx context.Context, key string) (string, error) {
	if err := r.rdb.Get(ctx, key).Err(); err != nil {
		r.logger.WithContext(ctx).Error("userRepository.Get error", zap.Error(err))
		return "", err
	}
	return r.rdb.Get(ctx, key).Result()
}

func (r *userRepository) Delete(ctx context.Context, key string) error {
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		r.logger.WithContext(ctx).Error("userRepository.Delete error", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) CreateUserAuth(ctx context.Context, userAuth *model.UserAuth) error {
	if err := r.DB(ctx).Table("sys_user_auths").Create(userAuth).Error; err != nil {
		r.logger.WithContext(ctx).Error("userRepository.CreateUserAuth error", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) GetUserAuthByUserId(ctx context.Context, userId string) (*model.UserAuth, error) {
	var userAuth model.UserAuth
	if err := r.DB(ctx).Table("sys_user_auths").Where("user_id =?", userId).First(&userAuth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.WithContext(ctx).Error("userRepository.GetUserAuthByUserId error", zap.Error(err))
			return nil, nil
		}
	}
	return &userAuth, nil
}

func (r *userRepository) UpdateUserAuth(ctx context.Context, userAuth *model.UserAuth) error {
	if err := r.DB(ctx).Table("sys_user_auths").Save(userAuth).Error; err != nil {
		r.logger.WithContext(ctx).Error("userRepository.UpdateUserAuth error", zap.Error(err))
		return err
	}
	return nil
}
