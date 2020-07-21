package view

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"me/asdqwer/fund-valuation/data"
	"strconv"
)

var Green = color.RGBA{A: 0xff, R: 76, G: 175, B: 80}
var Red = color.RGBA{A: 0xff, R: 244, G: 67, B: 54}
var Black = color.RGBA{A: 0xff, R: 0, G: 0, B: 0}

type RemoveFunc func(index int)

type FundDataItem struct {
	widget.Box
	OnRemove      RemoveFunc
	fundNameLabel *widget.Label
	fundCodeLabel *widget.Label
	fundUdText    *canvas.Text
	index         int
	fundCode      string
}

func (v *FundDataItem) UpdateUdValue(ud string) {
	v.fundUdText.Text = ud
	v.fundUdText.Color = getUdColor(ud)
	v.fundUdText.Refresh()
}

func getUdColor(ud string) color.RGBA {
	udVal, err := strconv.ParseFloat(ud, 32)
	if err == nil {
		if udVal > 0 {
			return Red
		} else if udVal < 0 {
			return Green
		} else {
			return Black
		}
	} else {
		return Black
	}

}

func NewFundItem(data data.FundData) *FundDataItem {
	fundDataItem := &FundDataItem{}
	fundDataItem.ExtendBaseWidget(fundDataItem)
	fundDataItem.Horizontal = true

	fundDataItem.fundCode = data.FundCode

	fundCodeLabel := widget.Label{Text: data.FundCode}
	fundNameLabel := widget.Label{Text: data.FundName}
	fundUdText := canvas.NewText(data.FundUd, getUdColor(data.FundUd))
	removeButton := widget.Button{Icon: theme.DeleteIcon()}

	removeButton.OnTapped = func() {
		if fundDataItem.OnRemove != nil {
			fundDataItem.OnRemove(fundDataItem.index)
		}
	}

	fundDataItem.Append(&fundCodeLabel)
	fundDataItem.Append(&fundNameLabel)
	fundDataItem.Append(fundUdText)
	fundDataItem.Append(layout.NewSpacer())
	fundDataItem.Append(&removeButton)

	fundDataItem.fundCodeLabel = &fundCodeLabel
	fundDataItem.fundNameLabel = &fundNameLabel
	fundDataItem.fundUdText = fundUdText

	return fundDataItem
}
