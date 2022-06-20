package practice

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

func init() {
	fmt.Println(".... init next ...")
}
func TPractice05() {
	name := flag.String("name", "value", "--name 90 or --name=90")
	//
	flag.Parse()
	fmt.Println(*name)
	fmt.Println(os.Args)
}

func TPractice06() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(arr)
	index := 2
	arr1 := arr[0:index]
	arr2 := arr[index+2:]

	arr2 = append(arr1, arr2...)
	fmt.Println(arr)
	fmt.Println(arr2)
}

func TPractice07() {
	// 图片大小
	const size = 300
	// 根据给定大小创建灰度图
	pic := image.NewGray(image.Rect(0, 0, size, size))
	// 遍历每个像素
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			// 填充为白色
			pic.SetGray(x, y, color.Gray{255})
		}
	}
	// 从0到最大像素生成x坐标
	for x := 0; x < size; x++ {
		// 让sin的值的范围在0~2Pi之间
		s := float64(x) * 2 * math.Pi / size
		// sin的幅度为一半的像素。向下偏移一半像素并翻转
		y := size/2 - math.Sin(s)*size/2
		// 用黑色绘制sin轨迹
		pic.SetGray(x, int(y), color.Gray{0})
	}
	// 创建文件
	file, err := os.Create("sin.png")
	if err != nil {
		log.Fatal(err)
	}
	// 使用png格式将数据写入文件
	png.Encode(file, pic) //将image信息写入文件中
	// 关闭文件
	file.Close()
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TPractice08() {
	arr := []string{"I", "am", "stupid", "and", "weak"}
	b, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	per := new(person)
	per.Name = "p1"
	per.Age = 30
	c, err2 := json.Marshal(per)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(string(c))

	objStr := string(c)
	var obj interface{}
	err3 := json.Unmarshal([]byte(objStr), &obj)
	if err3 == nil {
		objMap, ok := obj.(map[string]interface{})
		if ok {
			for k, v := range objMap {
				switch value := v.(type) {
				case string:
					fmt.Printf("type of %s is string, value is %v\n", k, value)
				case interface{}:
					fmt.Printf("type of %s is interface{}, value is %v\n", k, value)
				default:
					fmt.Printf("type of %s is default, value is %v\n", k, value)
				}
			}
		}
	}

}

func TPractice09() {
	baseCtx := context.Background()
	// 1秒超时
	timeOutCtx, cancel := context.WithTimeout(baseCtx, time.Second)
	defer cancel()
	go func(ctx context.Context) {
		// 时钟每秒跑一次
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			// C 通道
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupted.")
				return
			default:
				fmt.Println("enter default")
			}
		}
	}(timeOutCtx)

	time.Sleep(time.Second * 1)

	select {
	case <-timeOutCtx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("main process exited.")
	}

}

// 生产者消费者
func TPractice10() {
	ch := make(chan int, 10)
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		// 加_ 和 不加_【延时需要加长】 运行有区别
		for range ticker.C {
			select {
			case <-done:
				fmt.Println("child process exited")
				return
			default:
				fmt.Printf("receive data: %d\n", <-ch)
			}
		}
	}()
	// 生产者
	for i := 0; i < 10; i++ {
		ch <- i
	}
	time.Sleep(5 * time.Second)
	close(done)
	time.Sleep(1 * time.Second)
	fmt.Println("exited.")

}
