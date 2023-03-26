// Code generated by hertz generator.

package main

import (
	"log"

	"github.com/hertz-contrib/gzip"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"

	"simplified-tik-tok/biz/dal"
	"simplified-tik-tok/biz/mw"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dal.Init()
	mw.InitJwt()
	h := server.Default()
	h.Use(
		gzip.Gzip(
			gzip.DefaultCompression,
			gzip.WithExcludedExtensions([]string{".pdf", ".mp4"}),
		),
	)

	register(h)
	h.Spin()
}
