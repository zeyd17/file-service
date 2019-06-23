package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/zeyd17/file-microservice/model"
	"github.com/zeyd17/file-microservice/repository"
)

type FileApi struct {
	repo repository.IFileRepo
}

func NewFileApi(repo repository.IFileRepo) *FileApi {
	return &FileApi{repo: repo}
}

// Create a new File
func (api *FileApi) Post(w http.ResponseWriter, r *http.Request) {

	fileModel := model.File{}

	r.ParseMultipartForm(32 << 20)
	tmpfile, handler, err := r.FormFile("file")
	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tmpfile.Close()

	fileModel.ID = strings.ToLower(uuid.New().String())
	fileModel.Name = handler.Filename
	fileModel.Size = handler.Size
	fileModel.Format = handler.Header.Get("Content-Type")
	extensions := strings.Split(fileModel.Name, ".")
	ext := ""
	if len(extensions) > 0 {
		ext = extensions[len(extensions)-1]
	}
	fileModel.Extension = ext

	fileName := fmt.Sprintf("./files/%s.%s", fileModel.ID, fileModel.Extension)

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	io.Copy(file, tmpfile)
	fmt.Println(fileModel)
	err = api.repo.Create(&fileModel)

	if err != nil {

		respondwithJSON(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, fileModel)
	}

}

func (api *FileApi) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondwithJSON(w, http.StatusBadRequest, "File id missing")
		return
	}
	payload, err := api.repo.GetByID(id)

	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, payload)
	}
}

func (api *FileApi) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload, err := api.repo.GetByID(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		f, err := os.Open(fmt.Sprintf("./files/%s.%s", payload.ID, payload.Extension))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", payload.Name))
		w.Header().Set("Content-Type", payload.Format)

		io.Copy(w, f)

		w.WriteHeader(http.StatusOK)
		respondwithJSON(w, http.StatusOK, payload)
	}
}

func (api *FileApi) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondwithJSON(w, http.StatusBadRequest, "File id missing")
		return
	}

	payload, _ := api.repo.GetByID(id)

	err := api.repo.Delete(id)
	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, "Server Error"+err.Error())
	} else {
		os.Remove(fmt.Sprintf("./files/%s.%s", payload.ID, payload.Extension))
		respondwithJSON(w, http.StatusOK, "Delete Successfully")
	}
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	isSuccess := false
	if code == 200 {
		isSuccess = true
	}
	response, _ := json.Marshal(model.Result{IsSuccess: isSuccess, Data: payload})
	w.Write(response)
}
