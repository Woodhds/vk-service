package handlers

import (
	"testing"
)

func TestSave(t *testing.T) {
	if client, e := NewClient("http://vk-predict.herokuapp.com:80"); e != nil {
		t.Error(e)
	} else {
		if e := client.SaveMessage(1, 2, "", ""); e != nil {
			t.Error(e)
		}
	}
}
