package go_redoc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Setting struct {
	// OpenAPI JSON definition content
	OpenAPIJson string

	// Site prefix
	UriPrefix string

	// Title of doc site
	Title string

	// Redoc options. https://github.com/Redocly/redoc#redoc-options-object
	// example: { "json-sample-expand-level": "all" }
	RedocOptions map[string]string
}

func (s *Setting) normalize() {
	if s.Title == "" {
		s.Title = "Redoc"
	}
	if s.RedocOptions == nil {
		s.RedocOptions = make(map[string]string)
	}
	s.RedocOptions["spec-url"] = s.UriPrefix + "/doc.json"
}

const htmlTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>%s</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc %s></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
  </body>
</html>`

func GinHandler(setting *Setting) gin.HandlerFunc {
	setting.normalize()
	controller := &ginController{setting: setting}

	return func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, setting.UriPrefix)
		handler, ok := map[string]gin.HandlerFunc{
			"/doc.json":   controller.doc,
			"/index.html": controller.index,
			"/":           controller.index,
			"":            controller.index,
		}[path]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		handler(c)
	}
}

type ginController struct {
	setting *Setting
}

func (c *ginController) index(ctx *gin.Context) {
	options := ""
	for k, v := range c.setting.RedocOptions {
		options += fmt.Sprintf(" %s=\"%s\"", k, v)
	}

	content := fmt.Sprintf(
		htmlTemplate,
		c.setting.Title,
		options,
	)
	ctx.Data(http.StatusOK, "text/html", []byte(content))
}

func (c *ginController) doc(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "application/json", []byte(c.setting.OpenAPIJson))
}
