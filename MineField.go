package main

import (
	"math/rand"
	"time"
)

/*地图*/

// 三元向量，存储地图信息
type Vector3 struct {
	X   int
	Y   int
	Num int
}

func (v3 Vector3) isEqual(point Vector3) bool {
	//fmt.Println("1")
	if v3.X == point.X && v3.Y == point.Y {
		return true
	} else {
		return false
	}
}

// 地图数据
type MineField struct {
	//地图长宽
	length int
	width  int
	//地图元素列表(-1代表有，其他数字代表附近有几个雷)
	mapData []Vector3
	mineNum int
}

// 添加地雷，并除重，添加成功返回true,否则返回false
func (mf *MineField) addMines(v3 Vector3) bool {
	//遍历mapData,如果相同位置返回false
	if len(mf.mapData) == 0 {
		v3.Num = -1
		mf.mapData = append(mf.mapData, v3)
	}
	//如果有相同的返回false
	for _, point := range mf.mapData {

		if point.isEqual(v3) {
			return false
		}
	}
	//没有相同的，则将Num值-1,加入地图数据
	v3.Num = -1
	mf.mapData = append(mf.mapData, v3)
	return true
}

// 如果该点在mapData中，返回index；不在则返回-1
func (mf MineField) getPointIndex(input Vector3) int {
	for i := 0; i < len(mf.mapData); i++ {
		if mf.mapData[i].isEqual(input) {
			return i
		}
	}
	return -1
}

// 查询新建单元是否出界 出界返回true，没出返回false
func (mf MineField) isOutOfBound(point Vector3) bool {
	length := mf.length - 1
	width := mf.width - 1
	if point.X < 0 || point.X > length {
		return true
	}
	if point.Y < 0 || point.Y > width {
		return true
	}
	return false
}

// 新建数字单元,检测是否越界，没有越界则加入成功，成功返回true,不成功返回false
func (mf *MineField) addNums(v3 Vector3) bool {
	if !mf.isOutOfBound(v3) {
		v3.Num++
		mf.mapData = append(mf.mapData, v3)
		return true
	}
	return false
}

// 计算附近地雷数
/*
算法简述：
1. 遍历列表中现有的地雷数据，将地雷单元格周边(8个方位)单元格Num+1（如果不存在则创建新单元格对象）；
2. 如果一周存在地雷单元，则跳过该单元格。
*/
func (mf *MineField) getAroundMinesNum() {
	//问题：随着Nums的数据加入，可能会遍历到非地雷成员，因此需要判断自身是否为雷
	//否则有些Nums数据会多加1
	//问题：发现多次重复访问单元，导致数据不对
	//解决方法：append只在slice末尾加入新元素，因此先获取地雷数据的slice长度，只遍历地雷数据即可
	mineList := mf.mapData
	for i := 0; i < len(mineList); i++ {
		//遍历周边单元
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				aroundPoint := &Vector3{
					X: mineList[i].X + x,
					Y: mineList[i].Y + y,
				}
				//如果是本身，则跳过
				if aroundPoint.isEqual(mineList[i]) {
					continue
				}
				//如果Data中存在
				//1. 判断是否是雷，如果是则跳过
				//2. 如果不是雷，则Num加一
				index := mf.getPointIndex(*aroundPoint)
				if index != -1 {
					if mf.mapData[index].Num == -1 {
						continue
					} else {
						mf.mapData[index].Num++
					}
				} else {
					//如果不存在，则添加；如果非法则放弃添加
					mf.addNums(*aroundPoint)
				}

			}
		}
	}

}

// 加载地图 (长，宽，地雷数)
func (mf *MineField) init(length int, width int, mineNum int) {
	mf.length = length
	mf.width = width
	mf.mineNum = mineNum
	rand.NewSource(time.Now().UnixNano())
	//填埋地雷
	for mineNum > 1 {
		v3 := &Vector3{} // 改为创建指向新的Vector3实例的指针
		v3.X = rand.Intn(mf.length)
		v3.Y = rand.Intn(mf.width)
		if mf.addMines(*v3) {
			mineNum--
		}
	}
	//计算地雷数
	mf.getAroundMinesNum()
}
