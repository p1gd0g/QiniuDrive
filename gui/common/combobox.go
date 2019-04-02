package common

import "github.com/p1gd0g/ui"

// NewCombobox creates a combobox with some strings.
func NewCombobox(strs ...string) *ui.Combobox {
	combobox := ui.NewCombobox()
	for _, v := range strs {
		combobox.Append(v)
	}
	return combobox
}
