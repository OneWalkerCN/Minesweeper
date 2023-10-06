package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GameViewer struct {
	mf      *MineField
	time    string
	mineNum int
	mapData [][]int
}

func (gv *GameViewer) Init(mf *MineField) {
	gv.mf = mf
	gv.time = "0:0"
	gv.mineNum = mf.mineNum
	gv.mapData = MapData2ArrayList(mf.length, mf.width, mf.mapData)
}

func (gv GameViewer) Show() {
	GUI := app.New()
	MainWindow := GUI.NewWindow("MineSwipper")
	layOut := container.NewGridWithColumns(gv.mf.width)
	for _, row := range gv.mapData {
		for _, point := range row {
			//max只显示堆叠在最上层的图像
			maxContainer := container.NewStack()
			//通过声明获取btn自身
			var btn *widget.Button
			btn = widget.NewButton(" ", func() {
				btn.Hidden = true
			})
			txt := " "
			if point != 0 {
				txt = fmt.Sprintf("%d", point)
			}
			label := canvas.NewText(txt, getColor(point))
			maxContainer.Add(label)
			maxContainer.Add(btn)

			layOut.Add(maxContainer)
		}
	}
	//GUI.Settings().SetTheme(theme.LightTheme())
	MainWindow.SetContent(layOut)
	MainWindow.ShowAndRun()

}

// 根据数字改变颜色
func getColor(num int) color.Color {
	switch num {
	case 1:
		return color.White
	case 2:
		return color.RGBA{0, 0, 255, 255}
	case 3:
		return color.RGBA{0, 255, 0, 255}
	case 4:
		return color.RGBA{255, 255, 0, 255}
	case 5:
		return color.RGBA{255, 0, 0, 255}
	case 6:
		return color.RGBA{254, 67, 101, 255}
	default:
		return color.Black
	}
}
