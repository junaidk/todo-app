package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"syscall"
	"time"
	"todo/config"
	"todo/reports"
)

const (
	CountTasksKind        = "CountTasksKind"
	CalculateAvgKind      = "CalculateAvgKind"
	CountMaxCompletedKind = "CountMaxCompletedKind"
	CountMaxCreatedKind   = "CountMaxCreatedKind"
)

func writeResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// CountTasks godoc
// @Summary Count of total tasks, completed tasks, and remaining tasks
// @Description Count of total tasks, completed tasks, and remaining tasks
// @Tags Reports
// @Accept  json
// @Produce  json
// @Success 200 {object} reports.CountTaskReport
// @Failure 500 {object} HTTPError
// @Router /count-tasks [get]
func CountTasks(w http.ResponseWriter, r *http.Request) {

	cData, ok := cacheRead(CountTasksKind)
	if ok {
		writeResponse(w, cData)
		return
	}
	out, err := reports.CountTasks()

	if err != nil {
		NewError(w, 500, err)
		return
	}

	resp, _ := json.Marshal(out)
	cacheWrite(CountTasksKind, resp)
	writeResponse(w, resp)

}

// CalculateAvg godoc
// @Summary Average number of tasks completed per day
// @Description Average number of tasks completed per day
// @Tags Reports
// @Accept  json
// @Produce  json
// @Success 200 {object} CalculateAvgResp
// @Failure 500 {object} HTTPError
// @Router /calculate-avg [get]
func CalculateAvg(w http.ResponseWriter, r *http.Request) {

	cData, ok := cacheRead(CalculateAvgKind)
	if ok {
		writeResponse(w, cData)
		return
	}
	out, err := reports.CalculateAvg()

	if err != nil {
		NewError(w, 500, err)
		return
	}

	obj := CalculateAvgResp{AvgPerDay: out}

	resp, _ := json.Marshal(obj)
	cacheWrite(CalculateAvgKind, resp)
	writeResponse(w, resp)
}

type CalculateAvgResp struct {
	AvgPerDay float32 `json:"avg_per_day"`
}

// CountMaxCompleted godoc
// @Summary maximum number of tasks were completed in a single day
// @Description maximum number of tasks were completed in a single day
// @Tags Reports
// @Accept  json
// @Produce  json
// @Success 200 {object} CountMaxCompletedResp
// @Failure 500 {object} HTTPError
// @Router /count-max-completed [get]
func CountMaxCompleted(w http.ResponseWriter, r *http.Request) {

	cData, ok := cacheRead(CountMaxCompletedKind)
	if ok {
		writeResponse(w, cData)
		return
	}

	out, err := reports.CountMaxCompletedTasks()

	if err != nil {
		NewError(w, 500, err)
		return
	}

	obj := CountMaxCompletedResp{MaxCompleteDate: out}

	resp, _ := json.Marshal(obj)
	cacheWrite(CountMaxCompletedKind, resp)
	writeResponse(w, resp)
}

type CountMaxCompletedResp struct {
	MaxCompleteDate string `json:"max_complete_date"`
}

// CountMaxCreated godoc
// @Summary Count maximum number of tasks added on a particular day
// @Description Count maximum number of tasks added on a particular day
// @Tags Reports
// @Accept  json
// @Produce  json
// @Success 200 {object} reports.MaxTask
// @Failure 500 {object} HTTPError
// @Router /count-max-created [get]
func CountMaxCreated(w http.ResponseWriter, r *http.Request) {
	cData, ok := cacheRead(CountMaxCreatedKind)
	if ok {
		writeResponse(w, cData)
		return
	}
	out, err := reports.CountMaxCreatedTasks()

	if err != nil {
		NewError(w, 500, err)
		return
	}

	resp, _ := json.Marshal(out)
	cacheWrite(CountMaxCreatedKind, resp)
	writeResponse(w, resp)
}

func cacheRead(cacheKind string) ([]byte, bool) {

	fPath := path.Join(config.CacheDirPath, cacheKind+".json")
	fInfo, err := os.Stat(fPath)

	if os.IsNotExist(err) {
		return nil, false
	}

	statT := fInfo.Sys().(*syscall.Stat_t)
	creationTime := timespecToTime(statT.Ctim)

	if time.Now().Sub(creationTime) > (15 * time.Minute) {
		return nil, false
	}

	out, err := ioutil.ReadFile(fPath)

	if err != nil {
		return nil, false
	}

	return out, true
}

// https://topic.alibabacloud.com/a/how-to-use-golang-to-get-the-creationmodification-time-of-files-on-linux_1_16_30132202.html
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func cacheWrite(cacheKind string, data []byte) error {

	fPath := path.Join(config.CacheDirPath, cacheKind+".json")

	err := ioutil.WriteFile(fPath, data, 0755)

	return err
}
