package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range crudRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	for _, route := range reportRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var crudRoutes = Routes{
	Route{
		"Create ToDo",
		"POST",
		"/todo",
		CreateTaskHandler,
	},
	Route{
		"Get List",
		"GET",
		"/todo",
		GetListHandler,
	},
	Route{
		"Get Todo",
		"GET",
		"/todo/{todoId}",
		GetTaskHandler,
	},
	Route{
		"Edit ToDo",
		"PUT",
		"/todo/{todoId}",
		UpdateTaskHandler,
	},
	Route{
		"Delete ToDo",
		"DELETE",
		"/todo/{todoId}",
		DeleteTaskHandler,
	},
}

var reportRoutes = Routes{
	Route{
		"Count Tasks",
		"GET",
		"/count-tasks",
		CountTasks,
	},
	Route{
		"Calculate Avg",
		"GET",
		"/calculate-avg",
		CalculateAvg,
	},
	Route{
		"Count Max Completed Tasks",
		"GET",
		"/count-max-completed",
		CountMaxCompleted,
	},
	Route{
		"Count Max Created Tasks",
		"GET",
		"/count-max-created",
		CountMaxCreated,
	},
	Route{
		"Find Similar Task ",
		"GET",
		"/find-similar-task",
		FindSimilarTask,
	},
}
