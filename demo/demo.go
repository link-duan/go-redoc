package main

import (
	"github.com/gin-gonic/gin"
	goredoc "github.com/link-duan/go-redoc"

	_ "embed"
)

//go:embed doc.json
var doc string

func main() {
	r := gin.New()
	r.GET("/swagger/*any", goredoc.GinHandler(&goredoc.Setting{
		OpenAPIJson: doc,
		UriPrefix:   "/swagger",
		Title:       "Go Redoc",
	}))
	r.Run(":8888")
}
