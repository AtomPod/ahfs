package setting

var (
	LFS = struct {
		Storage
	}{}
)

func newLFSService() {
	LFS.Storage = getStorage("lfs")
}
