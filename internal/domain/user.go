package domain

import "time"

// User 用户实体
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Phone     string    `gorm:"size:11;uniqueIndex" json:"phone"` // 手机号
	Email     string    `gorm:"size:50;uniqueIndex" json:"email"` // 邮箱
	Password  string    `gorm:"size:100" json:"-"`                // 密码
	Nickname  string    `gorm:"size:50" json:"nickname"`          // 昵称
	Avatar    string    `gorm:"size:255" json:"avatar"`           // 头像
	Status    int       `gorm:"default:0" json:"status"`          // 1: 正常, 2: 禁用
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "mc_dashboard_users"
}

// UserService 用户服务接口
type UserService interface {
	Register(phone, email, password string) (*User, error)
	Login(loginType string, account, credential string) (string, error)
	GetUserInfo(id uint) (*User, error)
	UpdateUserInfo(user *User) error
	Logout(id uint) error
}

// UserRepository 用户仓储接口
type UserRepository interface {
	BaseRepository[User]
	FindByPhone(phone string) (*User, error)
	FindByEmail(email string) (*User, error)
}
