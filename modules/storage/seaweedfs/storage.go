package seaweedfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/czhj/ahfs/modules/storage"
)

const SeaweedfsStorageType storage.Type = "seaweedfs"

type SeaweedfsStorageConfig struct {
	Host string `json:"host"`
}

func (c *SeaweedfsStorageConfig) DirAssignUrl() string {
	url := &url.URL{
		Scheme: "http",
		Host:   c.Host,
		Path:   "/dir/assign",
	}
	return url.String()
}

func (c *SeaweedfsStorageConfig) DirLookupUrl(volumnID string) string {
	url := &url.URL{
		Scheme:   "http",
		Host:     c.Host,
		Path:     "/dir/lookup",
		RawQuery: fmt.Sprintf("volumeId=%s", volumnID),
	}
	return url.String()
}

type Storage struct {
	config SeaweedfsStorageConfig
}

func NewSeaweedfsStorage(ctx context.Context, cfg interface{}) (storage.Storage, error) {
	configInterface, err := storage.ToConfig(SeaweedfsStorageConfig{}, cfg)
	if err != nil {
		return nil, err
	}

	config := configInterface.(SeaweedfsStorageConfig)

	return &Storage{
		config: config,
	}, nil
}

func (s *Storage) Write(f *storage.Object, opts ...storage.WriteOption) (storage.ID, error) {
	dirAssign, err := s.requestDirAssign()
	if err != nil {
		return "", err
	}

	_, err = s.sendDataByDirAssign(dirAssign, f)
	if err != nil {
		return "", err
	}

	return storage.ID(dirAssign.FID), nil
}

func (s *Storage) Read(id storage.ID, opts ...storage.ReadOption) (*storage.Object, error) {
	url, err := s.makeRequestURLByFID(string(id))
	if err != nil {
		return nil, err
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("Failed to read object from seaweedfs, id: %s, status_code: %d", id, response.StatusCode)
	}

	return &storage.Object{
		Reader: response.Body,
		Size:   response.ContentLength,
	}, nil
}

func (s *Storage) Delete(id storage.ID) error {
	url, err := s.makeRequestURLByFID(string(id))
	if err != nil {
		return err
	}

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	defer func() {
		if response.Body != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return err
	}

	if http.StatusOK != response.StatusCode {
		if http.StatusNotFound == response.StatusCode {
			return storage.ErrNotFound
		}
		return fmt.Errorf("Failed to delete object from seaweedfs, id: %s, status_code: %d", id, response.StatusCode)
	}

	return nil
}

func (s *Storage) requestDirAssign() (*DirAssign, error) {
	var dirAssign DirAssign
	if err := s.GetJsonByURL("GET", s.config.DirAssignUrl(),
		nil, &dirAssign); err != nil {
		return nil, err
	}

	return &dirAssign, nil
}

func (s *Storage) sendDataByDirAssign(d *DirAssign, f *storage.Object) (*UploadResponse, error) {
	body, err := s.createBodyFromObject(f)
	if err != nil {
		return nil, err
	}

	response := &UploadResponse{}
	if err := s.GetJsonByURL("POST", d.ToHttpURL(), body, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Storage) createBodyFromObject(f *storage.Object) (io.Reader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", f.Name)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, f.Reader)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *Storage) GetJsonByURL(method, url string, body io.Reader, v interface{}) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	defer func() {
		if response.Body != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send request to seaweedfs,url: %s, status_code: %d", url, response.StatusCode)
	}

	body = response.Body
	if body == nil {
		return fmt.Errorf("Failed to get response body")
	}

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

func (s *Storage) makeRequestURLByFID(fid string) (string, error) {
	volumeID := s.fidToVolumeID(fid)
	dirLookupUrl := s.config.DirLookupUrl(volumeID)

	dirLookup := &DirLookup{}
	if err := s.GetJsonByURL("GET", dirLookupUrl, nil, dirLookup); err != nil {
		return "", err
	}

	if len(dirLookup.Locations) == 0 {
		return "", fmt.Errorf("VolumeID not found in any url")
	}

	dirAssign := &DirAssign{
		FID:       fid,
		PublicURL: dirLookup.Locations[0].PublicURL,
	}

	return dirAssign.ToHttpURL(), nil
}

func (s *Storage) fidToVolumeID(fid string) string {
	strArray := strings.Split(fid, ",")
	if len(strArray) != 2 {
		return ""
	}
	return strArray[0]
}

func init() {
	storage.RegisterStorageGenerator(SeaweedfsStorageType, NewSeaweedfsStorage)
}
