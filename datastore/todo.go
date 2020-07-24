package datastore

import "time"

type ToDo struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Attachments    []string  `json:"attachments"`
	Status         string    `json:"status"`
	CreationDate   time.Time `json:"creation_date"`
	DueDate        time.Time `json:"due_date"`
	CompletionDate time.Time `json:"completion_date"`
}
