package routers

import (
	"crypto/subtle"
	"net/http"

	"github.com/czhj/ahfs/modules/setting"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Metrics(ctx *gin.Context) {
	if setting.Metrics.Token == "" {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	header := ctx.GetHeader("Authorization")

	got := []byte(header)
	want := []byte("Bearer " + setting.Metrics.Token)

	if subtle.ConstantTimeCompare(got, want) != 1 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
}
