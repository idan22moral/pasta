package server

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

func RunServer(addr string, uploadsDir string) error {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(res, req)
			return
		}

		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			res.WriteHeader(http.StatusAccepted)
			return
		}
		res.Write(content)
	})

	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	http.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if len(req.MultipartForm.File["files"]) == 0 {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		uploadDirName := fmt.Sprintf("%d", time.Now().Unix())
		uploadDirPath, err := filepath.Abs(path.Join(uploadsDir, uploadDirName))

		if err != nil {
			fmt.Println(err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = os.MkdirAll(uploadDirPath, 0777)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, file := range req.MultipartForm.File["files"] {
			f, err := file.Open()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			uploadFileName := path.Join(uploadDirPath, file.Filename)
			uploadFile, err := os.Create(uploadFileName)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			defer uploadFile.Close()

			_, err = io.Copy(uploadFile, f)

			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	return http.ListenAndServe(addr, nil)
}
