package util

import(
    "sort"
)

type sorter struct {
    len int
    fnSwap func(i, j int)
    fnLess func(i, j int) bool
}
func (s sorter) Len() int { return s.len }
func (s sorter) Swap(i, j int) {
    s.fnSwap(i, j)
}
func (s sorter) Less(i, j int) bool {
    return s.fnLess(i, j)
}

//slice排序
//利用 closure 实现 Go 的泛型。fnSwap 和 fnLess 一般为闭包函数
// https://medium.com/capital-one-tech/closures-are-the-generics-for-go-cb32021fb5b5
func Sort(len int, fnSwap func(i, j int), fnLess func(i, j int) bool) {
    sort.Sort(sorter{len, fnSwap, fnLess})
}

//slice稳定排序，如果相同大小情况下，之前在前面的，排序后保证还在前面
func Stable(len int, fnSwap func(i, j int), fnLess func(i, j int) bool) {
    sort.Stable(sorter{len, fnSwap, fnLess})
}
