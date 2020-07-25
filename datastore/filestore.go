package datastore

import (
	"encoding/json"
	"errors"
	"github.com/teris-io/shortid"
	"io/ioutil"
	"os"
	"sync"
	"todo/config"
)

//var sid *shortid.Shortid
//var mux sync.Mutex

type FileStore struct {
	dbFile string
	sid    *shortid.Shortid
	mux    sync.Mutex
}

func (f *FileStore) Initialize() error {

	f.dbFile = config.DataStorePath

	_, err := os.Stat(f.dbFile)

	if os.IsNotExist(err) {
		file, err := os.Create(f.dbFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return err
	}

	f.sid = sid
	//f.mux = mux

	return nil
}

func (f *FileStore) WriteRecord(do ToDo) (ToDo, error) {

	data := make(map[string]ToDo)

	f.mux.Lock()
	defer f.mux.Unlock()
	dbFileData, err := ioutil.ReadFile(f.dbFile)
	if err != nil {
		return ToDo{}, err
	}
	if len(dbFileData) != 0 {
		err = json.Unmarshal(dbFileData, &data)
		if err != nil {
			return ToDo{}, err
		}
	}
	if do.ID == "" {
		uid := f.getUniqueId()
		do.ID = uid
	}

	_, ok := data[do.ID]
	if ok {
		return do, nil
	}
	data[do.ID] = do
	jOut, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return ToDo{}, err
	}
	ioutil.WriteFile(f.dbFile, jOut, 0755)

	return do, nil
}

func (f *FileStore) ReadRecord(id string) (ToDo, error) {

	data := make(map[string]ToDo)

	f.mux.Lock()
	defer f.mux.Unlock()
	dbFileData, err := ioutil.ReadFile(f.dbFile)
	if err != nil {
		return ToDo{}, err
	}
	if len(dbFileData) != 0 {
		err = json.Unmarshal(dbFileData, &data)
		if err != nil {
			return ToDo{}, err
		}
	}

	out, ok := data[id]
	if !ok {
		return ToDo{}, errors.New("No Record Found for id: " + id)
	}

	return out, nil

}

func (f *FileStore) ReadAllRecord() ([]ToDo, error) {

	data := make(map[string]ToDo)

	f.mux.Lock()
	defer f.mux.Unlock()
	dbFileData, err := ioutil.ReadFile(f.dbFile)
	if err != nil {
		return []ToDo{}, err
	}
	if len(dbFileData) != 0 {
		err = json.Unmarshal(dbFileData, &data)
		if err != nil {
			return []ToDo{}, err
		}
	}

	var out []ToDo

	for k, v := range data {
		v.ID = k
		out = append(out, v)
	}

	return out, nil

}

func (f *FileStore) UpdateRecord(do ToDo) (ToDo, error) {

	data := make(map[string]ToDo)

	f.mux.Lock()
	defer f.mux.Unlock()
	dbFileData, err := ioutil.ReadFile(f.dbFile)
	if err != nil {
		return ToDo{}, err
	}
	if len(dbFileData) != 0 {
		err = json.Unmarshal(dbFileData, &data)
		if err != nil {
			return ToDo{}, err
		}
	}

	_, ok := data[do.ID]
	if !ok {
		return ToDo{}, errors.New("No Record Found for id: " + do.ID)
	}
	data[do.ID] = do
	jOut, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return ToDo{}, err
	}
	ioutil.WriteFile(f.dbFile, jOut, 0755)

	return do, nil
}
func (f *FileStore) DeleteRecord(id string) error {
	data := make(map[string]ToDo)

	f.mux.Lock()
	defer f.mux.Unlock()
	dbFileData, err := ioutil.ReadFile(f.dbFile)
	if err != nil {
		return err
	}
	if len(dbFileData) != 0 {
		err = json.Unmarshal(dbFileData, &data)
		if err != nil {
			return err
		}
	}

	_, ok := data[id]
	if !ok {
		return errors.New("No Record Found for id: " + id)
	}
	delete(data, id)
	jOut, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(f.dbFile, jOut, 0755)

	return nil
}

func (f *FileStore) getUniqueId() string {
	id, _ := f.sid.Generate()
	return id
}
