package main

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
