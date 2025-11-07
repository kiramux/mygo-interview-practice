# **Go Generics 学习资料**

## **Go 中的泛型简介**

Go 1.18 引入了对泛型编程的支持，使开发者能够编写可以处理多种类型的代码，同时保持类型安全。这一特性在不牺牲编译时类型检查的情况下，提供了更灵活和可重用的代码。

### **为什么使用泛型？**

在引入泛型之前，Go 开发者有几种处理多类型的方式：

- **Interface{}**：使用空接口允许函数接受任何类型，但需要进行类型断言，且失去了编译时的类型检查。
- **代码生成**：像 `go generate` 这样的工具可以生成特定类型的实现，但增加了构建过程的复杂性。
- **复制粘贴**：为不同类型复制代码会导致维护问题。

泛型通过提供一种既类型安全又可在多种类型之间重用的方式，解决了这些问题。

### **基础语法**

定义一个泛型函数的基本语法：

``` go
func MyGenericFunction[T any](param T) T {
    // 函数体
    return param
}
```

以及一个泛型类型的定义：

``` go
type MyGenericType[T any] struct {
    Value T
}
```

### **类型参数和约束**

类型参数在方括号 `[T any]` 中指定，其中：

- T 是类型参数的名称
- `any` 是类型约束（在这种情况下，允许任何类型）

Go 提供了几个预定义的约束，可以在 `constraints` 包中找到：

``` go
import "golang.org/x/exp/constraints"

// 一个可以处理任何有序类型的函数
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

### **自定义类型约束**

你可以通过接口类型定义自定义约束：

``` go
// 定义一个要求有 String() 方法的约束
type Stringer interface {
    String() string
}

// 一个可以处理任何实现 String() 方法的类型的函数
func PrintValue[T Stringer](value T) {
    fmt.Println(value.String())
}
```

### **约束中的联合类型**

Go 泛型支持约束中的联合类型，允许一个参数接受多个特定类型：

``` go
// 一个接受 int 或 float64 的约束
type Number interface {
    int | float64
}

// 一个可以处理 int 或 float64 的函数
func Add[T Number](a, b T) T {
    return a + b
}
```

### **类型集合**

类型集合是 Go 泛型实现中的核心概念。约束定义了一组符合该约束的类型：

``` go
// 一个约束，要求类型可以与 == 和 != 操作符进行比较
type Comparable[T any] interface {
    comparable
}

// 一个检查两个值是否相等的函数
func AreEqual[T comparable](a, b T) bool {
    return a == b
}
```

### **泛型数据结构**

泛型特别适合用于实现数据结构：

``` go
// 泛型栈实现
type Stack[T any] struct {
    elements []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{elements: make([]T, 0)}
}

func (s *Stack[T]) Push(element T) {
    s.elements = append(s.elements, element)
}

func (s *Stack[T]) Pop() (T, error) {
    var zero T
    if len(s.elements) == 0 {
        return zero, errors.New("stack is empty")
    }
    
    lastIndex := len(s.elements) - 1
    element := s.elements[lastIndex]
    s.elements = s.elements[:lastIndex]
    return element, nil
}

func (s *Stack[T]) Peek() (T, error) {
    var zero T
    if len(s.elements) == 0 {
        return zero, errors.New("stack is empty")
    }
    
    return s.elements[len(s.elements)-1], nil
}

func (s *Stack[T]) Size() int {
    return len(s.elements)
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.elements) == 0
}
```

### **类型推断**

Go 可以根据函数参数自动推断类型参数：

``` go
func Identity[T any](value T) T {
    return value
}

// 类型推断示例
str := Identity("hello")    // T 被推断为 string
num := Identity(42)         // T 被推断为 int
```

### **多个类型参数**

函数和类型可以具有多个类型参数：

``` go
// 一个将一个类型的切片转换为另一个类型的映射函数
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// 使用示例
numbers := []int{1, 2, 3, 4}
squares := Map(numbers, func(x int) int { return x * x })
// squares: [1, 4, 9, 16]

// 将数字转换为字符串
strNumbers := Map(numbers, func(x int) string { return strconv.Itoa(x) })
// strNumbers: ["1", "2", "3", "4"]
```

### **将泛型与方法结合**

可以在泛型类型上定义方法：

``` go
type Pair[T, U any] struct {
    First  T
    Second U
}

func (p Pair[T, U]) Swap() Pair[U, T] {
    return Pair[U, T]{First: p.Second, Second: p.First}
}

// 使用示例
pair := Pair[string, int]{First: "answer", Second: 42}
swapped := pair.Swap() // Pair[int, string]{First: 42, Second: "answer"}
```

### **约束包**

`golang.org/x/exp/constraints` 包提供了有用的约束：

``` go
import "golang.org/x/exp/constraints"

