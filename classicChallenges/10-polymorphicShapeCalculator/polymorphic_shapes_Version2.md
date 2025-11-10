# 多态形状计算器学习材料
## Go 中的接口和多态

这个挑战专注于使用 Go 的接口来实现几何形状计算的多态性。

## 理解 Go 中的接口

在 Go 中，接口定义行为而不指定实现。接口是方法签名的集合：

```go
// 定义一个接口
// Define an interface
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

类型通过实现接口的方法隐式地实现接口：

```go
// Rectangle 实现 Shape 接口
// Rectangle implements the Shape interface
type Rectangle struct {
    Width  float64
    Height float64
}

// 实现 Area 方法
// Implement the Area method
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// 实现 Perimeter 方法
// Implement the Perimeter method
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle 也实现 Shape 接口
// Circle also implements the Shape interface
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}
```

## 使用接口实现多态

接口允许多态行为——不同类型可以基于它们的行为被统一对待：

```go
// 适用于任何 Shape 的函数
// Function that works with any Shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

// 使用
// Usage
rect := Rectangle{Width: 5, Height: 3}
circ := Circle{Radius: 2}

PrintShapeInfo(rect)  // 适用于 Rectangle | Works with Rectangle
PrintShapeInfo(circ)  // 适用于 Circle | Works with Circle
```

## 接口值

接口值由两个组件组成：

- **动态类型**：存储在接口中的具体类型
- **动态值**：该类型的实际值

```go
var s Shape                // nil 接口值（nil 类型，nil 值）
                          // nil interface value (nil type, nil value)
s = Rectangle{5, 3}        // s 的类型是 Rectangle，值是 Rectangle{5, 3}
                          // s has type Rectangle, value Rectangle{5, 3}
s = Circle{2.5}            // s 现在的类型是 Circle，值是 Circle{2.5}
                          // s now has type Circle, value Circle{2.5}
```

## 空接口

空接口 `interface{}` 或 `any`（Go 1.18+）没有方法，可以持有任何值：

```go
func PrintAny(a interface{}) {
    fmt.Println(a)
}

PrintAny(42)              // 适用于 int | Works with int
PrintAny("Hello")         // 适用于 string | Works with string
PrintAny(Rectangle{5, 3}) // 适用于 Rectangle | Works with Rectangle
```

## 类型断言

类型断言从接口中提取底层值：

```go
// 单返回值的类型断言
// Type assertion with single return value
rect := s.(Rectangle) // 如果 s 不是 Rectangle 会 panic
                     // Panics if s is not a Rectangle

// 带检查的类型断言
// Type assertion with check
rect, ok := s.(Rectangle)
if ok {
    fmt.Println("It's a rectangle with width:", rect.Width)
} else {
    fmt.Println("It's not a rectangle")
}
```

## 类型开关

类型开关处理多种类型：

```go
func Describe(s Shape) string {
    switch v := s.(type) {
    case Rectangle:
        return fmt.Sprintf("Rectangle with width %.2f and height %.2f", v.Width, v.Height)
    case Circle:
        return fmt.Sprintf("Circle with radius %.2f", v.Radius)
    case nil:
        return "nil shape"
    default:
        return fmt.Sprintf("Unknown shape of type %T", v)
    }
}
```

## 接口组合

接口可以由其他接口组成：

```go
type Sizer interface {
    Area() float64
}

type Perimeterer interface {
    Perimeter() float64
}

// 组合接口
// Composed interface
type Shape interface {
    Sizer
    Perimeterer
    String() string  // 额外方法 | Additional method
}
```

## 接口嵌入

Go 允许将一个接口嵌入到另一个接口中：

```go
type Stringer interface {
    String() string
}

type Shape interface {
    Area() float64
    Perimeter() float64
}

// CompleteShape 嵌入 Shape 和 Stringer
// CompleteShape embeds Shape and Stringer
type CompleteShape interface {
    Shape
    Stringer
}
```

## 使用指针接收器的接口实现

方法接收器类型对接口实现很重要：

```go
type Modifier interface {
    Scale(factor float64)
}

// 值接收器 - 不修改原值
// Value receiver - doesn't modify original
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// 指针接收器 - 修改原值
// Pointer receiver - modifies original
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

var m Modifier
r := Rectangle{5, 3}

// 这样可以 - r 是可寻址的
// This works - r is addressable
m = &r
m.Scale(2)

// 这样不行 - 接口期望指针接收器
// This doesn't work - interface expects pointer receiver
// m = r // 编译错误 | Compile error
```

## 接口最佳实践

- **保持接口小**：偏向少方法的接口（通常只有一个）
- **在使用点定义接口**：在使用它们的包中定义，而不是在实现它们的地方
- **接口表示行为，不是类型**：专注于某样东西做什么，而不是它是什么

```go
// 好的 - 定义行为
// Good - defines behavior
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 不太好 - 定义类型
// Less good - defines a type
type Car interface {
    Drive()
    Stop()
    Refuel()
}
```

## 里氏替换原则

里氏替换原则指出超类的对象应该可以被子类的对象替换而不影响程序的正确性：

```go
// 常见的违反是在子类型中添加要求
// A common violation is adding requirements in subtypes
type Parallelogram interface {
    SetWidth(w float64)
    SetHeight(h float64)
    Area() float64
}

