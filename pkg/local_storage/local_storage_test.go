package localstorage

import "testing"

func TestLocalStorage(t *testing.T) {
	localStorage := NewLocalStorage()
	localStorage.Set("test", "hello world")
	if localStorage.Get("test") != "hello world" {
		t.Log("localstore item isn't hello world", localStorage.Get("test"))
		t.Fail()
	}
}
