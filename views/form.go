package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var form *tview.Form

func GetForm() *tview.Form {
  form = tview.NewForm()

	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	form.SetLabelColor(tcell.ColorWhiteSmoke)
	form.SetButtonTextColor(tcell.ColorWhiteSmoke)
	form.SetButtonBackgroundColor(tcell.ColorBlack)
	form.SetButtonsAlign(tview.AlignCenter)

  return form
}
