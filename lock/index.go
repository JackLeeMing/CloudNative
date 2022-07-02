package lock

import (
	"fmt"
	"sync"
	"time"
)

func MutexTest(ch <-chan int, ch2 chan<- int) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Lock()
	p := <-ch
	ch2 <- 2
	fmt.Println(p)
	// 定义 channel 的时候 箭头指向 chan 时代表的是发送通道;箭头背向 chan 时代表的是接收通道
}

type Pod struct {
	Name string
}
type Group struct {
}

func (group *Group) GroupTest(pods []*Pod) []*Pod {
	var wg sync.WaitGroup
	ps := make([]*Pod, len(pods))
	for i, pod := range pods {
		wg.Add(1)
		go func(i int, pod *Pod) {
			defer wg.Done()
			fmt.Printf("end Pod %d\n", i)
			ps[i] = &Pod{Name: pod.Name + "- tag"}
		}(i, pod)
	}
	wg.Wait()
	return ps
}

type Queue struct {
	queue []string
	cond  *sync.Cond
}

func (q *Queue) Enqueue(item string) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.queue = append(q.queue, item)
	fmt.Println("入队列")
	q.cond.Broadcast()
}

func (q *Queue) Dequeue() string {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if len(q.queue) == 0 {
		fmt.Println("no data and waiting.")
		q.cond.Wait()
	}
	result := q.queue[0]
	q.queue = q.queue[1:]
	return result
}

func TryQueue() {
	q := Queue{
		queue: []string{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}
	go func() {
		for {
			q.Enqueue("a")
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		q.Dequeue()
		time.Sleep(time.Second)
	}
}
