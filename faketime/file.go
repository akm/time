package faketime

import (
	"os"

	"github.com/akm/time"
)

type File struct {
	FilePath string
	layout   string
}

func NewFile(filePath string, layout string) *File {
	return &File{
		FilePath: filePath,
		layout:   layout,
	}
}

func (f *File) Save(t time.Time) (rerr error) {
	content := t.Format(f.layout)
	file, err := os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		rerr = file.Close()
	}()
	if _, err = file.WriteString(content); err != nil {
		rerr = err
		return
	}
	return
}

func (f *File) Delete() error {
	if err := os.Remove(f.FilePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
