package main

/*GameManager 管理游戏运行的管理器*/
type GameMananger struct {
	GameTime   string
	GameScore  int
	GameMap    *MineField
	MapViewer  MapViewer
	GameViewer GameViewer
}

/*游戏加载：加载地图，加载视图*/
func (gm *GameMananger) GameInit(length int, width int, mineNum int) {

	//设置场地
	gm.GameMap = &MineField{}
	gm.GameMap.init(length, width, mineNum)
	//启动地图视图
	gm.MapViewer.init(length, width, *gm.GameMap)
	gm.GameViewer.Init(gm.GameMap)
}

/*开始游戏*/
func (gm GameMananger) Start() {
	gm.MapViewer.Show()
	gm.GameViewer.Show()
}

/*结束游戏*/
func (gm GameMananger) End() {}
