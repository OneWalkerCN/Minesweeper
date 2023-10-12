package main

import "image/color"

//将Vector3数据列表转换成二维数组
func MapData2ArrayList(length int, width int, mapData []Vector3) [][]int {
	result := make([][]int, width)
	for i := range result {

		result[i] = make([]int, length)
	}
	for _, point := range mapData {
		result[point.Y][point.X] = point.Num
	}
	return result
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
