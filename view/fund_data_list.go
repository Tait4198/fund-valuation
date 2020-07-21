package view

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"me/asdqwer/fund-valuation/data"
)

type FundDataList struct {
	widget.ScrollContainer
	box     *widget.Box
	fundMap map[string]*FundDataItem
}

func (v *FundDataList) PrependItem(item *FundDataItem) {
	if _, ok := v.fundMap[item.fundCode]; !ok {
		v.fundMap[item.fundCode] = item
		item.index = 0
		v.box.Prepend(item)
	}

}

func (v *FundDataList) AppendItem(item *FundDataItem) {
	if _, ok := v.fundMap[item.fundCode]; !ok {
		v.fundMap[item.fundCode] = item
		item.index = len(v.box.Children)
		v.box.Append(item)
	}
}

func (v *FundDataList) RemoveItem(index int) {
	removeItem := v.box.Children[index].(*FundDataItem)
	delete(v.fundMap, removeItem.fundCode)

	children := v.box.Children
	v.box.Children = append(children[:index], children[index+1:]...)
	for i, item := range v.box.Children {
		item.(*FundDataItem).index = i
	}
	v.box.Refresh()
}

func (v *FundDataList) GetItem(index int) *FundDataItem {
	return v.box.Children[index].(*FundDataItem)
}

func (v *FundDataList) UpdateItemByFundCode(fundCode string, fundData data.FundData) {
	for _, item := range v.box.Children {
		pItem := item.(*FundDataItem)
		if fundCode == pItem.fundCode {
			pItem.UpdateUdValue(fundData)
		}
	}
}

func (v *FundDataList) CheckItemByFundCode(fundCode string) bool {
	_, ok := v.fundMap[fundCode]
	return ok
}

func NewFundDataList(list ...*FundDataItem) *FundDataList {
	dataList := &FundDataList{}
	dataList.ExtendBaseWidget(dataList)
	dataList.Direction = widget.ScrollVerticalOnly
	dataList.SetMinSize(fyne.Size{Width: 600, Height: 400})

	dataList.fundMap = make(map[string]*FundDataItem)

	box := widget.Box{Horizontal: false}
	for i, item := range list {
		item.index = i
		dataList.fundMap[item.fundCode] = item
		box.Append(item)
	}
	dataList.Content = &box
	dataList.box = &box
	return dataList
}
