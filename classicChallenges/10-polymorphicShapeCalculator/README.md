# 挑战 10：多态图形计算器

## **问题说明**

实现一个系统，用于通过 Go 的接口计算各种几何图形的属性。此挑战重点在于理解并正确实现 Go 的接口系统，以实现**多态性（polymorphism）**。

---

## **要求**

实现一个 `Shape` 接口，包含以下方法：

* `Area() float64`：计算图形的面积
* `Perimeter() float64`：计算图形的周长（或圆的周长）
* `String() string`：返回图形的字符串表示（实现 `fmt.Stringer` 接口）

实现以下具体图形类型：

* **Rectangle（矩形）**：由宽度和高度定义
* **Circle（圆形）**：由半径定义
* **Triangle（三角形）**：由三条边定义（面积使用海伦公式计算）

实现一个 **ShapeCalculator（图形计算器）**，用于：

* 接收任意图形并返回其属性
* 计算多个图形的总面积
* 找出面积最大的图形
* 按面积升序或降序排序图形

---

## **函数签名**

```go
// Shape interface
type Shape interface {
    Area() float64
    Perimeter() float64
    fmt.Stringer // Includes String() string method
}

// Concrete types
type Rectangle struct {
    Width, Height float64
}

type Circle struct {
    Radius float64
}

type Triangle struct {
    SideA, SideB, SideC float64
}

// Constructor functions
func NewRectangle(width, height float64) (*Rectangle, error)
func NewCircle(radius float64) (*Circle, error)
func NewTriangle(a, b, c float64) (*Triangle, error)

// ShapeCalculator
type ShapeCalculator struct{}

func NewShapeCalculator() *ShapeCalculator
func (sc *ShapeCalculator) PrintProperties(s Shape)
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape
```

---

## **约束条件**

* 所有尺寸必须为正数。
* 三角形三边必须满足**三角不等式定理**（任意两边之和大于第三边）。
* 在构造函数中实现参数验证并返回合适的错误信息。
* 计算圆形属性时应使用常量 π（pi）。
* `String()` 方法应返回包含图形类型及尺寸的格式化字符串。

---

## **示例用法**

```go
// Create shapes
rect, _ := NewRectangle(5, 3)
circle, _ := NewCircle(4)
triangle, _ := NewTriangle(3, 4, 5)

// Use shapes polymorphically
calculator := NewShapeCalculator()
shapes := []Shape{rect, circle, triangle}

// Calculate total area
totalArea := calculator.TotalArea(shapes)
fmt.Printf("Total area: %.2f\n", totalArea)

// Sort shapes by area
sortedShapes := calculator.SortByArea(shapes, true)
for _, s := range sortedShapes {
    calculator.PrintProperties(s)
}

// Find largest shape
largest := calculator.LargestShape(shapes)
fmt.Printf("Largest shape: %s with area %.2f\n", largest, largest.Area())
```

---

这道题旨在巩固 Go 中的**接口实现与多态性**理解。
关键在于让不同图形类型实现同一接口，使得计算器可以在无需了解具体类型的情况下统一操作所有图形。
