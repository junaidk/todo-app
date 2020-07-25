package reports

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"sync"
	"time"
	"todo/config"
	"todo/datastore"
)

type CountTaskReport struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Remaining int `json:"remaining"`
}

func CountTasks() (CountTaskReport, error) {

	fileName := config.DataStorePath

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return CountTaskReport{}, err
	}

	out := make(map[string]datastore.ToDo)

	err = json.Unmarshal(data, &out)
	if err != nil {
		return CountTaskReport{}, err
	}

	totalChan := make(chan datastore.ToDo, len(out))
	completedChan := make(chan datastore.ToDo, len(out))

	totalCountChan := make(chan int)
	completedCountChan := make(chan int)

	// all tasks
	go func() {
		count := 0
		for range totalChan {
			count += 1
		}
		totalCountChan <- count
	}()

	// completed tasks
	go func() {
		count := 0
		for task := range completedChan {
			if task.Status == "done" {
				count += 1
			}
		}
		completedCountChan <- count
	}()

	for _, task := range out {
		totalChan <- task
		completedChan <- task
	}
	close(totalChan)
	close(completedChan)

	total := <-totalCountChan
	totalCompleted := <-completedCountChan
	remaining := total - totalCompleted
	close(totalCountChan)
	close(completedCountChan)

	res := CountTaskReport{
		Total:     total,
		Completed: totalCompleted,
		Remaining: remaining,
	}

	return res, nil

}

func CalculateAvg() (float32, error) {
	fileName := config.DataStorePath

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return -1, err
	}

	out := make(map[string]datastore.ToDo)

	err = json.Unmarshal(data, &out)
	if err != nil {
		return -1, err
	}

	countCompletedChan := make(chan int, len(out))
	countDaysChan := make(chan string, len(out))

	chanMap := make(map[string]chan datastore.ToDo)

	wg := &sync.WaitGroup{}

	for _, task := range out {
		date := task.CompletionDate.Format("2006-01-02")
		dayChan, ok := chanMap[date]
		if !ok {

			dayChan = make(chan datastore.ToDo, len(out))
			chanMap[date] = dayChan
			wg.Add(1)
			go func(dc chan datastore.ToDo) {
				for task := range dc {
					if task.Status == "done" {
						countCompletedChan <- 1
						countDaysChan <- task.CompletionDate.Format("2006-01-02")
					}
				}
				wg.Done()
			}(dayChan)
		}
		dayChan <- task

	}

	for _, tChan := range chanMap {
		close(tChan)
	}

	totalCompletedChan := make(chan int)
	go func() {
		total := 0
		for range countCompletedChan {
			total += 1
		}
		totalCompletedChan <- total
	}()

	totalChan := make(chan int)
	go func() {
		unique := make(map[string]bool)
		for date := range countDaysChan {
			unique[date] = true
		}
		totalChan <- len(unique)
	}()

	wg.Wait()
	close(countCompletedChan)
	close(countDaysChan)

	completed := <-totalCompletedChan
	total := <-totalChan
	avg := float32(completed) / float32(total)
	return avg, nil

}
func CountMaxCompletedTasks() (string, error) {

	fileName := config.DataStorePath

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	out := make(map[string]datastore.ToDo)

	err = json.Unmarshal(data, &out)
	if err != nil {
		return "", err
	}

	maxCount := make(map[time.Time]int)

	for _, v := range out {

		if v.Status == "done" {
			count, ok := maxCount[v.CompletionDate]
			if !ok {
				maxCount[v.CompletionDate] = 1
			}
			maxCount[v.CompletionDate] = count + 1
		}

	}

	var date time.Time
	count := 0
	for k, v := range maxCount {
		if v > count {
			count = v
			date = k
		}
	}

	return date.Format("2006-01-02"), nil

}

type MaxTask struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func CountMaxCreatedTasks() ([]MaxTask, error) {

	fileName := config.DataStorePath

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return []MaxTask{}, err
	}

	out := make(map[string]datastore.ToDo)

	err = json.Unmarshal(data, &out)
	if err != nil {
		return []MaxTask{}, err
	}

	maxCount := make(map[string]int, 0)

	for _, v := range out {
		date := v.CreationDate.Format("2006-01-02")
		count, ok := maxCount[date]
		if !ok {
			maxCount[date] = 1
		}
		maxCount[date] = count + 1
	}

	var maxList []MaxTask
	for k, v := range maxCount {
		m := MaxTask{
			Date:  k,
			Count: v,
		}
		maxList = append(maxList, m)

	}

	sort.Slice(maxList, func(i, j int) bool {
		return maxList[i].Count > maxList[j].Count
	})

	startVal := maxList[0]
	index := 0
	for i, val := range maxList {
		if val.Count != startVal.Count {
			index = i
			break
		}
	}

	return maxList[:index], nil
}
