package tui

import (
	"fmt"
	"strconv"

	"net/url"

	"github.com/gdamore/tcell/v2"
	"github.com/nojima/httpie-go/exchange"
	"github.com/nojima/httpie-go/input"
	"github.com/rivo/tview"
)

func inputHandler(app *tview.Application) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		the_rune := event.Rune()
		if the_rune == 'q' {
			app.Stop()
		}
		if the_rune == 't' {
			cb := func() {
				// now := time.Now()
				// entryDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				// res, _ := editor.CaptureInputFromEditor(editor.GetPreferredEditorFromEnvironment)
				// content := string(res[:])
				// entry := newEntry(-1, entryDate.Unix(), content)
				// saveEntry(db, *entry)
				// save in db
			}

			app.Suspend(cb)
		}
		return event
	}
}

type Query struct {
	Cmd string
}

type AppOptions struct {
	// type Input struct {
	// Method     Method
	// URL        *url.URL
	// Parameters []Field
	// Header     Header
	// Body       Body
	// }
	input input.Input
	// type Options struct {
	// JSON      bool
	// Form      bool
	// ReadStdin bool
	// }
	inputOptions input.Options
	// type Options struct {
	// Timeout         time.Duration
	// FollowRedirects bool
	// Auth            AuthOptions
	// SkipVerify      bool
	// ForceHTTP1      bool
	// }
	exchangeOptions exchange.Options
}

func Init() {
	app := tview.NewApplication()
	app.SetInputCapture(inputHandler(app))
	pages := tview.NewPages()

	txtView := tview.NewTextView()
	txtView.SetBorder(true)

	pages.AddPage("preview", txtView, true, false)

	list := tview.NewList()
	list.SetTitle("queries")
	list.SetBorder(true)
	pages.AddPage("query list", list, true, true)

	// type Field struct {
	// Name   string
	// Value  string
	// IsFile bool
	// }
	var p []input.Field
	// type Input struct {
	// Method     Method
	// URL        *url.URL
	// Parameters []Field
	// Header     Header
	// Body       Body
	// }
	// type Body struct {
	// BodyType      BodyType
	// Fields        []Field
	// RawJSONFields []Field // used only when BodyType == JSONBody
	// Files         []Field // used only when BodyType == FormBody
	// Raw           []byte  // used only when BodyType == RawBody
	// }
	// type BodyType int

	// const (
	// EmptyBody BodyType = iota
	// JSONBody
	// FormBody
	// RawBody
	// )
	b := input.Body{BodyType: 0}
	u, _ := url.Parse("https://www.reddit.com/r/newmoto.json")
	i := input.Input{Method: "GET", URL: u, Parameters: p, Body: b}
	// type Options struct {
	// Timeout         time.Duration
	// FollowRedirects bool
	// Auth            AuthOptions
	// SkipVerify      bool
	// ForceHTTP1      bool
	// }
	exchangeOptions := exchange.Options{}
	req, e := exchange.BuildHTTPRequest(&i, &exchangeOptions)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(req.Host)
	q := Query{Cmd: ""}
	queries := []Query{q}

	for j := 1; j <= 12; j++ {
		list.AddItem(strconv.Itoa(j), "", rune(0), nil)
	}

	listViewInputHandler := func(event *tcell.EventKey) *tcell.EventKey {
		key_rune := event.Rune()
		if key_rune == 'l' {
			// index := list.GetCurrentItem()
			// entry := entries[index]
			// txtView.SetText(entry.content)
			// date := getEntryDateString(&entry)
			// txtView.SetTitle(date)
			// pages.SwitchToPage("preview")
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
			list.SetCurrentItem(len(queries) - 1)
		}
		return event
	}
	list.SetInputCapture(listViewInputHandler)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
