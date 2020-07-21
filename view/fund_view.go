package view

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"me/asdqwer/fund-valuation/config"
	"me/asdqwer/fund-valuation/data"
)

type FundView struct {
	widget.Box
	dl *FundDataList
}

func (v *FundView) AddNewFund(code string, cfg *config.AppConfig) {
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

func (v *FundView) UpdateFund(cfg *config.AppConfig) {
	for _, fundConfig := range cfg.Funds {
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
	}
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

func NewFundView(cfg *config.AppConfig, win *fyne.Window) *FundView {
	fundView := &FundView{}
	fundView.ExtendBaseWidget(fundView)
	fundView.Horizontal = false

	fundCodeInput := widget.NewEntry()
	fundCodeInput.SetPlaceHolder("输入基金代码")

	progress := dialog.NewProgressInfinite("正在处理", "正在处理请稍等...", *win)
	progress.Hide()

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			progress.Show()
			fundView.UpdateFund(cfg)
			progress.Hide()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			fundCodeInput.Text = ""
			fundCodeInput.Refresh()
			dialog.ShowCustomConfirm("添加基金", "确定", "取消", fundCodeInput, func(b bool) {
				if b {
					progress.Show()
					fundView.AddNewFund(fundCodeInput.Text, cfg)
					progress.Hide()
				}
			}, *win)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	fundView.Append(toolbar)

	dl := NewFundDataList()
	fundView.dl = dl
	fundView.Append(dl)

	progress.Show()
	for _, fundConfig := range cfg.Funds {
		fundData, err := data.GetFundData(fundConfig.Code)
		if err == nil {
			item := NewFundItem(fundData)
			item.OnRemove = itemRemove(dl, cfg)
			dl.AppendItem(item)
		}
	}
	progress.Hide()

	return fundView
}
