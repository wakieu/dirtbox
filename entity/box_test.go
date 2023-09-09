package entity

import "testing"

func TestBox_UpdateText(t *testing.T) {
	box := NewBox("/foo")
	newText := "bar baz"
	box.UpdateText(newText)
	if box.Text != newText {
		t.Errorf("text not updated correctly!")
	}
}

func TestBox_IsEmpty(t *testing.T) {
	box := NewBox("/foo")
	empty := box.IsEmpty()
	if !empty {
		t.Errorf("wrong is empty value!")
	}
}