package views

import "github.com/rivo/tview"

var app = tview.NewApplication()

func StartApp() error {
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
	mainFlex.AddItem(requestModal, 0, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlF {
			app.SetFocus(requestModal)
		}
		return event
	})

  return app.SetRoot(mainFlex, true).SetFocus(form).Run()
}