type Rectangle struct {
    width, height float64
}

func (r *Rectangle) SetWidth(w float64) { r.width = w }
func (r *Rectangle) SetHeight(h float64) { r.height = h }
func (r Rectangle) Area() float64 { return r.width * r.height }

type Square struct {
    side float64
}

// 这个实现打破了期望！
// This implementation breaks expectations!
func (s *Square) SetWidth(w float64) {
    s.side = w
    // Square 在设置一个维度时改变两个维度
    // Square changes both dimensions when one is set
}

func (s *Square) SetHeight(h float64) {
    s.side = h
}

func (s Square) Area() float64 { return s.side * s.side }
```

## 实践示例：形状计算器

让我们实现一个完整的形状计算器：

```go
package shape

import (
    "fmt"
    "math"
)

// Shape 是基本接口
// Shape is the basic interface
type Shape interface {
    Area() float64
    Perimeter() float64
    String() string
}

// Circle 实现
// Circle implementation
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
    return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Rectangle 实现
// Rectangle implementation
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
    return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Triangle 实现
// Triangle implementation
type Triangle struct {
    SideA float64
    SideB float64
    SideC float64
}

func (t Triangle) Perimeter() float64 {
    return t.SideA + t.SideB + t.SideC
}

func (t Triangle) Area() float64 {
    // 海伦公式
    // Heron's formula
    s := t.Perimeter() / 2
    return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

func (t Triangle) String() string {
    return fmt.Sprintf("Triangle(sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator 处理多个形状
// ShapeCalculator handles multiple shapes
type ShapeCalculator struct {
    shapes []Shape
}

func NewCalculator() *ShapeCalculator {
    return &ShapeCalculator{shapes: make([]Shape, 0)}
}

func (c *ShapeCalculator) AddShape(s Shape) {
    c.shapes = append(c.shapes, s)
}

func (c *ShapeCalculator) TotalArea() float64 {
    total := 0.0
    for _, s := range c.shapes {
        total += s.Area()
    }
    return total
}

func (c *ShapeCalculator) TotalPerimeter() float64 {
    total := 0.0
    for _, s := range c.shapes {
        total += s.Perimeter()
    }
    return total
}

func (c *ShapeCalculator) ListShapes() []string {
    result := make([]string, len(c.shapes))
    for i, s := range c.shapes {
        result[i] = s.String()
    }
    return result
}
```

## 使用新形状进行扩展

接口的一个优势是可以添加新类型而不改变现有代码：

```go
// 添加新形状：正多边形
// Add a new shape: Regular Polygon
type RegularPolygon struct {
    Sides     int
    SideLength float64
}

func (p RegularPolygon) Perimeter() float64 {
    return float64(p.Sides) * p.SideLength
}

func (p RegularPolygon) Area() float64 {
    return (float64(p.Sides) * p.SideLength * p.SideLength) / (4 * math.Tan(math.Pi/float64(p.Sides)))
}

func (p RegularPolygon) String() string {
    return fmt.Sprintf("RegularPolygon(sides=%d, length=%.2f)", p.Sides, p.SideLength)
}

// 无需更改即可与现有计算器一起工作
// Works with the existing calculator without changes
calculator.AddShape(RegularPolygon{Sides: 6, SideLength: 5})
```

## 使用接口进行测试

接口通过允许模拟实现来促进测试：

```go
// 接口定义
// Interface definition
type AreaCalculator interface {
    Area() float64
}

// 使用接口的函数
// Function that uses the interface
func IsLargeShape(s AreaCalculator) bool {
    return s.Area() > 100
}

// 使用模拟进行测试
// Test with a mock
type MockShape struct{
    MockArea float64
}

func (m MockShape) Area() float64 {
    return m.MockArea
}

func TestIsLargeShape(t *testing.T) {
    small := MockShape{50}
    large := MockShape{150}
    
    if IsLargeShape(small) {
        t.Error("Expected small shape to not be large")
    }
    
    if !IsLargeShape(large) {
        t.Error("Expected large shape to be large")
    }
}
```

## 进一步阅读

- [Go 接口教程](https://tour.golang.org/methods/9)
- [Effective Go：接口](https://golang.org/doc/effective_go#interfaces)
- [Go 中的 SOLID 设计](https://dave.cheney.net/2016/08/20/solid-go-design)
- [反射定律](https://blog.golang.org/laws-of-reflection)（深入理解接口）

## 总结

Go 的接口系统提供了强大而灵活的多态机制，使得代码更加模块化、可测试和可扩展。通过定义小而专注的接口，并在使用点定义它们，你可以构建松散耦合且易于维护的系统。

接口是 Go 语言的核心特性之一，掌握它们对于编写惯用的 Go 代码至关重要。记住要保持接口简单，专注于行为而不是数据，并利用 Go 的隐式接口实现来创建灵活的设计。