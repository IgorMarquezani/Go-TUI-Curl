package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var responseBodyView *tview.TextView

func GetResponseBody() *tview.TextView {
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
			app.SetFocus(requestModal)
		}
		return event
	})

  return responseBodyView
}
