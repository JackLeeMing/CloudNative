package array

import (
	"fmt"
	"math/rand"
)

func compute(a, b int) int {
	return a + b + rand.Intn(2)
}

func ArraSort() {
	arr := []int{1, 4, 3, 2, 1, 6}
	fmt.Println(arr)
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				temp := arr[i]
				arr[i] = arr[j]
				arr[j] = temp
			}
		}
	}
	fmt.Println(arr)
	for i, n := range arr {
		fmt.Println(i, n)
	}
}
