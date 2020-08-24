package models

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/czhj/ahfs/modules/locker"
	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type FileType int

const (
	FileTypeFile FileType = iota
	FileTypeDir
)

type File struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	FileID   string `gorm:"unique_index"`
	FileDir  string
	FileName string

	FileType FileType
	FileSize uint64

	Owner    uint
	ParentID uint
}

func (f *File) IsDir() bool {
	return f.FileType == FileTypeDir
}

func (f *File) LocalPath() string {
	return filepath.Join(setting.FileUploadPath, f.FileID)
}

func (f *File) FilePath() string {
	return path.Join(f.FileDir, f.FileName)
}

func (f *File) ReadDir() ([]*File, error) {
	if !f.IsDir() {
		return nil, ErrFileNotDirectory{ID: f.ID, Path: f.FilePath()}
	}

	files := make([]*File, 0)
	err := engine.Where("parent_id=?", files).Find(&files).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return files, nil
}

func CreateUserRootFile(u *User) *File {
	return &File{
		FileID:   fmt.Sprintf("%d-root", u.ID),
		FileDir:  "/",
		FileType: FileTypeDir,
		ParentID: 0,
		Owner:    u.ID,
	}
}

func GetUserRootFile(uid uint) (*File, error) {
	file := new(File)
	result := engine.Where("id=? AND file_id=?", uid, fmt.Sprintf("%d-root", uid)).First(file)
	if err := result.Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrFileNotExist{Owner: uid}
		}
		return nil, err
	}
	return file, nil
}

func GetFileByID(id uint, uid uint) (*File, error) {
	return getFileByID(engine, id, uid)
}

func getFileByID(e *gorm.DB, id uint, uid uint) (*File, error) {
	if id == 0 && uid == 0 {
		return nil, ErrFileNotExist{}
	}

	var query *gorm.DB
	if id != 0 {
		query = e.Where("id=?", id)
	} else {
		root := CreateUserRootFile(&User{ID: uid})
		query = e.Where(root)
	}

	if uid != 0 {
		query = query.Where("owner=?", uid)
	}

	file := new(File)
	err := query.First(file).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrFileNotExist{ID: id, Owner: uid}
		}
		return nil, err
	}

	return file, nil
}

func CreateFile(f *File) error {
	return createFile(engine, f)
}

func createFile(e *gorm.DB, f *File) error {
	return e.Create(f).Error
}

func DeleteFile(f *File) error {
	uid := f.Owner
	id, err := LockUserFile(context.Background(), uid)
	if err != nil {
		return err
	}
	defer func() {
		if err := UnlockUserFile(context.Background(), uid, id); err != nil {
			log.Error("Failed to unlock user file", zap.Uint("id", f.ID), zap.Uint("uid", uid), zap.Error(err))
		}
	}()

	tx := engine.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	if err := deleteFile(tx, f); err != nil {
		return err
	}

	return tx.Commit().Error
}

func deleteFile(e *gorm.DB, f *File) error {

	if f.IsDir() {
		files, err := f.ReadDir()
		if err != nil {
			return err
		}

		for _, file := range files {
			if err := deleteFile(e, file); err != nil {
				return err
			}
		}
	}

	if err := e.Delete(f).Error; err != nil {
		return err
	}

	localPath := f.LocalPath()
	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("Failed to remove file [%s]: %v", localPath, err)
	}

	return nil
}

func LockUserFile(ctx context.Context, uid uint) (id string, err error) {
	key := fmt.Sprintf("user-%d-file", uid)
	return locker.Lock(ctx, key)
}

func UnlockUserFile(ctx context.Context, uid uint, id string) error {
	key := fmt.Sprintf("user-%d-file", uid)
	return locker.Unlock(ctx, key, id)
}
