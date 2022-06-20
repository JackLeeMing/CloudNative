package practice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func init() {
	fmt.Println(".... init index ...")
}

func TPractice01() {
	var a []int
	b := []int{1, 2, 3}
	c := a
	a = append(b, 4)
	// c != a
	c = a
	fmt.Println(c, a, len(c), len(a), cap(c), cap(a))
	a = append(a, 5, 6)
	fmt.Println(c, a, len(c), len(a), cap(c), cap(a))
	a = append(a, 7)
	// a = append(a, 8)
	// append 之后切片容量发生变化的话就代表着新生成的切片将不会和原切片共享一个底层数组
	a[1] = 100
	a[0] = 100
	fmt.Println(c, a, len(c), len(a), cap(c), cap(a))
	// len 是切片实际的数据量 cap 是切片最大数据量 超过这个量切片将重新分配存储空间(重新分配空间)

	for _, value := range a {
		fmt.Println(value)
		var v = &value
		*v *= 10
		fmt.Println(value)
		// value 的地址并不是 切片中对应位置的 value 的地址
	}
	fmt.Println(a)
}

func TPractice02() {
	m := make(map[string]func(a, b int) int)
	m["add"] = func(a, b int) int {
		return a + b
	}
	m["mul"] = func(a, b int) int {
		return a * b
	}
	fmt.Println(m["add"](3, 4))
}

type ServiceType string

type MyType struct {
	// 多个 tag 用 空格分割
	Name string `json:"name" gson:"gName"`
	Age  int    `json:"age" gorm:"type:json;NOT NULL;default:'{\"query\":[], \"path\":[],\"body\":[],\"header\":[]}'"`
}

func TPractice03() {
	mt := MyType{Name: "MyType"}
	myType := reflect.TypeOf(mt)
	println(myType.NumField())
	name := myType.Field(1)
	tag := name.Tag.Get("gorm")
	println(tag)
}

func TPractice04() {
	arr := []string{"I", "am", "stupid", "and", "weak"}
	for i, v := range arr {
		if v == "stupid" {
			arr[i] = "smart"
		}
		if v == "weak" {
			arr[i] = "strong"
		}
	}
	b, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	println(strings.Join(arr, ","))
}