// 一个可以处理任何整数类型的函数
func Sum[T constraints.Integer](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

// 一个可以处理任何浮点类型的函数
func Average[T constraints.Float](values []T) T {
    sum := T(0)
    for _, v := range values {
        sum += v
    }
    return sum / T(len(values))
}
```

常用的约束包括：

- Integer：任何整数类型
- Float：任何浮点数类型
- Complex：任何复数类型
- Ordered：支持 `<` 运算符的类型
- Signed：任何带符号整数类型
- Unsigned：任何无符号整数类型

### **泛型算法**

泛型非常适合用于实现多类型的算法：

``` go
// 泛型二分查找函数
func BinarySearch[T constraints.Ordered](slice []T, target T) int {
    left, right := 0, len(slice)-1
    
    for left <= right {
        mid := (left + right) / 2
        
        if slice[mid] == target {
            return mid
        } else if slice[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1 // 未找到
}
```

### **方法中的类型参数**

方法本身不能拥有与接收者类型分离的类型参数，但你可以使用泛型函数来解决这个问题：

``` go
// 这段代码无法编译 - 方法不能有自己的类型参数
// func (s *Stack[T]) ConvertTo[U any](converter func(T) U) []U { ... }

// 采用普通函数解决
func ConvertStack[T, U any](stack *Stack[T], converter func(T) U) []U {
    result := make([]U, stack.Size())
    for i, v := range stack.elements {
        result[i] = converter(v)
    }
    return result
}
```

### **泛型类型的零值**

在使用泛型时，经常需要生成类型参数的“零值”：

``` go
func GetZero[T any]() T {
    var zero T
    return zero
}

// 使用示例
zeroInt := GetZero[int]()       // 0
zeroString := GetZero[string]() // ""
```

### **性能考虑**

Go 中的泛型在实现时对性能进行了精心优化：

- **编译方式**：Go 使用混合方式，为每种类型实例化生成特定代码，同时尽可能共享代码。
- **运行时效率**：泛型代码在编译时优化，因此与手写的类型特定代码相比，几乎没有运行时开销。
- **代码大小**：使用大量类型实例化可能增加二进制文件的大小，但编译器会尽量减少这一影响。

### **使用泛型的最佳实践**

- **不要过度使用泛型**：只有在泛型能够明确提高代码重用性和类型安全时，才应使用泛型。
- **对约束要具体**：根据实际需求，尽量使用最具体的约束类型。
- **提供清晰的文档**：清晰地文档化泛型函数和类型的预期行为，方便他人理解和使用。
- **考虑性能影响**：在使用泛型时，要留意泛型对编译时间和二进制大小的影响。

### **实际应用案例**

## **泛型结果类型**

常见的一种模式是创建一个泛型结果类型，用于处理成功和错误的情况：

``` go
type Result[T any] struct {
    Value T
    Error error
}

func NewSuccess[T any](value T) Result[T] {
    return Result[T]{Value: value, Error: nil}
}

func NewError[T any](err error) Result[T] {
    var zero T
    return Result[T]{Value: zero, Error: err}
}

// 使用示例
func DivideInts(a, b int) Result[int] {
    if b == 0 {
        return NewError[int](errors.New("division by zero"))
    }
    return NewSuccess(a / b)
}
```

## **泛型集合实现**

``` go
type Set[T comparable] struct {
    elements map[T]struct{}
}

func NewSet[T comparable]() Set[T] {
    return Set[T]{elements: make(map[T]struct{})}
}

func (s *Set[T]) Add(element T) {
    s.elements[element] = struct{}{}
}

func (s *Set[T]) Remove(element T) {
    delete(s.elements, element)
}

func (s *Set[T]) Contains(element T) bool {
    _, exists := s.elements[element]
    return exists
}

func (s *Set[T]) Size() int {
    return len(s.elements)
}

func (s *Set[T]) Elements() []T {
    result := make([]T, 0, len(s.elements))
    for element := range s.elements {
        result = append(result, element)
    }
    return result
}

// 集合操作
func Union[T comparable](s1, s2 Set[T]) Set[T] {
    result := NewSet[T]()
    
    for element := range s1.elements {
        result.Add(element)
    }
    
    for element := range s2.elements {
        result.Add(element)
    }
    
    return result
}

func Intersection[T comparable](s1, s2 Set[T]) Set[T] {
    result := NewSet[T]()
    
    for element := range s1.elements {
        if s2.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}
```

## [补充]题目需要注意的地方

### Reduce 函数/归约函数

```go
// Reduce reduces a slice to a single value by applying a function to each element
// 归约函数，initial U：初始累积值，接收"当前累积值"和"下一个元素"，返回"新的累积值"
// 不止可以求和，还可以 求乘积、连接字符串、找最大值 等
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}
```

**求乘积**：

```go
product := Reduce(numbers, 1, func(acc, n int) int {
    return acc * n
})
```

**连接字符串**：

```go
sentence := Reduce(words, "", func(acc, word string) string {
    return acc + " " + word
})
```

**找最大值**：

```go
max := Reduce(numbers, numbers[0], func(acc, n int) int {
    if n > acc { return n }
    return acc
})
```

### map 判断某个键是否存在

``` go
value, ok := map[key]
```

### Union() 新建 set 不可直接修改传入参数

`func Union[T comparable](s1, s2 *Set[T]) *Set[T]` 中，要求返回一个新 set ，不可直接修改 s1, s2，（假定修改 s2 作为“新set”）否则会带来不可预料的副作用：

1. 数据意外污染
    ``` go
    func Union[T comparable](s1, s2 *Set[T]) *Set[T]
    ```
2. 并发安全问题
    多个 goroutine 同时调用 Union，可能导致数据竞争和不一致状态
3. 链式调用问题
    ``` go
    set1 := NewSet[int]()
    set1.Add(1)

    set2 := NewSet[int]()
    set2.Add(2)

    set3 := NewSet[int]()
    set3.Add(3)

    // 程序员可能期望这样链式调用
    result := Union(Union(set1, set2), set3)

    // 但是这会导致 set2 被修改两次！
    // 第一次 Union 后，set2 = {1, 2}
    // 第二次 Union 后，set2 = {1, 2, 3}
    ```

## **进一步阅读**

1. [Go Generics Design Document](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
2. [Go by Example: Generics](https://gobyexample.com/generics)
3. [Go Generics 101](https://go101.org/generics/101.html)
4. [The Go Programming Language Blog: Using Generics in Go](https://go.dev/blog/intro-generics)
5. [Go Generics in Practice](https://bitfieldconsulting.com/golang/generics)

