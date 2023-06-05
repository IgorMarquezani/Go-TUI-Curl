package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var MethodOptions = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

var CookieNums int = 1

var HowManyCookies int

type CookieOptions map[string]func(value string, cookie *http.Cookie)

type RequestForm struct {
	Url          string
	Method       string
	Body         []byte
	Request      http.Request
	Cookies      []string
	LastResponse http.Response
}

var app = tview.NewApplication()

var requestForm RequestForm

var responseFlex *tview.Flex

var statusCodeView *tview.TextView

var responseBodyView *tview.TextView

var form *tview.Form

var modal *tview.Modal

func GetCookieOptions() CookieOptions {
	options := make(CookieOptions)

	options["Expires"] = func(value string, cookie *http.Cookie) {
		expire, err := time.Parse(time.RFC1123, value)
		if err != nil {
			return
		}
		cookie.Expires = expire
	}

	options["Path"] = func(value string, cookie *http.Cookie) {
		cookie.Path = value
	}

	options["HttpOnly"] = func(value string, cookie *http.Cookie) {
		strTrue := [4]byte{'t', 'r', 'u', 'e'}
		for i, c := range strTrue {
			if value[i] == c || value[i] == c+32 {
				continue
			} else {
				return
			}
		}
		cookie.HttpOnly = true
	}

	options["Secure"] = func(value string, cookie *http.Cookie) {
		strTrue := [4]byte{'t', 'r', 'u', 'e'}
		for i, c := range strTrue {
			if value[i] == c || value[i] == c+32 {
				continue
			} else {
				return
			}
		}
		cookie.Secure = true
	}

	options["Domain"] = func(value string, cookie *http.Cookie) {
		cookie.Domain = value
	}

	return options
}

