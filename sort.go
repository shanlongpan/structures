package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var arrays []int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {

		arrays = append(arrays, int(r.Intn(1000000)))
	}
	//var arrays = []int{12, 423, 53, 765, 97, 432, 24}
	timer:=time.Now()
	bubbleSort(arrays)
	fmt.Println(time.Since(timer))

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	arrays= []int{}
	for i := 0; i < 10000; i++ {

		arrays = append(arrays, int(r.Intn(1000000)))
	}
	timer=time.Now()
	selectSort(arrays)
	fmt.Println(time.Since(timer))

	arrays= []int{}
	for i := 0; i < 10000; i++ {

		arrays = append(arrays, int(r.Intn(1000000)))
	}
	timer=time.Now()
	insertSort(arrays)
	fmt.Println(time.Since(timer))

	arrays= []int{}
	for i := 0; i < 10000; i++ {

		arrays = append(arrays, int(r.Intn(1000000)))
	}
	timer=time.Now()
	qsort(arrays)
	fmt.Println(time.Since(timer))
}
//冒泡排序 （排序10000个整数，用时约117ms）
func bubbleSort(nums []int) []int {
	arrayLen := len(nums)
	for i := 0; i < arrayLen-1; i++ {
		for j := 1; j < arrayLen-i; j++ {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
	return nums
}

//选择排序 （排序10000个整数，用时约101ms）
func selectSort(nums []int) []int {
	arrayLen := len(nums)
	for i := 0; i < arrayLen-1; i++ {
		minIndex := i
		for j := i + 1; j < arrayLen; j++ {
			if nums[minIndex] > nums[j] {
				minIndex = j
			}
		}
		if minIndex != i {
			nums[minIndex], nums[i] = nums[i], nums[minIndex]
		}
	}
	return nums
}

//插入排序（排序10000个整数，用时约14ms）
func insertSort(array []int) {
	n := len(array)
	if n < 2 {
		return
	}
	for i := 1; i < n; i++ {
		for j := i - 1; j >= 0; j-- {
			if array[j] > array[j+1] {
				array[j], array[j+1] = array[j+1],array[j]
			}else{
				break
			}
		}
	}
}

//快速排序
func qsort(data []int) []int{
	if len(data) <= 1 {
		return data
	}
	mid := data[0]
	head, tail := 0, len(data)-1
	for i := 1; i <= tail; {
		if data[i] > mid {
			data[i], data[tail] = data[tail], data[i]
			tail--
		} else {
			data[i], data[head] = data[head], data[i]
			head++
			i++
		}
	}
	qsort(data[:head])
	qsort(data[head+1:])

	return data
}
