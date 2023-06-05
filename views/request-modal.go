package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var requestModal *tview.Modal

func GetRequestModal() *tview.Modal {
	requestModal = tview.NewModal()

	requestModal.SetTitle("Extra response info")
	requestModal.SetTitleAlign(tview.AlignCenter)
	requestModal.SetBorder(true)
	requestModal.SetBackgroundColor(tcell.ColorBlack)
	requestModal.SetTextColor(tcell.ColorWhiteSmoke)

	requestModal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.SetFocus(responseBodyView)
		}
		return event
	})

  return requestModal
}
