package models

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/czhj/ahfs/modules/avatar"
	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserType int

const (
	UserTypeUser UserType = iota
)

type User struct {
	ID          uint       `gorm:"primary_key"`
	CreatedAt   time.Time  `sql:"index"`
	UpdatedAt   time.Time  `sql:"index"`
	DeletedAt   *time.Time `sql:"index"`
	LastLoginAt time.Time  `sql:"index"`

	Nickname string `gorm:"not null"`
	Username string `gorm:"unique_index,not null"`
	Email    string `gorm:"unique_index,not null"`
	Password string `gorm:"not null" json:"-"`

	MustChangePassword bool `gorm:"not null,default:false"`

	IsActive bool `sql:"index"`
	IsAdmin  bool

	LoginType LoginType
	Type      UserType

	Avatar string `gorm:"varchar(2048)"`

	UsedFileCapacity int64
	MaxFileCapacity  int64
}

func (u *User) AvatarLink() string {
	return strings.TrimRight(setting.AppSubURL, "/") + "/avatars/" + u.Avatar
}

func (u *User) AvatarPath() string {
	return filepath.Join(setting.AvatarUploadPath, u.Avatar)
}

func (u *User) IsOAuth2() bool {
	return u.LoginType == LoginOAuth2
}

func (u *User) IsPlain() bool {
	return u.LoginType == LoginPlain
}

func (u *User) RemainingFileSize() int64 {
	return u.MaxFileCapacity - u.UsedFileCapacity
}

func (u *User) CanIUploadFile(size int64) bool {
	return u.RemainingFileSize() >= size
}

func (u *User) HashPassword(pwd string) {
	hashedpwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	u.Password = string(hashedpwd)
}

func (u *User) ValidatePassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd)) == nil
}

func (u *User) IsPasswordSet() bool {
	return !u.ValidatePassword("")
}

func (u *User) UploadAvatar(data []byte) error {
	m, err := avatar.Prepare(data)
	if err != nil {
		return err
	}

	tx := engine.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	u.Avatar = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d-%x", u.ID, md5.Sum(data)))))
	if err := saveUser(tx, u); err != nil {
		return fmt.Errorf("saveUser: %v", err)
	}

	if err := os.MkdirAll(setting.AvatarUploadPath, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create dir %s: %v", setting.AvatarUploadPath, err)
	}

	fw, err := os.Create(u.AvatarPath())
	if err != nil {
		return fmt.Errorf("Cannot create file %s: %v", u.AvatarPath(), err)
	}
	defer fw.Close()

	if err := png.Encode(fw, *m); err != nil {
		return fmt.Errorf("Cannot encode png: %v", err)
	}

	return tx.Commit().Error
}

func (u *User) DeleteAvatar() error {
	if len(u.Avatar) > 0 {
		if err := os.Remove(u.AvatarPath()); err != nil {
			return fmt.Errorf("Failed to remove file %s: %v", u.AvatarPath(), err)
		}
	}

	u.Avatar = ""
	if err := SaveUser(u); err != nil {
		return fmt.Errorf("SaveUser: %v", err)
	}

	return nil
}

func (u *User) GetRootDir() ([]*File, error) {
	root, err := GetUserRootFile(u.ID)
	if err != nil {
		return nil, err
	}
	return root.ReadDir(ReadDirOption{})
}

func IsUserExists(uid int64) (bool, error) {
	return isUserExists(engine, uid)
}

func isUserExists(e *gorm.DB, uid int64) (bool, error) {
	var count int
	err := e.Model(&User{}).Where("id=?", uid).Limit(1).Count(&count).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return false, err
		}
	}
	return count != 0, nil
}

func CreateUser(u *User) error {
	tx := engine.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	used, err := isEmailUsed(tx, u.Email)
	if err != nil {
		return err
	} else if used {
		return ErrEmailAlreadyUsed{Email: u.Email}
	}

	used, err = isUsernameUsed(tx, u.Username)
	if err != nil {
		return err
	} else if used {
		return ErrUsernameAlreadyUsed{Username: u.Username}
	}

	u.HashPassword(u.Password)
	if u.MaxFileCapacity == 0 {
		u.MaxFileCapacity = setting.Service.MaxFileCapacitySize
	}

	count, err := countUser(tx)
	if err != nil {
		return err
	}

	// first user is admin
	if count == 0 {
		u.IsAdmin = true
	}

	log.Debug("Create user", zap.Any("user", u))
	if err := tx.Create(u).Error; err != nil {
		return err
	}

	root := CreateUserRootFile(u)
	if err := createFile(tx, root); err != nil {
		return err
	}

	return tx.Commit().Error
}

func DeleteUser(u *User) error {
	tx := engine.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	if err := deleteUser(tx, u); err != nil {
		return err
	}

	return tx.Commit().Error
}

