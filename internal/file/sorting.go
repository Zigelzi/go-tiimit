package file

import "sort"

type ByNewestDate []FileName

func (fn ByNewestDate) Len() int           { return len(fn) }
func (fn ByNewestDate) Less(i, j int) bool { return fn[i].Date.After(fn[j].Date) }
func (fn ByNewestDate) Swap(i, j int)      { fn[i], fn[j] = fn[j], fn[i] }

func sortByNewestDate(fileNames []FileName) {
	sort.Sort(ByNewestDate(fileNames))
}
