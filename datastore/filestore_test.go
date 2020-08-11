package datastore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
	"todo/config"
)

var seedData map[string]ToDo

func SeedItems(in time.Time) map[string]ToDo {

	seedTime := in
	data := map[string]ToDo{
		"1": {
			ID:             "1",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(24 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"2": {
			ID:             "2",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(48 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"3": {
			ID:             "3",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(24 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(24 * time.Hour),
		},
		"4": {
			ID:             "4",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(48 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(24 * time.Hour),
		},
		"5": {
			ID:             "5",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(24 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(48 * time.Hour),
		},
		"6": {
			ID:             "6",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(48 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"7": {
			ID:             "7",
			Description:    "",
			Status:         "done",
			CreationDate:   seedTime.Add(72 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: seedTime.Add(24 * time.Hour),
		},
		"8": {
			ID:             "8",
			Description:    "",
			Status:         "",
			CreationDate:   seedTime.Add(72 * time.Hour),
			DueDate:        time.Time{},
			CompletionDate: time.Time{},
		},
		"9": {
			ID:             "9",
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

	seedData = SeedItems(time.Now())
	config.DataStorePath = "/tmp/dat.json"

	os.Exit(m.Run())
}

func writeSeedData() {
	out, _ := json.Marshal(seedData)
	err := ioutil.WriteFile(config.DataStorePath, out, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func TestFileStore_WriteRecord(t *testing.T) {
	writeSeedData()
	t.Cleanup(func() {
		os.Remove(config.DataStorePath)
	})

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	id := "10"

	ti := time.Now()
	todo := ToDo{
		ID:             id,
		Title:          "abv",
		Description:    "asdasd",
		Attachments:    []string{"asdasd"},
		Status:         "done",
		CreationDate:   ti,
		DueDate:        ti,
		CompletionDate: ti,
	}

	todo, err := fs.WriteRecord(todo)
	if err != nil {
		t.Errorf("error was not expected while writing: %s", err.Error())
	}

	data, err := ioutil.ReadFile(fs.dbFile)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	var fileData = map[string]ToDo{}
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		t.Errorf("error was not expected while unmarshalling: %s", err.Error())
	}

	todoData, ok := fileData[id]
	if !ok {
		t.Errorf("written data is not as expected")
	}

	value, _ := json.Marshal(todoData)
	expected, _ := json.Marshal(todo)

	if string(expected) != string(value) {
		t.Errorf("written data is not as expected")
	}

}

func TestFileStore_ReadRecord(t *testing.T) {
	writeSeedData()
	t.Cleanup(func() {
		os.Remove(config.DataStorePath)
	})
	id := "1"

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	out, err := fs.ReadRecord(id)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	value, _ := json.Marshal(out)
	expected, _ := json.Marshal(seedData[id])

	if string(expected) != string(value) {
		t.Errorf("written data is not as expected")
	}

}

func TestFileStore_UpdateRecord(t *testing.T) {
	writeSeedData()
	t.Cleanup(func() {
		os.Remove(config.DataStorePath)
	})
	id := "2"

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	todo := seedData[id]
	todo.Description = "Updated"
	_, err := fs.UpdateRecord(todo)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	data, err := ioutil.ReadFile(fs.dbFile)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}
	var fileData = map[string]ToDo{}
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		t.Errorf("error was not expected while unmarshalling: %s", err.Error())
	}

	todoData, ok := fileData[id]
	if !ok {
		t.Errorf("written data is not as expected")
	}

	value, _ := json.Marshal(todoData)
	expected, _ := json.Marshal(todo)

	if string(expected) != string(value) {
		t.Errorf("written data is not as expected")
	}
}

func TestFileStore_DeleteRecord(t *testing.T) {
	writeSeedData()
	t.Cleanup(func() {
		os.Remove(config.DataStorePath)
	})
	id := "1"

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	err := fs.DeleteRecord(id)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	data, err := ioutil.ReadFile(fs.dbFile)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}
	var fileData = map[string]ToDo{}
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		t.Errorf("error was not expected while unmarshalling: %s", err.Error())
	}

	_, ok := fileData[id]
	if ok {
		t.Errorf("item is not deleted from db")
	}

}

func TestFileStore_ReadAllRecord(t *testing.T) {
	writeSeedData()
	t.Cleanup(func() {
		os.Remove(config.DataStorePath)
	})

	out, _ := json.Marshal(seedData)
	defer os.Remove(config.DataStorePath)
	err := ioutil.WriteFile(config.DataStorePath, out, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	_, err = fs.ReadAllRecord()
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	data, err := ioutil.ReadFile(fs.dbFile)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}
	var fileData = map[string]ToDo{}
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		t.Errorf("error was not expected while unmarshalling: %s", err.Error())
	}

	if len(fileData) != len(seedData) {
		t.Errorf("data size is not same %d!=%d", len(fileData), len(seedData))
	}

}
