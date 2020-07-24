package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
	"todo/datastore"
)

var ds = datastore.FileStore{}

func Init() error {
	err := ds.Initialize()
	return err
}

func errBadRequest(w http.ResponseWriter, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(msg)
}

func errNotFound(w http.ResponseWriter, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(msg)
}

// @Tags Create Task Handler
// Create Task Handler  godoc
// @Summary create kubernetes cluster
// @Description create kubernetes cluster on cloud nodes
// @Param cloud			path	string	true	"public cloud name e.g aws, gcp"
// @Param project_id	path	string	true	"id of the project"
// @Param	X-Profile-Id	header	string	false	"{X-Profile-Id}"
// @Param X-Auth-Token  header  string  false    "X-Auth-Token"
// @Produce json
// @Success 200 "{"status": "kubernetes creation in progress for project {name}"}"
// @Failure 400 "{"error": "error msg", "description" : "error description"}"
// @router /create/{cloud}/{project_id} [post]
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Invalid Request Body"}`))
		return
	}

	var todo datastore.ToDo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Invalid Request Body"}`))
		return
	}

	todo.CreationDate = time.Now()
	out, err := ds.WriteRecord(todo)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Error creating todo item"}`))
		return
	}

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["todoId"]
	if !ok {
		errBadRequest(w, []byte(`{"message": "item id not provided"}`))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Invalid Request Body"}`))
		return
	}

	var todo datastore.ToDo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Invalid Request Body"}`))
		return
	}

	todo.ID = id
	out, err := ds.UpdateRecord(todo)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Error updating todo item"}`))
		return
	}

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["todoId"]
	if !ok {
		errBadRequest(w, []byte(`{"message": "item id not provided"}`))
	}

	err := ds.DeleteRecord(id)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Error deleting todo item"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
func GetListHandler(w http.ResponseWriter, r *http.Request) {

	pageNum := 0
	limitNum := 0
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page != "" {
		t, err := strconv.Atoi(page)
		if err == nil {
			pageNum = t
		}
	}

	if limit != "" {
		t, err := strconv.Atoi(limit)
		if err == nil {
			limitNum = t
		}
	}

	out, err := ds.ReadAllRecord()
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Error creating todo item"}`))
		return
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreationDate.After(out[j].CreationDate)
	})

	start := pageNum * limitNum
	end := limitNum + start
	if end >= len(out) {
		pageNum = 0
		limitNum = 5
		start = pageNum * limitNum
		end = limitNum + start
	}

	out = out[start:end]

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["todoId"]
	if !ok {
		errBadRequest(w, []byte(`{"message": "item id not provided"}`))
	}

	out, err := ds.ReadRecord(id)
	if err != nil {
		errBadRequest(w, []byte(`{"message": "Error retrieving todo item"}`))
		return
	}

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
