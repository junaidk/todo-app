package datastore

import (
	"log"
	"testing"
	"time"
)

func TestFileStore_WriteRecord(t *testing.T) {

	fs := FileStore{}

	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	id := "testid"

	todo := ToDo{
		ID:             id,
		Title:          "abv",
		Description:    "asdasd",
		Attachments:    []string{"asdasd"},
		Status:         "done",
		CreationDate:   time.Now(),
		DueDate:        time.Now(),
		CompletionDate: time.Now(),
	}

	todo, err := fs.WriteRecord(todo)
	if err != nil {
		t.Errorf("error was not expected while writing: %s", err.Error())
	}

	log.Printf("%+v", todo)
}

func TestFileStore_ReadRecord(t *testing.T) {

	id := "testid"

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	out, err := fs.ReadRecord(id)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	log.Printf("%+v", out)
}

func TestFileStore_UpdateRecord(t *testing.T) {
	id := "testid"

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	todo := ToDo{
		ID:             id,
		Title:          "updated",
		Description:    "updated",
		Attachments:    nil,
		Status:         "",
		CreationDate:   time.Time{},
		DueDate:        time.Time{},
		CompletionDate: time.Time{},
	}
	out, err := fs.UpdateRecord(todo)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	log.Printf("%+v", out)
}

func TestFileStore_DeleteRecord(t *testing.T) {
	id := "testid"

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	err := fs.DeleteRecord(id)
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

}

func TestFileStore_ReadAllRecord(t *testing.T) {

	fs := FileStore{}
	if err := fs.Initialize(); err != nil {
		t.Errorf("error was not expected while init: %s", err.Error())
	}

	out, err := fs.ReadAllRecord()
	if err != nil {
		t.Errorf("error was not expected while reading: %s", err.Error())
	}

	log.Printf("%+v", out)

}
