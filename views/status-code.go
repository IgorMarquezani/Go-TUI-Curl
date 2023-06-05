package views

import "github.com/rivo/tview"

var statusCode *tview.TextView

func GetStatusCode() *tview.TextView {
	statusCode = tview.NewTextView()

	statusCode.SetTextAlign(tview.AlignCenter)
	statusCode.SetText("Status Code")
	statusCode.SetDynamicColors(true)

  return statusCode
}
