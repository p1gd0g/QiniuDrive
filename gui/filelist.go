package gui

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/comm"
	"github.com/p1gd0g/QiniuDrive/tool"
	"github.com/p1gd0g/ui"
)

// FileList contains name, type, size and checkbox.
type FileList struct {
	HBox *ui.Box

	name     *ui.Box
	mType    *ui.Box
	size     *ui.Box
	checkbox *ui.Box

	NameList     []string
	CheckboxList []*ui.Checkbox
}

//NewFileList creates a new file list.
func NewFileList() (l *FileList) {

	l = &FileList{}

	l.name = ui.NewVerticalBox()
	l.mType = ui.NewVerticalBox()
	l.size = ui.NewVerticalBox()
	l.checkbox = ui.NewVerticalBox()

	l.HBox = ui.NewHorizontalBox()
	l.HBox.SetPadded(true)

	l.HBox.Append(l.name, true)
	l.HBox.Append(ui.NewHorizontalSeparator(), false)
	l.HBox.Append(l.mType, false)
	l.HBox.Append(ui.NewHorizontalSeparator(), false)
	l.HBox.Append(l.size, false)
	l.HBox.Append(ui.NewHorizontalSeparator(), false)
	l.HBox.Append(l.checkbox, false)

	return
}

// Display the file list using imported data.
func (l *FileList) Display(
	accessKey *ui.Entry,
	secretKey *ui.Entry,
	bucket *ui.Entry) (err error) {

	list, err := comm.Refresh(
		accessKey.Text(), secretKey.Text(), bucket.Text())

	if err != nil {
		return
	}
	log.Println("Displaying the list.")

	l.name.Clear()
	l.mType.Clear()
	l.size.Clear()
	l.checkbox.Clear()
	log.Println("Boxes cleared.")

	l.NameList = []string{}
	l.CheckboxList = []*ui.Checkbox{}
	log.Println("Lists cleared.")

	for _, item := range list {
		l.name.Append(ui.NewLabel(item.Key), true)
		l.NameList = append(l.NameList, item.Key)

		l.mType.Append(ui.NewLabel(item.MimeType), true)
		l.size.Append(
			ui.NewLabel(tool.FormatSize(item.Fsize)), true)

		tempCheckbox := ui.NewCheckbox("")
		l.checkbox.Append(tempCheckbox, true)
		l.CheckboxList = append(l.CheckboxList, tempCheckbox)
	}
	log.Println("Displayed the list.")

	return
}
