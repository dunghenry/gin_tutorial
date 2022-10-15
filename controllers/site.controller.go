package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
)

func GetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "sites/index.tmpl", gin.H{
		"title": "Home Page",
	})
}

func GetUploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sites/upload.tmpl", gin.H{})
}
func PostUpload(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	filename := filepath.Base(file.Filename)
	data := time.Now().Unix()
	namefile := strconv.Itoa(int(data)) + "." + strings.Split(filename, ".")[1]
	if err := c.SaveUploadedFile(file, "public/images/"+namefile); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	c.String(http.StatusOK, "File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email)
}
func GetPostPage(c *gin.Context) {
	c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{})
}
