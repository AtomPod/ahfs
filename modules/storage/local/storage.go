package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/czhj/ahfs/modules/setting"
	"github.com/czhj/ahfs/modules/storage"
	"github.com/czhj/ahfs/modules/utils"
)

const LocalStorageType storage.Type = "local"

type LocalStorageConfig struct {
	Directory string `json:"directory"`
}

func (c LocalStorageConfig) GetDirectory() string {
	if len(c.Directory) == 0 {
		return setting.FileUploadPath
	}
	return c.Directory
}

type Storage struct {
	config LocalStorageConfig
}

func NewStorage(cfg LocalStorageConfig) *Storage {
	return &Storage{
		config: cfg,
	}
}

func NewLocalStorage(ctx context.Context, cfg interface{}) (storage.Storage, error) {
	lsc, err := storage.ToConfig(LocalStorageConfig{}, cfg)

	if err != nil {
		return nil, err
	}

	config := lsc.(LocalStorageConfig)

	return NewStorage(config), nil
}

func (s *Storage) Write(f *storage.Object, opts ...storage.WriteOption) (storage.ID, error) {
	wo := s.makeWriteOptions(opts...)

	if err := s.checkDir(); err != nil {
		return "", fmt.Errorf("LocalStorage: %v", err)
	}

	id := utils.GenerateFileID(wo.ID)
	localPath := filepath.Join(s.config.GetDirectory(), id)

	file, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("Failed to create local file [%s]: %v", localPath, err)
	}

	if _, err := io.Copy(file, f.Reader); err != nil {
		return "", fmt.Errorf("Failed to copy file to local file [%s]: %v", localPath, err)
	}

	return storage.ID(id), nil
}

func (s *Storage) Read(id storage.ID, opts ...storage.ReadOption) (*storage.Object, error) {

	localPath := filepath.Join(s.config.GetDirectory(), string(id))

	file, err := os.Open(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &storage.Object{
		Reader: file,
	}, nil

}

func (s *Storage) Delete(id storage.ID) error {
	localPath := filepath.Join(s.config.GetDirectory(), string(id))

	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("Failed to remove file [%s]: %v", localPath, err)
	}

	return nil
}

func (s *Storage) checkDir() error {
	if err := os.MkdirAll(s.config.GetDirectory(), os.ModePerm); err != nil {
		return fmt.Errorf("Failed to run MkdirAll [%s]: %v", os.ModePerm, err)
	}
	return nil
}

func (s *Storage) makeWriteOptions(opts ...storage.WriteOption) *storage.WriteOptions {
	wo := &storage.WriteOptions{}

	for _, o := range opts {
		o(wo)
	}
	return wo
}

func init() {
	storage.RegisterStorageGenerator(LocalStorageType, NewLocalStorage)
}
