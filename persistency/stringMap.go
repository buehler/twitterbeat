package persistency

import (
	"github.com/dustin/gojson"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type StringMap struct {
	list     map[string]string
	fileName string
	mux      sync.Mutex
}

func NewStringMap() *StringMap {
	return &StringMap{}
}

func (m *StringMap) Contains(key string) bool {
	_, contains := m.list[key]
	return contains
}

func (m *StringMap) Get(key string) string {
	data, contains := m.list[key]
	if !contains {
		return "" //TODO: errorhandling
	}
	return data
}

func (m *StringMap) Set(key, value string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.list[key] = value
	m.save()
}

func (m *StringMap) Delete(key string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.list, key)
	m.save()
}

func (m *StringMap) Load(fileName string) {
	defer func() {
		m.fileName = fileName
	}()

	if _, err := os.Stat(fileName); os.IsNotExist(err) || fileName == "" {
		m.list = make(map[string]string)
		if fileName == "" {
			fileName = "persistentStringMap.json"
		}
		return
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(data, &m.list)
	if m.list == nil {
		m.list = make(map[string]string)
	}
}

func (m *StringMap) save() {
	err := os.MkdirAll(path.Dir(m.fileName), 0766)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(m.fileName)

	defer f.Close()

	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(m.list)

	if err != nil {
		panic(err)
	}

	_, err = f.Write(data)

	if err != nil {
		panic(err)
	}
}
