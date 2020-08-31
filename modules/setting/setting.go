package setting

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CustomPath        string   = "custom"
	CustomConfigName  string   = "config"
	CustomConfigPaths []string = []string{".", "config"}
	CustomConfigType  string   = "yaml"
	UsedConfigFile    string   = "config.yaml"

	IsWindows bool

	AppName     string
	AppPath     string
	AppWorkPath string
	AppDataPath string
	AppURL      string
	AppSubURL   string

	ServerMode        string
	Protocol          string
	Domain            string
	HTTPAddr          string
	HTTPPort          string
	CertFile          string
	KeyFile           string
	EnableLetsEncrypt bool
	LetsEncryptTOS    bool
	LetsEncryptDir    string
	LetsEncryptHost   []string
	EnableGzip        bool
	StaticRootPath    string

	AvatarMaxWidth   int
	AvatarMaxHeight  int
	AvatarUploadPath string
	FileUploadPath   string

	API struct {
		DefaultPagingSize int
		MaxPagingSize     int
	}
)

func getAppPath() (string, error) {
	var appPath string
	var err error

	if IsWindows && filepath.IsAbs(os.Args[0]) {
		appPath = filepath.Clean(os.Args[0])
	} else {
		appPath, err = exec.LookPath(os.Args[0])
	}

	if err != nil {
		return "", err
	}

	appPath, err = filepath.Abs(appPath)
	if err != nil {
		return "", err
	}

	return strings.Replace(appPath, "\\", "/", -1), nil
}

func getAppWorkPath(appPath string) string {
	workPath := AppWorkPath

	if hWorkPath, ok := os.LookupEnv("AHFS_WORK_DIR"); ok {
		workPath = hWorkPath
	}

	if len(workPath) == 0 {
		s := strings.LastIndex(appPath, "/")
		if s == -1 {
			workPath = appPath
		} else {
			workPath = appPath[:s]
		}
	}

	return strings.Replace(workPath, "\\", "/", -1)

}

func init() {
	IsWindows = runtime.GOOS == "windows"
	log.AddLogger("console", "console", fmt.Sprintf(`{"level": "fatal", "stderr": true}`))
	log.New(fmt.Sprintf(`{"encoding": "console"}`))

	var err error
	if AppPath, err = getAppPath(); err != nil {
		log.Fatal("Failed to get app path", zap.Error(err))
	}
	AppWorkPath = getAppWorkPath(AppPath)
	SetCustomPath()
}

func SetCustomPath() {
	if custom, ok := os.LookupEnv("AHFS_CUSTOM_PATH"); ok {
		CustomPath = custom
	}

	if len(CustomPath) == 0 {
		CustomPath = path.Join(AppWorkPath, "custom")
	} else if !path.IsAbs(CustomPath) {
		CustomPath = path.Join(AppWorkPath, CustomPath)
	}

	UsedConfigFile = path.Join(AppWorkPath, UsedConfigFile)

	for i, p := range CustomConfigPaths {
		CustomConfigPaths[i] = path.Join(CustomPath, p)
	}
}

