package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"me/asdqwer/fund-valuation/config"
	"me/asdqwer/fund-valuation/view"
	"os"
)

func main() {

	cfg := config.GetAppConfig()

	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exPath := filepath.Dir(ex)
	//fontPath := fmt.Sprintf("%s/%s", exPath, cfg.Font)
	// 绝对路径
	fontPath := "/Users/xxx/fund-valuation/WeiRuanYaHei.ttf"

	os.Setenv("FYNE_FONT", fontPath)
	os.Setenv("FYNE_THEME", "light")
	defer os.Unsetenv("FYNE_FONT")
	defer os.Unsetenv("FYNE_THEME")

	fApp := app.New()
	fWin := fApp.NewWindow("Title")
	size := fyne.NewSize(600, 400)
	fWin.Resize(size)

	fundContainer := view.NewFundContainer(cfg, &fWin)
	fWin.SetContent(fundContainer.Container)
	fundContainer.InitFund()
	fWin.ShowAndRun()

}
