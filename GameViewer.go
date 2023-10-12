package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type GameViewer struct {
	mf *MineField
	//用来结束changeTimeString 的Goroutine
	Done          chan bool
	timer         time.Ticker
	timeString    binding.String
	mineNum       int
	mineNumString binding.ExternalInt
	mapData       [][]int
	GUI           fyne.App
	MainWindow    fyne.Window
	ManuWindow    fyne.Window
	btnArray      []*widget.Button
	isFlag        bool
}

// 创建app
func (gv *GameViewer) AppInit() {
	gv.GUI = app.New()
}

// 加载游戏数据
func (gv *GameViewer) GameInit(length int, width int, mineNum int) {
	//建设雷场
	gv.mf = &MineField{}
	gv.mf.init(length, width, mineNum)
	//启动地图视图
	//mapViewer := &MapViewer{}
	//mapViewer.init(gv.mf.length, gv.mf.width, *gv.mf)
	//mapViewer.Show()
	//GameViewer其他设置
	gv.mineNum = gv.mf.mineNum
	gv.mineNumString = binding.BindInt(&gv.mineNum)
	gv.timeString = binding.NewString()
	gv.timeString.Set("0:0")
	gv.mapData = MapData2ArrayList(gv.mf.length, gv.mf.width, gv.mf.mapData)
	gv.isFlag = false
	gv.timer = *time.NewTicker(1 * time.Second)
	gv.Done = make(chan bool)
}

// 显示菜单主页面
func (manu *GameViewer) ShowManu() {
	manu.ManuWindow = manu.GUI.NewWindow("manu")
	manu.ManuWindow.Resize(fyne.NewSize(300, 300))
	label := widget.NewLabelWithStyle("MineSweepr", fyne.TextAlignCenter, widget.RichTextStyleCodeBlock.TextStyle)
	label.Resize(fyne.NewSize(300, 200))
	levelEzBtn := widget.NewButton("easy level", func() {
		manu.GameInit(9, 9, 10)
		manu.Start()
		manu.ManuWindow.Hide()
	})
	levelMiddleBtn := widget.NewButton("middle level", func() {
		manu.GameInit(16, 16, 40)
		manu.Start()
		manu.ManuWindow.Hide()
	})
	levelHardBtn := widget.NewButton("hard level", func() {
		manu.GameInit(30, 16, 99)
		manu.Start()
		manu.ManuWindow.Hide()
	})
	container := container.NewVBox(label, levelEzBtn, levelMiddleBtn, levelHardBtn)
	manu.ManuWindow.SetContent(container)
	manu.ManuWindow.CenterOnScreen()
	//run()只能在主Goroutine,按钮新建窗口Show()就行
	manu.ManuWindow.ShowAndRun()
}

// 游戏退出，关闭主界面，显示菜单
func (gv GameViewer) QuitGame(isQuit bool) {
	gv.MainWindow.Close()
	gv.ManuWindow.Show()

}

// 接收定时器消息，更新timeString
func (gv *GameViewer) changeTimeString() {
	time := 0
	second := 0
	minute := 0
	isDone := false
	for !isDone {
		select {
		case <-gv.Done:
			isDone = true
		case <-gv.timer.C:
			time++
			if time == 60 {
				minute++
				second = 0
				time = 0
			}
			second = time
			gv.timeString.Set(fmt.Sprintf("%d:%d", minute, second))
		}
	}
}

// 游戏结束 true 表示赢了 false 表示踩到地雷
func (gv *GameViewer) GameOver(isWin bool) {
	message := " "
	if isWin {
		message = "You Win!"
	} else {
		message = "You've touched a mine!"
	}
	//展示所有地雷
	for y := 0; y < gv.mf.width; y++ {
		for x := 0; x < gv.mf.length; x++ {
			if gv.mapData[y][x] == -1 {
				gv.btnArray[gv.mf.length*y+x].Hide()
			}
		}
	}
	//杀掉时间
	gv.Done <- true
	gv.timer.Stop()
	dia := dialog.NewConfirm("GameOver", message, gv.QuitGame, gv.MainWindow)
	dia.Show()
}

