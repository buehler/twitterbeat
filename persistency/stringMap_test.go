package persistency

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := NewStringMap()
	if m == nil {
		t.Errorf("Constructor returned nil")
	}
}

func TestLoad(t *testing.T) {
	m := NewStringMap()
	m.Load("")

	if m.fileName != "persistentStringMap.json" {
		t.Errorf("Wrong fileName! Exprected 'persistentStringMap.json' got: '%s'", m.fileName)
	}
}

func TestLoadFromFile(t *testing.T) {
	m := NewStringMap()
	m.Load("./stringMap_test.json")

	data, ok := m.list["hello"]

	if !ok {
		t.Errorf("Value hello not in list")
	}

	if data != "world" {
		t.Errorf("Wrong data provided! Exprected 'world' got: '%s'", data)
	}

	if m.fileName != "./stringMap_test.json" {
		t.Errorf("Wrong fileName! Exprected './stringMap_test.json' got: '%s'", m.fileName)
	}
}

func TestGet(t *testing.T) {
	m := NewStringMap()
	m.Load("./stringMap_test.json")

	data := m.Get("hello")

	if data != "world" {
		t.Errorf("Wrong data provided! Exprected 'world' got: '%s'", data)
	}
}

func TestSet(t *testing.T) {
	m := NewStringMap()
	m.Load("./stringMap_test.json")

	m.Set("data", "myvalue")

	data := m.Get("data")
	if data != "myvalue" {
		t.Errorf("Wrong data provided! Exprected 'myvalue' got: '%s'", data)
	}
}

func TestDelete(t *testing.T) {
	m := NewStringMap()
	m.Load("./stringMap_test.json")

	contains := m.Contains("data")

	if !contains {
		t.Errorf("Data not persisted")
	}

	m.Delete("data")

	contains = m.Contains("data")

	if contains {
		t.Errorf("Data not deleted")
	}
}
