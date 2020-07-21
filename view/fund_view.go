package view

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"me/asdqwer/fund-valuation/config"
	"me/asdqwer/fund-valuation/data"
)

type FundContainer struct {
	Container  *fyne.Container
	dl         *FundDataList
	cfg        *config.AppConfig
	fWin       *fyne.Window
	progress   *dialog.ProgressInfiniteDialog
	upProgress *dialog.ProgressDialog
}

type StockBox struct {
	widget.Box
	stockCode   string
	priceText   *canvas.Text
	udText      *canvas.Text
	updateLabel *widget.Label
}

func (*StockBox) CalcUd(stockData data.StockData) string {
	diff := stockData.TodayPrice - stockData.YesterdayPrice
	val := fmt.Sprintf("%.2f", (diff/stockData.YesterdayPrice)*100)
	return val
}

func (v *StockBox) UpdateUd() {
	stockData, err := data.GetStockData(v.stockCode)
	if err == nil {
		ud := v.CalcUd(stockData)
		udColor := GetUdColor(ud)
		v.udText.Text = ud + "%"
		v.udText.Color = udColor
		v.priceText.Text = fmt.Sprintf("%.2f", stockData.TodayPrice)
		v.priceText.Color = udColor
		v.updateLabel.Text = stockData.UpdateTime
		v.Refresh()
	}
}

func newStockBox(stockCode string) *StockBox {
	stockBox := &StockBox{}
	stockBox.ExtendBaseWidget(stockBox)
	stockBox.Horizontal = true
	stockBox.stockCode = stockCode
	stockData, err := data.GetStockData(stockCode)
	if err == nil {
		ud := stockBox.CalcUd(stockData)

		nameLabel := widget.Label{Text: stockData.Name}
		updateLabel := widget.Label{Text: stockData.UpdateTime}
		udText := canvas.NewText(ud+"%", GetUdColor(ud))
		priceText := canvas.NewText(fmt.Sprintf("%.2f", stockData.TodayPrice), GetUdColor(ud))

		stockBox.Append(&nameLabel)
		stockBox.Append(priceText)
		stockBox.Append(udText)
		stockBox.Append(layout.NewSpacer())
		stockBox.Append(&updateLabel)

		stockBox.updateLabel = &updateLabel
		stockBox.priceText = priceText
		stockBox.udText = udText
	}
	return stockBox
}

func (v *FundContainer) addNewFund(code string, cfg *config.AppConfig) {
	if !v.dl.CheckItemByFundCode(code) {
		fundData, err := data.GetFundData(code)
		if err == nil {
			item := NewFundItem(fundData)
			item.OnRemove = itemRemove(v.dl, cfg)
			v.dl.AppendItem(item)

			newFundCfg := config.FundConfig{Code: code}
			cfg.Funds = append(cfg.Funds, newFundCfg)
			config.UpdateAppConfig()
		}
	}
}

func (v *FundContainer) updateFund(cfg *config.AppConfig, up func(p float64)) {
	length := float64(len(cfg.Funds))
	for i, fundConfig := range cfg.Funds {
		fundData, err := data.GetFundData(fundConfig.Code)
		if err == nil {
			if v.dl.CheckItemByFundCode(fundConfig.Code) {
				v.dl.UpdateItemByFundCode(fundData.FundCode, fundData)
			} else {
				item := NewFundItem(fundData)
				item.OnRemove = itemRemove(v.dl, cfg)
				v.dl.AppendItem(item)
			}
		}
		if up != nil {
			up(float64(i) / length)
		}
	}
}

func (v *FundContainer) InitFund() {
	length := float64(len(v.cfg.Funds))
	v.upProgress.SetValue(0)
	v.upProgress.Show()
	for i, fundConfig := range v.cfg.Funds {
		fundData, err := data.GetFundData(fundConfig.Code)
		if err == nil {
			item := NewFundItem(fundData)
			item.OnRemove = itemRemove(v.dl, v.cfg)
			v.dl.AppendItem(item)
		}
		v.upProgress.SetValue(float64(i) / length)
	}
	v.upProgress.Hide()
}

func itemRemove(dl *FundDataList, cfg *config.AppConfig) RemoveFunc {
	return func(index int) {
		fundCode := dl.GetItem(index).fundCode
		cfgIndex := -1
		for i, item := range cfg.Funds {
			if item.Code == fundCode {
				cfgIndex = i
				break
			}
		}
		if cfgIndex != -1 {
			cfg.Funds = append(cfg.Funds[:cfgIndex], cfg.Funds[cfgIndex+1:]...)
			dl.RemoveItem(index)
			config.UpdateAppConfig()
		}
	}
}

func NewFundContainer(cfg *config.AppConfig, win *fyne.Window) *FundContainer {
	fundContainer := &FundContainer{}

	fundCodeInput := widget.NewEntry()
	fundCodeInput.SetPlaceHolder("输入基金代码")

	progress := dialog.NewProgressInfinite("正在添加", "正在添加请稍等...", *win)
	progress.Hide()

	upProgress := dialog.NewProgress("正在更新", "正在更新请稍等...", *win)
	upProgress.Hide()

	shStockBox := newStockBox("0000001")

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			upProgress.Show()
			upProgress.SetValue(0)
			shStockBox.UpdateUd()
			fundContainer.updateFund(cfg, func(p float64) {
				upProgress.SetValue(p)
			})
			upProgress.Hide()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			fundCodeInput.Text = ""
			fundCodeInput.Refresh()
			dialog.ShowCustomConfirm("添加基金", "确定", "取消", fundCodeInput, func(b bool) {
				if b {
					progress.Show()
					fundContainer.addNewFund(fundCodeInput.Text, cfg)
					progress.Hide()
				}
			}, *win)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	dl := NewFundDataList()

	fundContainer.dl = dl
	fundContainer.progress = progress
	fundContainer.upProgress = upProgress
	fundContainer.fWin = win
	fundContainer.cfg = cfg

	fundContainer.Container = fyne.NewContainerWithLayout(layout.NewVBoxLayout(), toolbar, shStockBox, dl)
	return fundContainer
}
