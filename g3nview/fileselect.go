package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets"
	"github.com/g3n/engine/math32"
	"os"
	"path/filepath"
	"sort"
)

type FileSelect struct {
	gui.Panel
	path *gui.Label
	list *gui.List
	bok  *gui.Button
	bcan *gui.Button
}

func NewFileSelect(width, height float32) *FileSelect {

	fs := new(FileSelect)
	fs.Panel.Initialize(width, height)
	fs.SetBorders(2, 2, 2, 2)
	fs.SetPaddings(4, 4, 4, 4)
	fs.SetColor(&math32.White)

	// Set vertical box layout for the whole panel
	l := gui.NewVBoxLayout()
	l.SetSpacing(4)
	fs.SetLayout(l)

	// Creates path label
	fs.path = gui.NewLabel("path")
	fs.Add(fs.path)

	// Creates list
	fs.list = gui.NewVList(0, 0)
	fs.list.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 5, AlignH: gui.AlignWidth})
	fs.Add(fs.list)

	// Button container panel
	bc := gui.NewPanel(0, 0)
	bcl := gui.NewHBoxLayout()
	bcl.SetAlignH(gui.AlignWidth)
	bc.SetLayout(bcl)
	bc.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 1, AlignH: gui.AlignWidth})
	fs.Add(bc)

	fs.bok = gui.NewButton("OK")
	fs.bok.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	fs.bok.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fs.Dispatch("OnOK", nil)
	})
	bc.Add(fs.bok)

	fs.bcan = gui.NewButton("Cancel")
	fs.bcan.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	fs.bcan.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fs.Dispatch("OnCancel", nil)
	})

	fs.list.Subscribe(gui.OnChange, func(evname string, ev interface{}) {
		fs.onSelect()

	})

	bc.Add(fs.bcan)

	path, err := os.Getwd()
	if err != nil {
		log.Error("%s", err)
	} else {
		fs.SetPath(path)
	}
	return fs
}

func (fs *FileSelect) SetPath(path string) error {

	// Open path file or dir
	f, err := os.Open(path)
	if err != nil {
		log.Error("%s", err)
		return err
	}
	defer f.Close()

	// Checks if it is a directory
	files, err := f.Readdir(0)
	if err != nil {
		log.Error("%s", err)
		return err
	}
	fs.path.SetText(path)

	// Sort files by name
	sort.Sort(listFileInfo(files))

	// Reads directory contents and loads into the list
	fs.list.Clear()
	// Adds previous directory
	prev := gui.NewImageLabel("..")
	prev.SetIcon(assets.FolderOpen)
	fs.list.Add(prev)
	// Adds directory files
	for i := 0; i < len(files); i++ {
		item := gui.NewImageLabel(files[i].Name())
		if files[i].IsDir() {
			item.SetIcon(assets.FolderOpen)
		} else {
			item.SetIcon(assets.InsertPhoto)
		}
		fs.list.Add(item)
	}
	return nil
}

func (fs *FileSelect) Selected() string {

	sel := fs.list.Selected()[0]
	label := sel.(*gui.ImageLabel)
	text := label.Text()
	return filepath.Join(fs.path.Text(), text)
}

func (fs *FileSelect) onSelect() {

	// Get selected image label and its txt
	sel := fs.list.Selected()[0]
	label := sel.(*gui.ImageLabel)
	text := label.Text()

	// Checks if previous directory
	if text == ".." {
		dir, _ := filepath.Split(fs.path.Text())
		fs.SetPath(filepath.Dir(dir))
		return
	}

	// Checks if it is a directory
	path := filepath.Join(fs.path.Text(), text)
	s, err := os.Stat(path)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if s.IsDir() {
		fs.SetPath(path)
	}
}

// For sorting panels by Z coordinate
type listFileInfo []os.FileInfo

func (fi listFileInfo) Len() int      { return len(fi) }
func (fi listFileInfo) Swap(i, j int) { fi[i], fi[j] = fi[j], fi[i] }
func (fi listFileInfo) Less(i, j int) bool {

	return fi[i].Name() < fi[j].Name()
}
