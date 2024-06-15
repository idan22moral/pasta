package server

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const SESSION_ID_COOKIE_NAME string = "deviceID"

type FileUpload struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

//go:embed static/*
var staticFiles embed.FS

func RunServer(addr string, uploadsDir string) error {
	uploadsDirAbs, err := filepath.Abs(uploadsDir)
	if err != nil {
		return err
	}

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

	http.HandleFunc("/existingUploads", func(res http.ResponseWriter, req *http.Request) {
		var filePaths []FileUpload
		filepath.Walk(uploadsDirAbs, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(uploadsDirAbs, path)
			if err != nil {
				fmt.Print(err)
				return nil
			}

			fileUrl := filepath.Join("/existingUploads/", relPath)
			fileUpload := FileUpload{Name: info.Name(), Path: fileUrl}
			filePaths = append(filePaths, fileUpload)
			return nil
		})

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(filePaths)
	})

	http.Handle("/existingUploads/", http.StripPrefix("/existingUploads/", http.FileServer(http.Dir(uploadsDirAbs))))

	return http.ListenAndServe(addr, nil)
}
