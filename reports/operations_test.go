package reports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
	"todo/config"
	"todo/datastore"
)

func SeedItems(in time.Time) map[string]datastore.ToDo {

	seedTime := in
	data := map[string]datastore.ToDo{
		"1": {
			ID:             "1",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(24 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"2": {
			ID:             "1",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(48 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"3": {
			ID:             "1",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(24 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(24 * time.Hour),
		},
		"4": {
			ID:             "1",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(48 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(24 * time.Hour),
		},
		"5": {
			ID:             "1",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(24 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(48 * time.Hour),
		},
		"6": {
			ID:             "1",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(48 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"7": {
			ID:             "1",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(72 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(24 * time.Hour),
		},
		"8": {
			ID:             "1",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(72 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"9": {
			ID:             "1",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(96 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(96 * time.Hour),
		},
	}

	return data
}

func TestMain(m *testing.M) {

	data := SeedItems(time.Now())
	config.DataStorePath = "/tmp/dat.json"

	out, _ := json.Marshal(data)
	os.Remove(config.DataStorePath)
	err := ioutil.WriteFile(config.DataStorePath, out, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	os.Exit(m.Run())
}
func TestCountTasks(t *testing.T) {

	out, err := CountTasks()

	if err != nil {
		t.Errorf("error was not expected : %s", err.Error())
	}
	if out.Completed != 4 {
		t.Errorf("Completed count is not correct: %d != %d", out.Completed, 4)
	}
	if out.Total != 9 {
		t.Errorf("Total count is not correct: %d != %d", out.Total, 9)
	}
	if out.Remaining != 5 {
		t.Errorf("Remaining count is not equal: %d != %d", out.Remaining, 5)
	}

}

func TestCalculateAvg(t *testing.T) {
	var out float32
	out, err := CalculateAvg()
	if err != nil {
		t.Errorf("error was not expected : %s", err.Error())
	}
	expected := float32(1.33)
	if fmt.Sprintf("%.2f", out) != fmt.Sprintf("%.2f", 1.33) {
		t.Errorf("Avg is not correct: %.2f != %.2f", out, expected)
	}
}

func TestCountMaxCompletedTasks(t *testing.T) {
	out, err := CountMaxCompletedTasks()
	if err != nil {
		t.Errorf("error was not expected : %s", err.Error())
	}
	expected := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	if out != expected {
		t.Errorf("Day is not correct: %s != %s", out, expected)
	}
}

func TestCountMaxCreatedTasks(t *testing.T) {
	out, err := CountMaxCreatedTasks()
	if err != nil {
		t.Errorf("error was not expected : %s", err.Error())
	}

	if len(out) != 2 {
		t.Errorf("max count is not correct : %d != %d", 2, len(out))
	}
}
