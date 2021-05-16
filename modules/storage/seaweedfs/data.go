package seaweedfs

import "fmt"

type DirAssign struct {
	Count     int    `json:"count"`
	FID       string `json:"fid"`
	URL       string `json:"url"`
	PublicURL string `json:"publicUrl"`
}

func (d *DirAssign) ToHttpURL() string {
	return fmt.Sprintf("http://%s/%s", d.PublicURL, d.FID)
}

type UploadResponse struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	ETag string `json:"eTag"`
}

type Location struct {
	PublicURL string `json:"publicUrl"`
	URL       string `json:"url"`
}

type DirLookup struct {
	VolumnID  string     `json:"volumeId"`
	Locations []Location `json:"locations"`
}
