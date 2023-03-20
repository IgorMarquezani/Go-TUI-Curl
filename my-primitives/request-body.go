package myprimitives

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewRequestBody() *tview.Box {
  RequestBody := tview.NewBox()

  RequestBody.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

    return x, y, width, height
  })

  return RequestBody
}
