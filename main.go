package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/gif"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MainWindow struct {
	app        fyne.App
	window     fyne.Window
	chkR       []RuneCheckboxes
	chkP       []string
	rw         []RW
	rwFiltered []RW
	checkListR []string
	checkListP []string
	db         *gorm.DB
}

type RuneCheckboxes struct {
	rc   string
	chck *widget.Check
}

type R struct {
	Id    uint
	Name  string
	Image []byte
}

type RW struct {
	Id          uint
	Name        string
	ItemClass   string
	Ladder      bool
	Sockets     uint8
	LVL         uint8
	Runes       string
	Description string
}

func NewMW(db *gorm.DB) *MainWindow {
	mw := &MainWindow{
		app:  app.New(),
		chkR: make([]RuneCheckboxes, 0),
		chkP: []string{"CR", "FHR", "FCR", "+Skills"},
		db:   db,
	}
	mw.window = mw.app.NewWindow("D2Runewords")
	mw.window.Resize(fyne.NewSize(700, 480))
	mw.window.SetFixedSize(true)
	mw.window.CenterOnScreen()

	icon, _ := fyne.LoadResourceFromPath("icon.png")
	mw.app.SetIcon(icon)

	err := mw.getRuneWords(&mw.rw)
	if err != nil {
		fmt.Println(err)
	}

	mw.ui()

	return mw
}

// connects to the database
func connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error while connect to sqlite:", err)
		return nil, err
	}
	return db, err
}

func (mw *MainWindow) getRunes(r *[]R) error {
	rows := mw.db.Find(r)
	if rows.Error != nil {
		fmt.Println("No rows:", rows.Error)
		return rows.Error
	}
	return nil
}

func (mw *MainWindow) setRunes() *fyne.Container {
	var r []R
	err := mw.getRunes(&r)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rContainer := container.New(layout.NewGridLayout(3))
	for _, v := range r {
		c := widget.NewCheck(v.Name, func(checked bool) { mw.checker(v.Name, checked, &mw.checkListR) })
		image := canvas.NewImageFromResource(
			fyne.NewStaticResource(v.Name+".gif", v.Image),
		)
		image.FillMode = canvas.ImageFillOriginal
		mw.chkR = append(mw.chkR, RuneCheckboxes{v.Name, c})
		rContainer.Add(container.NewHBox(c, layout.NewSpacer(), image))
	}
	return rContainer
}

func (mw *MainWindow) getRuneWords(rw *[]RW) error {
	rows := mw.db.Find(rw)
	if rows.Error != nil {
		fmt.Println("No rows:", rows.Error)
		return rows.Error
	}
	return nil
}

