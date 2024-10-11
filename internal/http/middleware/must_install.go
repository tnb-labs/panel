package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/internal/data"
)

// MustInstall 确保已安装应用
func MustInstall(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var slugs []string
		if strings.HasPrefix(r.URL.Path, "/api/website") {
			slugs = append(slugs, "openresty")
		} else if strings.HasPrefix(r.URL.Path, "/api/container") {
			slugs = append(slugs, "podman", "docker")
		} else if strings.HasPrefix(r.URL.Path, "/api/apps/") {
			pathArr := strings.Split(r.URL.Path, "/")
			if len(pathArr) < 4 {
				render := chix.NewRender(w)
				render.Status(http.StatusForbidden)
				render.JSON(chix.M{
					"message": "应用不存在",
				})
				return
			}
			slugs = append(slugs, pathArr[3])
		}

		flag := false
		for _, s := range slugs {
			if installed, _ := data.NewAppRepo().IsInstalled("slug = ?", s); installed {
				flag = true
				break
			}
		}
		if !flag && len(slugs) > 0 {
			render := chix.NewRender(w)
			render.Status(http.StatusForbidden)
			render.JSON(chix.M{
				"message": fmt.Sprintf("应用 %s 未安装", slugs),
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
