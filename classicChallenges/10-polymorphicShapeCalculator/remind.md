## 用 make 创建 slice

``` go
    var slice []type = make([]type, len)
    slice  := make([]type, len)
    slice  := make([]type, len, cap)
```
## sort 用法

``` go
package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
// ByAge 实现了基于 Age 字段对 []Person 进行 sort.Interface 排序
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func main() {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	fmt.Println(people)
    // There are two ways to sort a slice. First, one can define
	// a set of methods for the slice type, as with ByAge, and
	// call sort.Sort. In this first example we use that technique.
    // 排序一个 slice 有两种方法。
    // 首先，可以为 slice 类型定义一组方法，就像 ByAge 一样，
    // 然后调用 sort.Sort。在这个第一个示例中，我们使用了这种技术。
	sort.Sort(ByAge(people))
	fmt.Println(people)

    // The other way is to use sort.Slice with a custom Less
	// function, which can be provided as a closure. In this
	// case no methods are needed. (And if they exist, they
	// are ignored.) Here we re-sort in reverse order: compare
	// the closure with ByAge.Less.
    // 另一种方法是使用 sort.Slice 和自定义的 Less 函数，
    // 该函数可以作为闭包提供。
    // 在这种情况下，不需要方法。(如果它们存在，则会被忽略。)
    // 这里我们按相反的顺序重新排序：比较闭包和 ByAge.Less。
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	fmt.Println(people)

}
Output:

[Bob: 31 John: 42 Michael: 17 Jenny: 26]
[Michael: 17 Jenny: 26 Bob: 31 John: 42]
[John: 42 Bob: 31 Jenny: 26 Michael: 17]
```

例如本题中的用法：

在排序前先创建一个副本切片，不修改原切片。

``` go
// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sortedShapes := make([]Shape, len(shapes))
	copy(sortedShapes, shapes)

	sort.Slice(sortedShapes, func(i, j int) bool {
		areaI := sortedShapes[i].Area()
		areaJ := sortedShapes[j].Area()
		if ascending {
			return areaI < areaJ
		} else {
			return areaI > areaJ
		}
	})

	return sortedShapes
}
```