func SetCookies(text string) http.Cookie {
	var (
		already       bool
		cookie        http.Cookie
		cookieOptions = GetCookieOptions()
		fields        = strings.Split(text, ";")
	)

	for _, field := range fields {
		infos := strings.Split(field, "=")
		if len(infos) < 2 {
			continue
		}

		key, value := infos[0], infos[1]
		if key[0] == '\n' || key[0] == ' ' {
			key = key[1:]
		}

		assign, ok := cookieOptions[key]
		if ok {
			assign(value, &cookie)
		}

		if !ok && !already {
			cookie.Name = key
			cookie.Value = value
			already = true
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
	HeaderTextArea := tview.NewTextArea()
	HeaderTextArea.SetLabel("Headers")
	HeaderTextArea.SetSize(4, 80)

	// Item 4
	BodyTextArea := tview.NewTextArea()
	BodyTextArea.SetLabel("Body:")
	BodyTextArea.SetSize(16, 80)
	BodyTextArea.SetMaxLength(80 * 16)
	BodyTextArea.SetChangedFunc(func() {
		requestForm.Body = []byte(BodyTextArea.GetText())
	})
	form.AddFormItem(BodyTextArea)

	// Buttons does not count as items
	form.AddButton("Make Request", func() {
		body := bytes.NewBuffer(requestForm.Body)

		request, err := http.NewRequest(requestForm.Method, requestForm.Url, body)
		if err != nil {
			responseBodyView.SetText("Error: " + err.Error())
			return
		}

		for i := 0; i < HowManyCookies; i++ {
			cookie := SetCookies(requestForm.Cookies[i])
			request.AddCookie(&cookie)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			responseBodyView.SetText("[red]Error: " + err.Error())
			return
		}

		bodyArray, err := ioutil.ReadAll(response.Body)
		if err != nil {
			responseBodyView.SetText("[red]Error: " + err.Error())
			return
		}

		headerValues := make([]byte, 0)
		for name, values := range response.Header {
			headerValues = append(headerValues, []byte(name+": ")...)
			for _, value := range values {
				headerValues = append(headerValues, []byte(value+" ")...)
			}
			headerValues = append(headerValues, '\n')
		}

		var cookies []byte
		for i := range response.Cookies() {
			num := strconv.Itoa(i + 1)
			cookies = append(cookies, "Cookie "+num+": "+response.Cookies()[i].String()...)
			cookies = append(cookies, '\n')
		}

		statusCodeView.SetText(response.Status)
		responseBodyView.SetText("[yellow]" + string(bodyArray))
		modal.SetText("Cookies:\n" + string(cookies) + "Headers:\n" + string(headerValues) + "Testing")
		app.SetFocus(responseBodyView)
	})

	form.AddButton("Add Cookie", func() {
		if HowManyCookies < 5 {
			cookieIndex := HowManyCookies
			cookieField := tview.NewTextArea()
			cookieField.SetLabel("Cookie " + strconv.Itoa(HowManyCookies+1))
			cookieField.SetSize(5, 80)
			cookieField.SetChangedFunc(func() {
				i := cookieIndex
				requestForm.Cookies[i] = cookieField.GetText()
			})
			form.AddFormItem(cookieField)
			cookieNum := strconv.Itoa(HowManyCookies)
			cookieField.SetText("Cookie number: "+cookieNum, true)
			HowManyCookies++
		}
	})

	form.AddButton("Remove Cookie", func() {
		if HowManyCookies > 0 {
			index := form.GetFormItemIndex("Cookie " + strconv.Itoa(HowManyCookies))
			requestForm.Cookies[HowManyCookies-1] = ""
			form.RemoveFormItem(index)
			HowManyCookies--
		}
	})

	form.AddButton("Exit", func() {
		app.Stop()
	})

	return form
}

func main() {
	modal = tview.NewModal()
	modal.SetTitle("Extra response info")
	modal.SetTitleAlign(tview.AlignCenter)
	modal.SetBorder(true)
	modal.SetBackgroundColor(tcell.ColorBlack)
	modal.SetTextColor(tcell.ColorWhiteSmoke)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.SetFocus(responseBodyView)
		}
		return event
	})

	statusCodeView = tview.NewTextView()
	statusCodeView.SetTextAlign(tview.AlignCenter)
	statusCodeView.SetText("Status Code")
	statusCodeView.SetDynamicColors(true)

	responseBodyView = tview.NewTextView()
	responseBodyView.SetScrollable(true)
	responseBodyView.SetBorder(true)
	responseBodyView.SetBorderAttributes(tcell.AttrDim)
	responseBodyView.SetDynamicColors(true)
	responseBodyView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlL || event.Key() == tcell.KeyCtrlA {
			app.SetFocus(form)
		}

		if event.Rune() == 'f' || event.Rune() == 'F' {
			app.SetFocus(modal)
		}

		return event
	})

	responseFlex = tview.NewFlex().SetDirection(tview.FlexRow)
	responseFlex.AddItem(statusCodeView, 1, 1, false)
	responseFlex.AddItem(responseBodyView, 0, 1, false)

	helpBar := tview.NewTextView()
	helpBar.SetDynamicColors(true)
	helpBar.SetText("Type F for more informations like cookies and response header")
	helpBar.SetTextAlign(tview.AlignLeft)

	form = StartRequestForm()
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR || event.Key() == tcell.KeyCtrlD {
			app.SetFocus(responseBodyView)
		}
		return event
	})

	formFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	formFlex.AddItem(form, 0, 1, false)
	formFlex.AddItem(responseFlex, 0, 1, false)

	subFlex := tview.NewFlex()
	subFlex.SetDirection(tview.FlexRow)
	subFlex.AddItem(formFlex, 0, 1, false)
	subFlex.AddItem(helpBar, 1, 1, false)

	mainFlex := tview.NewFlex()
	mainFlex.SetDirection(tview.FlexRow)
	mainFlex.AddItem(subFlex, 0, 1, false)
	mainFlex.SetBorder(false)
	mainFlex.AddItem(modal, 0, 0, false)

	app = tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlF {
			app.SetFocus(modal)
		}
		return event
	})

	if err := app.SetRoot(mainFlex, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
