package usecase

import (
	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/candbright/mc-dashboard/internal/pkg/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository
	log  *logrus.Logger
}

// NewUserService 创建用户服务实例
func NewUserService(repo domain.UserRepository, log *logrus.Logger) domain.UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (s *UserService) Register(phone, email, password string) (*domain.User, error) {
	// 检查手机号是否已存在
	if _, err := s.repo.FindByPhone(phone); err == nil {
		return nil, errors.ErrPhoneRegistered
	}

	// 如果提供了邮箱，检查邮箱是否已存在
	if email != "" {
		if _, err := s.repo.FindByEmail(email); err == nil {
			return nil, errors.ErrEmailRegistered
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrGeneratePassword
	}

	user := &domain.User{
		Phone:    phone,
		Email:    email,
		Password: string(hashedPassword),
		Status:   0,
	}

	// 13800138000暂时作为admin号码
	if phone == "13800138000" {
		user.Status = 1
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.ErrCreateUser
	}

	return user, nil
}

func (s *UserService) Login(loginType string, account, credential string) (string, error) {
	var user *domain.User
	var err error

	switch loginType {
	case "phone_password":
		user, err = s.repo.FindByPhone(account)
		if err != nil {
			return "", errors.ErrUserNotFound
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential)); err != nil {
			return "", errors.ErrWrongPassword
		}
	case "email_password":
		user, err = s.repo.FindByEmail(account)
		if err != nil {
			return "", errors.ErrUserNotFound
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential)); err != nil {
			return "", errors.ErrWrongPassword
		}
	case "phone_code":
		user, err = s.repo.FindByPhone(account)
		if err != nil {
			return "", errors.ErrUserNotFound
		}
		// TODO: 实现短信验证码验证逻辑
	default:
		return "", errors.ErrUnsupportedLoginType
	}

	return s.generateToken(user)
}

func (s *UserService) GetUserInfo(id uint) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateUserInfo(user *domain.User) error {
	if err := s.repo.Update(user); err != nil {
		return errors.ErrUpdateUser
	}
	return nil
}

// Logout 用户登出
func (s *UserService) Logout(id uint) error {
	// 由于使用 JWT，服务端不需要维护会话状态
	// 客户端只需要清除本地存储的 token 即可
	// 这里可以添加一些额外的清理逻辑，比如清除用户的在线状态等
	s.log.WithField("user_id", id).Info("用户登出")
	return nil
}

// generateToken 生成JWT令牌
func (s *UserService) generateToken(user *domain.User) (string, error) {
	return middleware.GenerateToken(user.ID, user.Phone, user.Status)
}