// 开始游戏
// 游戏主显示逻辑
func (gv *GameViewer) Start() {
	gv.MainWindow = gv.GUI.NewWindow("MineSweeper")
	gv.MainWindow.Resize(fyne.NewSize(400, 400))
	vbox := container.NewVBox()
	gv.btnArray = make([]*widget.Button, 0)
	//时间&flag&地雷数显示
	TimeLabel := widget.NewLabelWithData(gv.timeString)
	MineLabelString := binding.IntToString(gv.mineNumString)
	MineLabel := widget.NewLabelWithData(MineLabelString)
	var FlagBtn *widget.Button
	FlagBtn = widget.NewButton("🚩", func() {
		if gv.isFlag {
			FlagBtn.Text = "🚩"
			gv.isFlag = false
			FlagBtn.Refresh()
		} else {
			gv.isFlag = true
			FlagBtn.Text = "🏴‍☠️"
			FlagBtn.Refresh()
		}
	})
	dataContainer := container.NewGridWithColumns(3, MineLabel, FlagBtn, TimeLabel)

	//地图显示
	mapContainer := container.NewGridWithColumns(gv.mf.length)
	for y := 0; y < len(gv.mapData); y++ {
		for x := 0; x < len(gv.mapData[y]); x++ {
			stackContainer := container.NewStack()
			//label,在btn之下
			var label *canvas.Text
			pointNum := gv.mapData[y][x]
			if pointNum == -1 {
				label = canvas.NewText("💣", nil)
			} else if pointNum == 0 {
				label = canvas.NewText(" ", nil)
			} else {
				label = canvas.NewText(fmt.Sprintf("%d", pointNum), getColor(pointNum))
			}
			label.Alignment = fyne.TextAlignCenter
			stackContainer.Add(label)
			//btn点击逻辑
			var btn *widget.Button
			btn = widget.NewButton(" ", func() {
				//如果是旗子
				if gv.isFlag {
					if btn.Text != "🚩" {
						btn.SetText("🚩")
						i, _ := gv.mineNumString.Get()
						gv.mineNumString.Set(i - 1)
					} else {
						btn.SetText(" ")
						i, _ := gv.mineNumString.Get()
						gv.mineNumString.Set(i + 1)
					}
					//否则
				} else {
					//加判断防止插旗之后可以点开
					if btn.Text != "🚩" {
						btn.Hidden = true
						position := gv.getBtnPosition(btn)
						btnNum := gv.mapData[position[0]][position[1]]
						if btnNum == 0 {
							gv.openField(position)
						}
						if btnNum == -1 {
							gv.GameOver(false)
						}
					}
				}
				//判定是否胜利
				if gv.isWin() {
					gv.GameOver(true)
				}
			})
			stackContainer.Add(btn)
			gv.btnArray = append(gv.btnArray, btn)
			mapContainer.Add(stackContainer)
		}
	}
	vbox.Add(dataContainer)
	vbox.Add(mapContainer)
	go gv.changeTimeString()
	gv.MainWindow.SetContent(vbox)
	gv.MainWindow.CenterOnScreen()
	gv.MainWindow.Show()
	//关闭时展示菜单界面
	gv.MainWindow.SetOnClosed(func() {
		gv.ManuWindow.Show()
	})

}

// 是否胜利
func (gv GameViewer) isWin() bool {
	btnNum := 0
	for _, btn := range gv.btnArray {
		if !btn.Hidden {
			btnNum++
		}
	}
	return btnNum == gv.mf.mineNum
}

// 在点击到0时打开一段范围（到周围非数字的地区）的雷区，帮助玩家更好的游玩
// 参数：自身所在的xy坐标
// 算法描述：查找自身四周范围是否为数字单元，如果隐藏过了就跳过，如是则隐藏按钮，如果是0则将0位置传给openField继续调用
func (gv *GameViewer) openField(positon []int) {
	posX := positon[1]
	posY := positon[0]
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			pointX := posX + x
			pointY := posY + y
			//地址合法
			if pointX >= 0 && pointX < gv.mf.length && pointY >= 0 && pointY < gv.mf.width {
				pointNum := gv.mapData[pointY][pointX]
				btn := gv.btnArray[pointY*gv.mf.length+pointX]
				//如果没有这个判断就会形成死循环
				if btn.Hidden || btn.Text == "🚩" {
					continue
				}
				btn.Hide()
				if pointNum == 0 {
					newPos := make([]int, 0)
					newPos = append(newPos, pointY)
					newPos = append(newPos, pointX)
					gv.openField(newPos)
				}

			}
		}
	}
}

// 获取btn位置数组
func (gv *GameViewer) getBtnPosition(btn *widget.Button) []int {
	for i := 0; i < len(gv.btnArray); i++ {
		if btn == gv.btnArray[i] {
			y := i / gv.mf.length
			x := i % gv.mf.length
			result := make([]int, 0)
			result = append(result, y)
			result = append(result, x)
			return result
		}
	}
	return nil
}
