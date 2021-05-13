package journal

import (
	"database/sql"
	"fmt"
	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
	"gopkg.in/ini.v1"
	// "strconv"
	"github.com/curioussavage/journal-go/editor"
	"time"
)

type config struct {
	Journal_dir          string `ini:"journal_dir"`
	Editor               string `ini:"editor"`
	File_ext             string `ini:"file_ext"`
	Default_ls_day_range int    `ini:"default_ls_day_range"`
}

func initConfig() *config {
	c := new(config)
	err := ini.MapTo(c, "/home/awinn/.config/com.github.curioussavage.journal/config.ini")
	if err != nil {
		fmt.Println("error loading config")
	}
	return c
}

type Entry struct {
	id      int
	time    int64
	content string
}

func newEntry(id int, time int64, content string) *Entry {
	e := Entry{time: time, content: content}
	return &e
}

func (entry Entry) isExisting() bool {
	return entry.id != -1
}

func getEntryDateString(entry *Entry) string {
	return time.Unix(entry.time, 0).Format("2006-01-02")
}

func initDB(db_path string) *sql.DB {
	db, _ := sql.Open("sqlite3", db_path)
	db.Exec("create table if not exists entries (id INTEGER PRIMARY KEY, date INT, content TEXT)")
	return db
}

func getEntries(db *sql.DB /*, page int*/) []Entry {
	// now := time.Now().AddDate(0, 0, -days).Unix()
	rows, err := db.Query("SELECT * FROM entries ORDER BY date DESC")
	if err != nil {
		fmt.Println("problem with query")
	}
	entries := []Entry{}
	for rows.Next() {
		var tempEntry Entry
		err =
			rows.Scan(&tempEntry.id, &tempEntry.time, &tempEntry.content)
		if err != nil {
			fmt.Println("error scanning row")
		}
		entries = append(entries, tempEntry)
	}
	if err := rows.Err(); err != nil {
		// log.Fatal(err)
		fmt.Println("foo")
	}
	return entries
}

func saveEntry(db *sql.DB, entry Entry) {
	if entry.isExisting() {
		_, err := db.Query("UPDATE entries SET content = ? WHERE id = ?", entry.content, entry.id)
		if err != nil {
			fmt.Println("problem with query")
		}
	} else {
		fmt.Println("entry content")
		fmt.Println(entry.content)
		_, err := db.Query("INSERT INTO entries (date, content) VALUES (?, ?)", entry.time, entry.content)
		if err != nil {
			fmt.Println("problem with query")
		}
	}

}

func GetEntry(db *sql.DB, id int) Entry {
	var entry Entry
	rows, _ := db.Query("SELECT * FROM entries WHERE Id = ?", id)
	for rows.Next() {
		rows.Scan(&entry.id, &entry.time, &entry.content)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("error")
	}
	return entry
}

func onEntrySelected(pages *tview.Pages, txt *tview.TextView, entry Entry) func() {
	return func() {
		pages.SwitchToPage("preview")
		txt.SetText(entry.content)
		date := getEntryDateString(&entry)
		txt.SetTitle(date)
	}
}

func inputHandler(app *tview.Application, db *sql.DB) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		the_rune := event.Rune()
		if the_rune == 'q' {
			app.Stop()
		}
		if the_rune == 't' {
			cb := func() {
				now := time.Now()
				entryDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				res, _ := editor.CaptureInputFromEditor(editor.GetPreferredEditorFromEnvironment)
				content := string(res[:])
				entry := newEntry(-1, entryDate.Unix(), content)
				saveEntry(db, *entry)
				// save in db
			}

			app.Suspend(cb)
		}
		return event
	}
}

// editEntry launches editor to change an entry

// init starts the tui app
func Init() {
	config := *initConfig()
	db := initDB(config.Journal_dir)
	db.SetMaxOpenConns(1)
	app := tview.NewApplication()
	app.SetInputCapture(inputHandler(app, db))
	pages := tview.NewPages()

	entries := getEntries(db)
	txtView := tview.NewTextView()
	txtView.SetBorder(true)
	pages.AddPage("preview", txtView, true, false)
	list := tview.NewList()
	pages.AddPage("list", list, true, true)
	list.SetTitle("Journal")
	list.SetBorder(true)

	monthList := tview.NewList()
	pages.AddPage("monthlist", list, true, true)
	for j := 1; j <= 12; j++ {
		monthList.AddItem(time.Month(j).String(), "", rune(0), nil)
	}

	txtViewInputHandler := func(event *tcell.EventKey) *tcell.EventKey {
		key_rune := event.Rune()
		if key_rune == 'h' {
			pages.SwitchToPage("list")
		} else if key_rune == 'k' {
			index := list.GetCurrentItem()
			if index != 0 {
				list.SetCurrentItem(index - 1)
				entry := entries[index-1]
				txtView.SetText(entry.content)
				date := getEntryDateString(&entry)
				txtView.SetTitle(date)
			}
		} else if key_rune == 'j' {
			index := list.GetCurrentItem()
			if index != len(entries)-1 {
				list.SetCurrentItem(index + 1)
				entry := entries[index+1]
				txtView.SetText(entry.content)
				date := getEntryDateString(&entry)
				txtView.SetTitle(date)
			}
		}
		return nil
	}
	txtView.SetInputCapture(txtViewInputHandler)

	list.ShowSecondaryText(false)
	list.SetWrapAround(false)
	listViewInputHandler := func(event *tcell.EventKey) *tcell.EventKey {
		key_rune := event.Rune()
		if key_rune == 'l' {
			index := list.GetCurrentItem()
			entry := entries[index]
			txtView.SetText(entry.content)
			date := getEntryDateString(&entry)
			txtView.SetTitle(date)
			pages.SwitchToPage("preview")
			return nil
		} else if key_rune == 'j' {
			index := list.GetCurrentItem()
			list.SetCurrentItem(index + 1)
			return nil
		} else if key_rune == 'k' {
			index := list.GetCurrentItem()
			if index != 0 {
				list.SetCurrentItem(index - 1)
				return nil
			}
		} else if key_rune == 'G' {
			list.SetCurrentItem(len(entries) - 1)
		}
		return event
	}
	list.SetInputCapture(listViewInputHandler)

	for _, entry := range entries {
		handler := onEntrySelected(pages, txtView, entry)
		list.AddItem(getEntryDateString(&entry), string(entry.id), rune(0), handler)
	}

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
