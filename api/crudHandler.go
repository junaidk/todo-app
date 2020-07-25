package api

import (
	"encoding/json"
	"errors"
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

// CreateTodo godoc
// @Summary Creates a Task item
// @Description creates a Task item
// @Tags Crud
// @Accept  json
// @Produce  json
// @Param account body datastore.ToDo true "Creates a Task"
// @Success 200 {object} datastore.ToDo
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /todo [post]
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewError(w, 400, err)
		return
	}

	var todo datastore.ToDo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		NewError(w, 400, err)
		return
	}

	todo.CreationDate = time.Now()
	out, err := ds.WriteRecord(todo)
	if err != nil {
		NewError(w, 500, err)
		return
	}

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// UpdateTodo godoc
// @Summary Updates a ToDo item
// @Description updates a ToDo item
// @Tags Crud
// @Accept  json
// @Produce  json
// @Param  todoId path string true "ToDO task ID"
// @Success 200 {object} datastore.ToDo
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /todo/{todoId} [update]
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["todoId"]
	if !ok {
		NewError(w, 400, errors.New("item id not provided"))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewError(w, 400, err)
		return
	}

	var todo datastore.ToDo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		NewError(w, 400, err)
		return
	}

	todo.ID = id
	out, err := ds.UpdateRecord(todo)
	if err != nil {
		NewError(w, 500, err)
		return
	}

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// DeleteTodo godoc
// @Summary Deletes a ToDo item
// @Description delete a ToDo item
// @Tags Crud
// @Accept  json
// @Produce  json
// @Param  todoId path string true "ToDO task ID"
// @Success 200 {object} datastore.ToDo
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /todo/{todoId} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["todoId"]
	if !ok {
		NewError(w, 400, errors.New("item id not provided"))
	}

	err := ds.DeleteRecord(id)
	if err != nil {
		NewError(w, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

// GetTodoList godoc
// @Summary Get List of ToDo items
// @Description get List of ToDo items
// @Tags Crud
// @Accept  json
// @Produce  json
// @Success 200 {array} datastore.ToDo
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /todo [get]
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
		NewError(w, 500, err)
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

// GetTodo godoc
// @Summary Get ToDo item
// @Description get List of ToDo items
// @Tags Crud
// @Accept  json
// @Produce  json
// @Param  todoId path string true "ToDO task ID"
// @Success 200 {object} datastore.ToDo
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /todo/{todoId} [get]
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["todoId"]
	if !ok {
		NewError(w, 400, errors.New("item id not provided"))
		return
	}

	out, err := ds.ReadRecord(id)
	if err != nil {
		NewError(w, 500, err)
		return
	}

	resp, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
