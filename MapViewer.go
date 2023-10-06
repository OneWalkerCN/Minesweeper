package main

import (
	"fmt"
)

/*MapViewer 用来观看地图数据的一个视图，DebugTool*/
type MapViewer struct {
	field  [][]int
	length int
	width  int
}

// 加载视图
func (v *MapViewer) init(length int, width int, mf MineField) {
	v.length = length
	v.width = width
	v.field = MapData2ArrayList(length, width, mf.mapData)
}

// 显示视图
func (v MapViewer) Show() {
	for i := 0; i < v.length; i++ {
		fmt.Printf("%3s", "——")
	}
	fmt.Println()
	for _, verticalList := range v.field {
		fmt.Print("|")
		for _, point := range verticalList {
			fmt.Printf("%3d", point)
		}
		fmt.Print("|\n")
	}
	for i := 0; i < v.length; i++ {
		fmt.Printf("%3s", "——")
	}
}
