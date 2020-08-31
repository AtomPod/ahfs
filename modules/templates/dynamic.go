package templates

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/unknwon/com"
	"go.uber.org/zap"
)

var (
	bodyTemplates = template.New("")
)

func Mailer() *template.Template {
	staticDir := filepath.Join(setting.StaticRootPath, "templates", "mail")

	if com.IsDir(staticDir) {
		files, err := com.StatDir(staticDir)
		if err != nil {
			log.Warn("Failed to read templates", zap.String("dir", staticDir), zap.Error(err))
		} else {
			for _, file := range files {
				if !strings.HasSuffix(file, ".tpl") {
					continue
				}

				content, err := ioutil.ReadFile(filepath.Join(staticDir, file))
				if err != nil {
					log.Warn("Failed to read custom templates", zap.String("filename", file), zap.Error(err))
					continue
				}

				name := strings.TrimSuffix(file, ".tpl")
				if _, err := bodyTemplates.New(name).Parse(string(content)); err != nil {
					log.Warn("Failed to parse templates", zap.String("name", name), zap.Error(err))
				} else {
					log.Debug("parse template content", zap.String("name", file))
				}
			}
		}
	}
	return bodyTemplates
}
