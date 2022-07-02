package main

import (
	"fmt"

	"github.com/JackLeeMing/CloudNative/lock"
)

func init() {
	fmt.Println(".... init main ...")
}

func main() {
	lock.TryQueue()
}
