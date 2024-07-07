package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	tmpFilePath string
)

func init() {
	tmpFilePath = os.Getenv("TMP_FILE_PATH")
	if tmpFilePath == "" {
		tmpFilePath = "./tmpFiles"
	}
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/file", func(c *gin.Context) {
		namespace, pod, container, path := c.DefaultQuery("namespace", ""), c.DefaultQuery("pod", ""), c.DefaultQuery("container", ""), c.DefaultQuery("path", "")

		tmpFilePath := fmt.Sprintf("%s/%s/%s", tmpFilePath, uuid.New(), filepath.Base(path))

		// kubectl cp <some-namespace>/<some-pod>:/tmp/foo /tmp/bar
		b, err := exec.Command("kubectl", "cp", fmt.Sprintf("%s/%s:%s", namespace, pod, path), tmpFilePath, "-c", container).CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("kubectl cp exec error: %w %s", err, b))
		}
		defer os.Remove(tmpFilePath)

		c.File(tmpFilePath)
	})

	r.POST("/file", func(c *gin.Context) {
		namespace, pod, container, path := c.DefaultQuery("namespace", ""), c.DefaultQuery("pod", ""), c.DefaultQuery("container", ""), c.DefaultQuery("path", "")

		file, err := c.FormFile("file")
		if err != nil {
			panic(err)
		}

		tmpFilePath := fmt.Sprintf("%s/%s/%s", tmpFilePath, uuid.New(), filepath.Base(path))
		if err := c.SaveUploadedFile(file, tmpFilePath); err != nil {
			panic(err)
		}

		// kubectl cp /tmp/foo <some-namespace>/<some-pod>:/tmp/bar
		b, err := exec.Command("kubectl", "cp", tmpFilePath, fmt.Sprintf("%s/%s:%s", namespace, pod, path), "-c", container).CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("kubectl cp exec error: %w %s", err, b))
		}
		defer os.Remove(tmpFilePath)

		c.String(http.StatusCreated, "Uploaded")
	})

	r.Run()
}
