package datastore

import (
	"time"
)

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

//type JSONTime struct {
//	time.Time
//}
//
//func (t JSONTime) MarshalJSON() ([]byte, error) {
//	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))
//	return []byte(stamp), nil
//}
//func (t *JSONTime) UnmarshalJSON(data []byte) error {
//	var err error
//
//	str := strings.Trim(string(data), "\"")
//	t.Time, err = time.Parse("2006-01-02", str)
//	return err
//}
