package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
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

	file := model.File{}
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		respondwithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	err = api.repo.Create(&file)

	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, "Successfully Created")
	}
}

func (p *FileApi) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondwithJSON(w, http.StatusBadRequest, "File id missing")
		return
	}
	payload, err := p.repo.GetByID(id)

	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, payload)
	}
}

func (p *FileApi) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload, err := p.repo.GetByID(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		f, err := os.Open(fmt.Sprintf(".files/%s.%s", payload.ID, payload.Extension))
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

func (p *FileApi) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondwithJSON(w, http.StatusBadRequest, "File id missing")
		return
	}
	_, err := p.repo.Delete(id)
	if err != nil {
		respondwithJSON(w, http.StatusInternalServerError, "Server Error"+err.Error())
	} else {
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
