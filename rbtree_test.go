package rbtree

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"testing"
)

type insertTest struct {
	Nodes int
	Color map[int][]bool
	Nums  map[int][]int
}

var testInsertCases = []insertTest{
	{
		1, map[int][]bool{0: {black}}, map[int][]int{0: {1}},
	},
	{
		2, map[int][]bool{0: {black}, 1: {black, red}}, map[int][]int{0: {1}, 1: {0, 2}},
	},
	{
		3, map[int][]bool{0: {black}, 1: {red, red}}, map[int][]int{0: {2}, 1: {1, 3}},
	},
	{
		4, map[int][]bool{0: {black}, 1: {black, black}, 2: {black, black, black, red}}, map[int][]int{0: {2}, 1: {1, 3}, 2: {0, 0, 0, 4}},
	},
}

type deleteTest struct {
	Nodes   int
	Deleted []int
	Color   map[int][]bool
	Nums    map[int][]int
}

var testDeleteCases = []deleteTest{
	{
		1, []int{1}, map[int][]bool{}, map[int][]int{},
	},
	{
		2, []int{1}, map[int][]bool{0: {black}}, map[int][]int{0: {2}},
	},
	{
		3, []int{2}, map[int][]bool{0: {black}, 1: {red, black}}, map[int][]int{0: {3}, 1: {1, 0}},
	},
}

type newInt struct {
	value int
}

func (i newInt) Less(than Item) bool {
	j := than.(newInt)
	return i.value < j.value
}
func insertAndDeleteRB(info deleteTest) *Tree {
	tree := Init()
	for i := 1; i < info.Nodes+1; i++ {
		var tmp newInt
		tmp.value = i
		tree.InsertNoReplace(tmp)
	}
	for _, v := range info.Deleted {
		var tmp newInt
		tmp.value = v
		tree.Delete(tmp)
	}
	return tree

}
func TestDelete(t *testing.T) {
	for _, v := range testDeleteCases {
		tree := insertAndDeleteRB(v)
		//	draw(tree, v.Nodes)
		color := getColor(tree, v.Nodes-len(v.Deleted))
		nums := getNums(tree, v.Nodes-len(v.Deleted))
		if tree.count != v.Nodes-len(v.Deleted) {
			t.Error("nodes count not equal ")
			return
		}
		for k, colors := range v.Color {
			getColors, ok := color[k]
			if !ok {
				t.Error("not exist color:", k, colors, v.Nodes)
				return
			}
			if len(getColors) != len(colors) {
				t.Error("color length not equal")
				return
			}
			for index := range colors {
				if colors[index] != getColors[index] {
					t.Error("color not equal")
					return
				}
			}
		}
		for k, num := range v.Nums {
			calcNums, ok := nums[k]
			if !ok {
				t.Error("not exist nums")
				return
			}
			if len(calcNums) != len(num) {
				t.Error("num length not equal")
				return
			}
			for index := range num {
				if num[index] != calcNums[index] {
					t.Error("num not equal")
					return
				}
			}
		}

	}
}
func TestInsertNoReplace(t *testing.T) {
	for _, v := range testInsertCases {
		tree := insertRB(v)
		//	draw(tree, v.Nodes)
		color := getColor(tree, v.Nodes)
		nums := getNums(tree, v.Nodes)
		if tree.count != v.Nodes {
			t.Error("nodes count not equal")
			return
		}
		for k, colors := range v.Color {
			getColors, ok := color[k]
			if !ok {
				t.Error("not exist color")
				return
			}
			if len(getColors) != len(colors) {
				t.Error("color length not equal")
				return
			}
			for index := range colors {
				if colors[index] != getColors[index] {
					t.Error("color not equal")
					return
				}
			}
		}
		for k, num := range v.Nums {
			calcNums, ok := nums[k]
			if !ok {
				t.Error("not exist nums")
				return
			}
			if len(calcNums) != len(num) {
				t.Error("num length not equal")
				return
			}
			for index := range num {
				if num[index] != calcNums[index] {
					t.Error("num not equal")
					return
				}
			}
		}

	}
}

//TODO make it  right now just for test
func draw(tree *Tree, nodes int) {
	root := tree.root
	var slice []*Node

	slice = append(slice, root)
	for {
		if nodes == 0 {
			break
		}
		length := len(slice)
		var lines string
		for i := 0; i < length; i++ {
			if slice[i] != nil {
				if tree.isBlack(slice[i]) {
					color.New(color.FgBlue).Fprintf(os.Stdout, " %d ", slice[i].Item.(newInt).value)

				} else {
					color.New(color.FgRed).Fprintf(os.Stdout, "  %d ", slice[i].Item.(newInt).value)
					//color.Red("%d  ", slice[i].Item.(newInt).value)
				}
				slice = append(slice, leftOf(slice[i]), rightOf(slice[i]))
				nodes--
				lines += "/  \\  "
			} else {
				color.New(color.FgBlack).Fprintf(os.Stdout, "nil")
				lines += "    "
			}
		}
		fmt.Println()
		fmt.Println(lines)

		slice = slice[length:]

	}
}
func getNums(tree *Tree, nodes int) map[int][]int {
	root := tree.root
	if root == nil {
		return nil
	}
	nums := make(map[int][]int)
	slice := make([]*Node, 0)
	slice = append(slice, root)
	length := 1
	count := 0
	for {
		if nodes == 0 {
			break
		}
		length = len(slice)
		for i := 0; i < length; i++ {
			if slice[i] != nil {
				nums[count] = append(nums[count], slice[i].Item.(newInt).value)
			} else {
				nums[count] = append(nums[count], 0)
			}

			slice = append(slice, leftOf(slice[i]), rightOf(slice[i]))
			if slice[i] != nil {
				nodes--
			}
		}
		slice = slice[length:]
		count++
	}
	return nums
}
func getColor(tree *Tree, nodes int) map[int][]bool {
	root := tree.root
	if root == nil {
		return nil
	}
	color := make(map[int][]bool)
	slice := make([]*Node, 0)
	slice = append(slice, root)
	length := 1
	count := 0
	fmt.Printf("root :%+v %p\n", root, root)
	for {
		if nodes == 0 {
			break
		}
		length = len(slice)
		for i := 0; i < length; i++ {

			color[count] = append(color[count], tree.isBlack(slice[i]))

			slice = append(slice, leftOf(slice[i]), rightOf(slice[i]))
			if slice[i] != nil {
				nodes--
			}
		}
		slice = slice[length:]
		count++
	}
	return color
}
func insertRB(info insertTest) *Tree {
	tree := Init()
	for i := 1; i < info.Nodes+1; i++ {
		var tmp newInt
		tmp.value = i
		tree.InsertNoReplace(tmp)
	}
	return tree

}
