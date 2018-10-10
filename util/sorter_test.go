package util

import (
    "testing"
    "fmt"
)

type exampleObj struct {
    name string
    f1   int
    f2   int
}

func TestSort(t *testing.T) {
    data := []int{2, 1, 5, 3, 4}
    Sort(len(data), func(i, j int) {
        data[i], data[j] = data[j], data[i]
    }, func(i, j int) bool {
        return data[i] < data[j]
    })
    fmt.Printf("sorted data: %+v\n", data)

    objs := []exampleObj {
        exampleObj{"obj1", 2, 5},
        exampleObj{"obj2", 1, 10},
        exampleObj{"obj3", 9, 2},
        exampleObj{"obj4", 1, 10},
        exampleObj{"obj5", 7, 8},
        exampleObj{"obj5", 9, 9},
    }
    Stable(len(objs), func(i, j int) {
        objs[i], objs[j] = objs[j], objs[i]
    }, func(i, j int) bool {
        if objs[i].f1 > objs[j].f1 {
            return true
        }
        if objs[i].f2 < objs[j].f2 {
            return true
        }
        return false
    })

    fmt.Printf("sorted objs: %+v\n", objs)
}
