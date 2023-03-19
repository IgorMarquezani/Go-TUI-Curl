package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var MethodOptions = []string{"GET", "POST", "UPDATE", "PATCH", "DELETE"}

var CookieNums int = 1

type RequestForm struct {
	Url     string
	Method  string
	Body    []byte
	Request http.Request
	Cookies []string
}

var app = tview.NewApplication()

var requestForm RequestForm

var responseFlex *tview.Flex

var statusCodeView *tview.TextView

var responseBodyView *tview.TextView

var form *tview.Form

func SetCookieByField(field string, cookies *[]string) http.Cookie {
	cookie := http.Cookie{}

	fields := strings.Split(field, " ")

	for i := range fields {
		var name []byte
		var data []byte
		var swap bool
		for j := 0; j < len(fields[i])-1; j++ {
			if fields[i][j] == '=' && !swap {
				swap = true
				continue
			}
			if swap {
				data = append(data, byte(fields[i][j]))
			}
			name = append(name, byte(fields[i][j]))
		}
	}

	return cookie
}

func StartRequestForm() *tview.Form {
	requestForm.Cookies = make([]string, 5)
	form := tview.NewForm()

	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	form.SetLabelColor(tcell.ColorWhiteSmoke)
	form.SetButtonTextColor(tcell.ColorWhiteSmoke)
	form.SetButtonBackgroundColor(tcell.ColorBlack)
	form.SetButtonsAlign(tview.AlignCenter)

	// Item 1
	form.AddInputField("URL:", "", 80, nil, func(text string) {
		requestForm.Url = text
	})

	// Item 2
	dropDown := tview.NewDropDown()
	dropDown.SetFieldBackgroundColor(tcell.ColorBlack)
	dropDown.SetLabel("Methods")
	dropDown.SetOptions(MethodOptions, func(option string, optionIndex int) {
		requestForm.Method = option
	})
	dropDown.SetCurrentOption(0)

	form.AddFormItem(dropDown)
	// Item 3
	form.AddTextArea("Body:", "", 80, 16, 80*16, func(text string) {
		requestForm.Body = []byte(text)
	})

	// Buttons does not count as items
	form.AddButton("Make Request", func() {
		body := bytes.NewBuffer(requestForm.Body)
		request, err := http.NewRequest(requestForm.Method, requestForm.Url, body)
		if err != nil {
			panic(err)
		}

		for i := 0; i < CookieNums-1; i++ {
			cookie := SetCookieByField(requestForm.Cookies[i], &requestForm.Cookies)
			request.AddCookie(&cookie)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		bodyArray, _ := ioutil.ReadAll(response.Body)
		statusCodeView.SetText(response.Status)
		responseBodyView.SetText(string(bodyArray))
		app.SetFocus(responseBodyView)
	})

	form.AddButton("Add Cookie", func() {
		cookieIndex := CookieNums - 1
		form.AddTextArea("Cookie "+strconv.Itoa(CookieNums), "", 80, 5, 80*5, func(text string) {
			i := cookieIndex
			requestForm.Cookies[i] = text
		})
		CookieNums++
	})

	form.AddButton("Remove Cookie", func() {
		if CookieNums > 1 {
			index := form.GetFormItemIndex("Cookie " + strconv.Itoa(CookieNums-1))
			requestForm.Cookies = requestForm.Cookies[:len(requestForm.Cookies)-1]
			form.RemoveFormItem(index)
			CookieNums--
		}
	})

	form.AddButton("Exit", func() {
		app.Stop()
	})

	return form
}

func main() {
	responseFlex = tview.NewFlex().SetDirection(tview.FlexRow)
	responseFlex.SetBorder(true)

	statusCodeView = tview.NewTextView()
	statusCodeView.SetBorder(true).SetTitle("Status Code")

	responseBodyView = tview.NewTextView()
	responseBodyView.SetScrollable(true)
	responseBodyView.SetBorder(true).SetTitle("Response Body")
	responseBodyView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlL {
			app.SetFocus(form)
		}
		return event
	})

	form = StartRequestForm()
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR {
			app.SetFocus(responseBodyView)
		}
		return event
	})

	app = tview.NewApplication()

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(form, 0, 1, false).
			AddItem(responseFlex.
				AddItem(statusCodeView, 3, 1, false).
				AddItem(responseBodyView, 39, 1, false),
				0, 1, false),
			0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