// without function gives an error
func (mw *MainWindow) decodeImage(b []byte) (*gif.GIF, error) {
	image, err := gif.DecodeAll(bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return image, nil
}

func (mw *MainWindow) showAllRW(c *fyne.Container, data *[]RW) {
	var desc, ls string
	c.RemoveAll()

	for _, v := range *data {
		desc = strings.ReplaceAll(v.Description, "  ", "\n")

		lna := canvas.NewText(v.Name, color.RGBA{R: 160, G: 144, B: 96, A: 255})
		lna.TextSize = 24
		lna.Alignment = fyne.TextAlignCenter

		lic := canvas.NewText(fmt.Sprintf("Items class: %s", v.ItemClass), color.White)
		lic.Alignment = fyne.TextAlignCenter

		lso := canvas.NewText(fmt.Sprintf("Sockets: %d", v.Sockets), color.White)
		lso.Alignment = fyne.TextAlignCenter

		llv := canvas.NewText(fmt.Sprintf("Level: %d", v.LVL), color.White)
		llv.Alignment = fyne.TextAlignCenter

		lru := canvas.NewText(fmt.Sprintf("Runes: %s", v.Runes), color.White)
		lru.Alignment = fyne.TextAlignCenter

		if v.Ladder {
			ls = fmt.Sprintf("Only Ladder\n%s\n", desc)
		} else {
			ls = fmt.Sprintf("%s\n", desc)
		}
		lde := widget.NewLabel(ls)
		lde.Wrapping = fyne.TextWrapWord
		lde.Alignment = fyne.TextAlignCenter
		lde.TextStyle = fyne.TextStyle{Italic: true}

		c.Add(lna)
		c.Add(lic)
		c.Add(lso)
		c.Add(llv)
		c.Add(lru)
		c.Add(lde)
	}
}

// make slice with active checkboxes
func (mw *MainWindow) checker(cName string, checked bool, cl *[]string) {
	var temp []string

	switch cName {
	case "CR":
		cName = "Crushing Blow"
	case "FCR":
		cName = "Faster Cast Rate"
	case "FHR":
		cName = "Faster Hit Recovery"
	case "+Skills":
		cName = "Skills"
	}

	if checked {
		*cl = append(*cl, cName)
	} else {
		for _, v := range *cl {
			if v != cName {
				temp = append(temp, v)
			}
		}
		*cl = temp
	}
}

func (mw *MainWindow) findRW(c *fyne.Container, listR *[]string, listP *[]string) {
	var temp []string
	var str, desc, ls string
	c.RemoveAll()

	for _, v := range *listR {
		if len(*listR) > 0 {
			temp = append(temp, fmt.Sprintf("runes LIKE \"%%%s%%\"", v))
		}
	}
	for _, v := range *listP {
		if len(*listP) > 0 {
			temp = append(temp, fmt.Sprintf("description LIKE \"%%%s%%\"", v))
		}
	}
	
	str = strings.Join(temp, " AND ")
	rows := mw.db.Where(str).Find(&mw.rwFiltered)

	if rows.Error != nil {
		fmt.Println("No rows:", rows.Error)
	}
	
	for _, v := range mw.rwFiltered {
		desc = strings.ReplaceAll(v.Description, "  ", "\n")

		lna := canvas.NewText(v.Name, color.RGBA{R: 160, G: 144, B: 96, A: 255})
		lna.TextSize = 24
		lna.Alignment = fyne.TextAlignCenter

		lic := canvas.NewText(fmt.Sprintf("Items class: %s", v.ItemClass), color.White)
		lic.Alignment = fyne.TextAlignCenter

		lso := canvas.NewText(fmt.Sprintf("Sockets: %d", v.Sockets), color.White)
		lso.Alignment = fyne.TextAlignCenter

		llv := canvas.NewText(fmt.Sprintf("Level: %d", v.LVL), color.White)
		llv.Alignment = fyne.TextAlignCenter

		lru := canvas.NewText(fmt.Sprintf("Runes: %s", v.Runes), color.White)
		lru.Alignment = fyne.TextAlignCenter

		if v.Ladder {
			ls = fmt.Sprintf("Only Ladder\n%s\n", desc)
		} else {
			ls = fmt.Sprintf("%s\n", desc)
		}
		lde := widget.NewLabel(ls)
		lde.Wrapping = fyne.TextWrapWord
		lde.Alignment = fyne.TextAlignCenter
		lde.TextStyle = fyne.TextStyle{Italic: true}

		c.Add(lna)
		c.Add(lic)
		c.Add(lso)
		c.Add(llv)
		c.Add(lru)
		c.Add(lde)
	}
}

func (mw *MainWindow) ui() {
	l1 := canvas.NewLine(color.Black)
	l2 := canvas.NewLine(color.Black)

	// container with checkboxes for ["CR", "FCR", "FHR", "+Skills"]
	propContainer := container.New(layout.NewGridLayout(4))
	for _, v := range mw.chkP {
		cp := widget.NewCheck(v, func(checked bool) { mw.checker(v, checked, &mw.checkListP) })
		propContainer.Add(cp)
	}
	// container with all checkboxes
	fContainer := container.NewVBox(
		mw.setRunes(),
		l1,
		propContainer,
	)
	// scroll
	contvScroll := container.New(layout.NewVBoxLayout())
	vScroll := container.NewVScroll(contvScroll)
	vScroll.SetMinSize(fyne.NewSize(360.00, 430.00))
	// buttons
	fButton := widget.NewButton("All runewords", func() { mw.showAllRW(contvScroll, &mw.rw) })
	sButton := widget.NewButton("Search", func() { mw.findRW(contvScroll, &mw.checkListR, &mw.checkListP) })
	buttContainer := container.NewHBox(
		fButton,
		layout.NewSpacer(),
		sButton,
	)
	// container with data of runewords
	sContainer := container.NewVBox(vScroll)
	sContainer.Add(buttContainer)

	mainContainer := container.NewHBox(fContainer, l2, sContainer)
	mw.window.SetContent(mainContainer)
}

func (mw *MainWindow) Run() {
	mw.window.ShowAndRun()
}

func main() {
	db, err := connect()
	if err != nil {
		return
	}

	mw := NewMW(db)
	mw.Run()
}
