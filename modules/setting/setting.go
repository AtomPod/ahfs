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

	PasswordComplexity []string
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
	if err := log.AddLogger("console", "console", fmt.Sprintf(`{"level": "fatal", "stderr": true}`)); err != nil {
		panic(err)
	}
	if err := log.New(fmt.Sprintf(`{"encoding": "console"}`)); err != nil {
		panic(err)
	}

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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
		"app_name":            "AHFS",
		"domain":              "localhost",
		"http_addr":           "0.0.0.0",
		"http_port":           "6270",
		"enable_lets_encrypt": false,
		"lets_encrypt_tos":    false,
		"lets_encrypt_dir":    "https",
		"mode":                "debug",
		"avatar_max_width":    256,
		"avatar_max_height":   256,
		"password_complexity": []string{},
	})
	serverCfg := viper.Sub("server")
	AppName = serverCfg.GetString("app_name")

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
	EnableLetsEncrypt = serverCfg.GetBool("enable_lets_encrypt")
	LetsEncryptTOS = serverCfg.GetBool("lets_encrypt_tos")
	if !LetsEncryptTOS {
		log.Warn("Failed to enable Let's Encrypt due to Let's Encrypt TOS not beging accpeted")
		EnableLetsEncrypt = false
	}
	LetsEncryptDir = serverCfg.GetString("lets_encrypt_dir")
	LetsEncryptHost = serverCfg.GetStringSlice("let_encrypt_host")

	Domain = serverCfg.GetString("domain")

	HTTPAddr = serverCfg.GetString("http_addr")
	HTTPPort = serverCfg.GetString("http_port")
	EnableGzip = serverCfg.GetBool("enable_gzip")

	serverCfg.SetDefault("static_root_path", AppWorkPath)
	serverCfg.SetDefault("app_data_path", path.Join(AppWorkPath, "data"))
	StaticRootPath = serverCfg.GetString("static_root_path")
	AppDataPath = serverCfg.GetString("app_data_path")

	viper.SetDefault("server.avatar_upload_path", filepath.Join(AppDataPath, "avatar"))
	AvatarUploadPath = viper.GetString("server.avatar_upload_path")

	AvatarMaxWidth = serverCfg.GetInt("avatar_max_width")
	AvatarMaxHeight = serverCfg.GetInt("avatar_max_height")

	viper.SetDefault("server.file_upload_path", filepath.Join(AppDataPath, "file"))
	FileUploadPath = viper.GetString("server.file_upload_path")

	defaultAppURL := Protocol + "://" + Domain
	if (Protocol == "http" && HTTPPort != "80") || (Protocol == "https" && HTTPPort != "443") {
		defaultAppURL = defaultAppURL + ":" + HTTPPort
	}
	viper.SetDefault("server.root_url", defaultAppURL)
	AppURL = viper.GetString("server.root_url")
	AppURL = strings.TrimSuffix(AppURL, "/") + "/"

	appURL, err := url.Parse(AppURL)
	if err != nil {
		log.Fatal("Invalid RootURL", zap.String("url", AppURL), zap.Error(err))
	}
	AppSubURL = strings.TrimSuffix(appURL.Path, "/")

	PasswordComplexity = serverCfg.GetStringSlice("password_complexity")
	newAPIService()
}

func newAPIService() {
	viper.SetDefault("api", map[string]interface{}{
		"default_paging_size": 16,
		"max_paging_size":     32,
	})

	apiCfg := viper.Sub("api")

	API.DefaultPagingSize = apiCfg.GetInt("default_paging_size")
	API.MaxPagingSize = apiCfg.GetInt("max_paging_size")
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
	newMetricsService()
	newPprofService()
}
