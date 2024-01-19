package master

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}
func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func HandleStaticFile(f embed.FS, router *gin.Engine) {
	root := EmbedFolder(f, "out")
	router.Use(static.Serve("/", root))
	staticServer := static.Serve("/", root)
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet &&
			!strings.ContainsRune(c.Request.URL.Path, '.') &&
			!strings.HasPrefix(c.Request.URL.Path, "/v1/") {
			if strings.HasSuffix(c.Request.URL.Path, "/") {
				c.Request.URL.Path += "index.html"
				staticServer(c)
				return
			}
			if !strings.HasSuffix(c.Request.URL.Path, ".html") {
				c.Request.URL.Path += ".html"
				staticServer(c)
				return
			}
		}
	})
}
