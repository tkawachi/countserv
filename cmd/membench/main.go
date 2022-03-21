package main

import (
	"fmt"
	"os"

	"github.com/axiomhq/hyperloglog"
	"github.com/google/uuid"
)

func main() {
	arr := []*hyperloglog.Sketch{}
	for i := 0; i < 100; i++ {
		sketch := hyperloglog.New()
		arr = append(arr, sketch)
		// Generate random UUID
		for j := 0; j < 10000; j++ {
			id := uuid.New().String()
			// id := fmt.Sprintf("user-%d", j)
			sketch.Insert([]byte(id))
		}
	}
	fmt.Println(len(arr))

	b, err := os.ReadFile("/proc/self/status")
	if err != nil {
		panic(err)
	}
	// $ go run cmd/membench/main.go | grep -e VmPeak -e VmHWM
	// VmPeak:  1223084 kB
	// VmHWM:     13944 kB
	fmt.Println(string(b))
}
