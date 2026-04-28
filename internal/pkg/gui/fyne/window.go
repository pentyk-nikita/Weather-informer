package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	guisettings "github.com/pentyk-nikita/Weather-informer/internal/domain/gui_settings"
)

type window struct {
	w             fyne.Window
	temperatureLabel *widget.Label
}

func NewW(w fyne.Window) *window {
	return &window{w: w}
}

func (win *window) Resize(ws guisettings.WindowSize) error {
	if !ws.IsFull() {
		win.w.Resize(fyne.NewSize(float32(ws.Width()), float32(ws.Height())))
	}
	return nil
}

func (win *window) UpdateTemperature(t float32) error {
	if win.temperatureLabel != nil {
		win.temperatureLabel.SetText("Температура: " + formatTemperature(t))
	}
	return nil
}

func (win *window) SetTemperatureWidget(tw guisettings.TextWidget) error {

	content := win.w.Content()
	if content == nil {
		win.w.SetContent(tw.Render().(fyne.CanvasObject))
	} else {
		win.w.SetContent(tw.Render().(fyne.CanvasObject))
	}
	return nil
}

func (win *window) Render() error {
	return nil
}

func formatTemperature(t float32) string {
	return string(rune(t)) 
}