func NewSetting() {

	for i := range CustomConfigPaths {
		viper.AddConfigPath(CustomConfigPaths[i])
	}

	viper.SetConfigName(CustomConfigName)
	viper.SetConfigType(CustomConfigType)

	viper.SetEnvPrefix("ahfs")
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("Cannot found config, use default")
		} else {
			log.Fatal("Failed to load config", zap.Error(err))
		}
	}

	fileUsed := viper.ConfigFileUsed()
	if len(fileUsed) > 0 {
		UsedConfigFile = fileUsed
		log.Info("Load config file successfully", zap.String("file", fileUsed))
	}

	viper.SetDefault("server", map[string]interface{}{
		"appName":           "AHFS",
		"domain":            "localhost",
		"httpAddr":          "0.0.0.0",
		"httpPort":          "6270",
		"enableLetsEncrypt": false,
		"letsEncryptTOS":    false,
		"letsEncryptDir":    "https",
		"mode":              "debug",
		"avatarMaxWidth":    256,
		"avatarMaxHeight":   256,
	})
	serverCfg := viper.Sub("server")
	AppName = serverCfg.GetString("appName")

	Protocol = "http"
	switch serverCfg.GetString("protocol") {
	case "https":
		Protocol = "https"
		CertFile = serverCfg.GetString("certFile")
		KeyFile = serverCfg.GetString("keyFile")

		if !filepath.IsAbs(CertFile) && len(CertFile) > 0 {
			CertFile = filepath.Join(CustomPath, CertFile)
		}

		if !filepath.IsAbs(KeyFile) && len(KeyFile) > 0 {
			KeyFile = filepath.Join(CustomPath, KeyFile)
		}
	}

	ServerMode = serverCfg.GetString("mode")
	EnableLetsEncrypt = serverCfg.GetBool("enableLetsEncrypt")
	LetsEncryptTOS = serverCfg.GetBool("letsEncryptTOS")
	if !LetsEncryptTOS {
		log.Warn("Failed to enable Let's Encrypt due to Let's Encrypt TOS not beging accpeted")
		EnableLetsEncrypt = false
	}
	LetsEncryptDir = serverCfg.GetString("letsEncryptDir")
	LetsEncryptHost = serverCfg.GetStringSlice("letsEncryptHost")

	Domain = serverCfg.GetString("domain")

	HTTPAddr = serverCfg.GetString("httpAddr")
	HTTPPort = serverCfg.GetString("httpPort")
	EnableGzip = serverCfg.GetBool("enableGZIP")

	serverCfg.SetDefault("staticRootPath", AppWorkPath)
	serverCfg.SetDefault("appDataPath", path.Join(AppWorkPath, "data"))
	StaticRootPath = serverCfg.GetString("staticRootPath")
	AppDataPath = serverCfg.GetString("appDataPath")

	viper.SetDefault("server.avatarUploadPath", filepath.Join(AppDataPath, "avatar"))
	AvatarUploadPath = viper.GetString("server.avatarUploadPath")

	AvatarMaxWidth = serverCfg.GetInt("avatarMaxWidth")
	AvatarMaxHeight = serverCfg.GetInt("avatarMaxHeight")

	viper.SetDefault("server.fileUploadPath", filepath.Join(AppDataPath, "file"))
	FileUploadPath = viper.GetString("server.fileUploadPath")

	defaultAppURL := Protocol + "://" + Domain
	if (Protocol == "http" && HTTPPort != "80") || (Protocol == "https" && HTTPPort != "443") {
		defaultAppURL = defaultAppURL + ":" + HTTPPort
	}
	viper.SetDefault("server.rootURL", defaultAppURL)
	AppURL = viper.GetString("server.rootURL")
	AppURL = strings.TrimSuffix(AppURL, "/") + "/"

	appURL, err := url.Parse(AppURL)
	if err != nil {
		log.Fatal("Invalid RootURL", zap.String("url", AppURL), zap.Error(err))
	}
	AppSubURL = strings.TrimSuffix(appURL.Path, "/")

	newAPIService()
}

func newAPIService() {
	viper.SetDefault("api", map[string]interface{}{
		"defaultPagingSize": 16,
		"maxPagingSize":     32,
	})

	apiCfg := viper.Sub("api")

	API.DefaultPagingSize = apiCfg.GetInt("defaultPagingSize")
	API.MaxPagingSize = apiCfg.GetInt("maxPagingSize")
}

func SaveSetting() {
	if viper.ConfigFileUsed() == "" {
		if err := viper.WriteConfigAs(UsedConfigFile); err != nil {
			log.Error("Cannot save as setting to file", zap.String("filepath", UsedConfigFile), zap.Error(err))
		}
	} else if err := viper.WriteConfig(); err != nil {
		log.Error("Cannot save setting to file", zap.String("filepath", viper.ConfigFileUsed()), zap.Error(err))
	}
}

func NewServices() {
	newLogService()
	newDBService()
	newCacheService()
	newSessionService()
	newMailService()
	newService()
}