func deleteUser(e *gorm.DB, u *User) error {
	db := e.Delete(u)
	if err := db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return ErrUserNotExist{ID: u.ID}
	}

	root, err := getFileByID(e, 0, u.ID)
	if err != nil {
		if !IsErrFileNotExist(err) {
			return err
		}
	}

	id, err := LockUserFile(context.Background(), u.ID)
	if err != nil {
		return err
	}
	defer func() {
		if err := UnlockUserFile(context.Background(), u.ID, id); err != nil {
			log.Error("Failed to unlock user file", zap.Uint("user_id", u.ID), zap.Error(err))
		}
	}()

	if err := deleteFile(e, root); err != nil {
		return err
	}

	return nil
}

func IsEmailUsed(email string) (bool, error) {
	return isEmailUsed(engine, email)
}

func isEmailUsed(e *gorm.DB, email string) (bool, error) {
	var count int64
	err := e.Model(&User{}).Where("email=?", email).Limit(1).Count(&count).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return false, err
	}
	return count != 0, nil
}

func IsUsernameUsed(username string) (bool, error) {
	return isUsernameUsed(engine, username)
}

func isUsernameUsed(e *gorm.DB, username string) (bool, error) {
	var count int64
	err := e.Model(&User{}).Where("username=?", username).Limit(1).Count(&count).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return false, err
	}
	return count != 0, nil
}

func CountUser() (int64, error) {
	return countUser(engine)
}

func countUser(e *gorm.DB) (int64, error) {
	var count int64
	err := e.Model(&User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := new(User)

	err := engine.Where("email=?", email).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrUserNotExist{Email: email}
		}
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := new(User)

	err := engine.Where("username=?", username).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrUserNotExist{Username: username}
		}
		return nil, err
	}
	return user, nil
}

func GetUserByID(uid uint) (*User, error) {
	return getUserByID(engine, uid)
}

func getUserByID(e *gorm.DB, uid uint) (*User, error) {
	user := new(User)

	err := e.Where("id=?", uid).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrUserNotExist{ID: uid}
		}
		return nil, err
	}
	return user, nil
}

func SaveUser(u *User) error {
	return saveUser(engine, u)
}

func saveUser(e *gorm.DB, u *User) error {
	return e.Save(u).Error
}

func UserSignIn(username string, password string) (*User, error) {

	tx := engine.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	user := &User{}

	if strings.Contains(username, "@") {
		user.Email = username
	} else {
		user.Username = username
	}

	result := tx.Where(user).First(user)
	if err := result.Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	switch user.LoginType {
	case LoginNoType, LoginPlain, LoginOAuth2:
		if user.IsPasswordSet() && user.ValidatePassword(password) {

			user.LastLoginAt = time.Now()
			err := tx.Model(user).UpdateColumns(map[string]interface{}{
				"last_login_at": user.LastLoginAt,
			}).Error

			if err != nil {
				return nil, err
			}

			return user, tx.Commit().Error

		}
		return nil, ErrUserNotExist{Email: user.Email}
	}

	return nil, ErrUserNotExist{Email: user.Email}
}

type SearchUserOptions struct {
	ListOptions
	Keyword       string
	Type          UserType
	UID           uint
	IsActive      sql.NullBool
	SearchByEmail bool
	OrderBy       SearchOrderBy
	Actor         *User
}

func (opts *SearchUserOptions) Apply(e *gorm.DB) *gorm.DB {
	db := e.Where("type=?", opts.Type)
	if len(opts.Keyword) > 0 {
		lowerKeyword := strings.ToLower(opts.Keyword)
		if opts.SearchByEmail {
			db = db.Where("(LOWER(nickname) = ?) Or (LOWER(email) = ?) Or (LOWER(username) = ?)", lowerKeyword, lowerKeyword, lowerKeyword)
		} else {
			db = db.Where("(LOWER(nickname) = ?) Or (LOWER(username) = ?)", lowerKeyword, lowerKeyword)
		}
	}

	if opts.IsActive.Valid {
		db = db.Where("is_active=?", opts.IsActive.Valid)
	}

	if opts.UID > 0 {
		db = db.Where("id = ?", opts.UID)
	}
	return db
}

func SearchUser(opts *SearchUserOptions) (users []*User, _ int64, _ error) {
	var count int64

	db := engine.Model(&User{})
	db = opts.Apply(db)
	db = db.Count(&count)
	if err := db.Error; err != nil {
		return nil, 0, fmt.Errorf("Count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByAlphabetically
	}

	db = engine.Model(&User{})
	db = opts.Apply(db)
	db = db.Order(opts.OrderBy.String())
	if opts.Page != 0 {
		db = opts.SetEnginePagination(db)
	}

	users = make([]*User, 0, opts.PageSize)
	err := db.Find(&users).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return users, count, nil
}
