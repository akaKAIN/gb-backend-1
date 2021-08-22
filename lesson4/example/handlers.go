package example

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// HandlerRoot Handler for base URI
type HandlerRoot struct{}

func (h *HandlerRoot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.FormValue("name")
		if _, err := fmt.Fprintf(w, "Get query 'name': %s\n", name); err != nil {
			http.Error(w, "Response error", http.StatusBadRequest)
		}
	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			user := new(User)
			if err := json.NewDecoder(r.Body).Decode(user); err != nil {
				http.Error(w, "Wrong user data\n", http.StatusBadRequest)
				return
			}
			if _, err := fmt.Fprintf(w, "new User:\nname: %s\nwallet: %v\n", user.Name, user.Wallet); err != nil {
				http.Error(w, "Response error", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unexpected content-type", http.StatusBadRequest)
			return
		}

	}
}

// HandlerUpload Handler for upload files
type HandlerUpload struct {
	dir string
}

func (h *HandlerUpload) UploadPath(fileName, uniqStr string) string {
	return fmt.Sprintf("%s%v%s-%s", h.dir, string(os.PathSeparator), uniqStr, fileName)
}

func (h *HandlerUpload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileReader, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file in request", http.StatusBadRequest)
		return
	}
	defer fileReader.Close()
	var uniqStr string
	uniq := time.Now().UnixNano()

	uniqStr = strconv.Itoa(int(uniq))
	filePath := h.UploadPath(fileHeader.Filename, uniqStr)

	data, err := ioutil.ReadAll(fileReader)
	if !isDirExist(h.dir) {
		if err := CreateDir(h.dir); err != nil {
			http.Error(w, "Can't init dir path", http.StatusInternalServerError)
		}
	}
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, err = fmt.Fprintf(w, "file '%s' was uploaded\n", fileHeader.Filename); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type HandlerUploadList struct {
	dir string
}

func (h *HandlerUploadList) UploadDir() string {
	return h.dir
}

func (h *HandlerUploadList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, err := os.ReadDir(h.UploadDir())
	if err != nil {
		http.Error(w, "Error of upload read", http.StatusInternalServerError)
		return
	}
	uploads := new(Uploads)
	for _, file := range dir {
		if !file.IsDir() {
			fi, err := file.Info()
			if err != nil {
				continue
			}
			upload := new(Upload)
			upload.Name = fi.Name()
			upload.Size = fi.Size()

			uploads.Uploads = append(uploads.Uploads, *upload)
		}
	}

	// Filter by query
	filter := r.FormValue("filter")
	switch filter {
	case "jpg":
	case "jpeg":
	case "png":
		uploads = FilterUploadBy(uploads, filter)
	default:
		http.Error(w, "Wrong filter value", http.StatusBadRequest)
		return
	}

	// Send json
	if err = json.NewEncoder(w).Encode(uploads); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FilterUploadBy(uploads *Uploads, suffix string) *Uploads {
	filteredUploads := make([]Upload, 0)
	for _, upload := range (*uploads).Uploads {
		if filepath.Ext(upload.Name) == "."+suffix {
			filteredUploads = append(filteredUploads, upload)
		}
	}

	uploads.Uploads = filteredUploads
	return uploads
}
