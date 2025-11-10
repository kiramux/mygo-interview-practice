# **挑战 30：上下文管理实现**

## **概述**

实现一个上下文管理器，演示 Go 中上下文包的基本模式。上下文包在 Go 应用中用于管理取消信号、超时和请求范围内的值。

## **任务**

实现一个包含 6 个核心方法和 2 个辅助函数的 `ContextManager` 接口：

### `ContextManager` 接口

```go
type ContextManager interface {
    // 从父上下文创建一个可取消的上下文
    CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)
    
    // 创建一个带超时的上下文
    CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
    
    // 向上下文添加一个值
    AddValue(parent context.Context, key, value interface{}) context.Context
    
    // 从上下文获取一个值
    GetValue(ctx context.Context, key interface{}) (interface{}, bool)
    
    // 执行一个带有上下文取消支持的任务
    ExecuteWithContext(ctx context.Context, task func() error) error
    
    // 等待一段时间或直到上下文被取消
    WaitForCompletion(ctx context.Context, duration time.Duration) error
}
```

### 辅助函数

```go
// 模拟可以取消的工作
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error

// 处理多个项目时考虑上下文的影响
func ProcessItems(ctx context.Context, items []string) ([]string, error)
```

## **要求**

### 核心功能

- **上下文取消**：通过 `context.WithCancel` 处理手动取消。
- **上下文超时**：通过 `context.WithTimeout` 实现超时行为。
- **值存储**：通过 `context.WithValue` 存储和检索值。
- **任务执行**：在支持取消的上下文下执行函数。
- **等待操作**：在等待时考虑取消信号。

### 实现细节

- 使用 Go 的标准上下文包函数。
- 处理 `context.Canceled` 和 `context.DeadlineExceeded` 错误。
- 返回适当的布尔标志，指示值是否存在。
- 支持基于 goroutine 的任务执行，确保适当的同步。
- 批量处理项目时，检查每个项目之间的取消信号。

### 测试覆盖范围

你的实现将通过 13 个测试用例进行测试，涵盖以下内容：

- 上下文创建与取消
- 超时行为
- 值存储与检索
- 任务执行场景（成功、错误、取消）
- 等待操作（完成与取消）
- 辅助函数行为
- 集成场景

### 开始实现

- 检查解决方案模板和测试文件。
- 从简单的方法开始，例如 `AddValue` 和 `GetValue`。
- 逐步实现取消和超时上下文。
- 实现带有适当 goroutine 处理的任务执行。
- 经常运行测试，使用命令 `go test -v`。

**提示**：可以参考 `learning.md` 文件，了解上下文模式和示例的详细内容！