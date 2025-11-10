# 银行账户错误处理学习材料
## Go 中的错误处理

错误处理是编写健壮 Go 程序的关键方面。这个挑战专注于实现一个具有正确错误处理技术的银行系统。

## 基础错误处理

Go 使用显式错误处理，通过返回值而不是异常：

```go
// 可能返回错误的函数
// Function that may return an error
func divideNumbers(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 使用 if 检查进行错误处理
// Error handling with if checks
result, err := divideNumbers(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return // 或处理错误 | or handle the error
}
fmt.Println("Result:", result)
```

## 创建自定义错误

创建错误的标准方式：

```go
// 使用 errors.New 创建简单错误消息
// Using errors.New for simple error messages
err := errors.New("insufficient funds")

// 使用 fmt.Errorf 创建格式化错误消息
// Using fmt.Errorf for formatted error messages
amount := 100
err := fmt.Errorf("insufficient funds: need $%d more", amount)
```

## 自定义错误类型

创建自定义错误类型允许更详细的错误处理：

```go
// 定义自定义错误类型
// Define a custom error type
type InsufficientFundsError struct {
    Balance float64
    Amount  float64
}

// 实现 error 接口
// Implement the error interface
func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf("insufficient funds: balance $%.2f, attempted to withdraw $%.2f", 
        e.Balance, e.Amount)
}
```

自定义错误的关键概念：
- 实现 `Error() string` 方法
- 在错误字段中包含相关上下文
- 检查错误类型时使用指针接收器
- 使用带 ok 模式的类型断言进行错误类型检查

## 错误包装（Go 1.13+）

Go 1.13 引入了错误包装以提供更好的错误上下文：

关键概念：
- 使用带 `%w` 动词的 `fmt.Errorf` 来包装错误
- 使用 `errors.As()` 检查链中的特定错误类型
- 使用 `errors.Is()` 检查链中的特定错误值
- 在添加有意义信息的同时保留原始错误上下文

## 哨兵错误

可以直接比较的预定义错误：

```go
// 将哨兵错误定义为包级变量
// Define sentinel errors as package-level variables
var (
    ErrAccountNotFound   = errors.New("account not found")
    ErrInsufficientFunds = errors.New("insufficient funds")
    ErrInvalidAmount     = errors.New("invalid amount")
)
```

关键概念：
- 使用包级变量创建可重用的错误
- 使用 `==` 或 `errors.Is()` 比较错误
- 提供清晰、描述性的错误消息

## 错误处理模式

### 1. 提前返回模式
首先验证输入并立即返回错误以避免深层嵌套。

### 2. 错误处理函数
创建能够依次处理多个容易出错操作的函数。

### 3. 错误上下文
在返回或包装错误时始终提供有意义的上下文。

## 银行应用特定考虑

### 账户操作
- **余额验证**：在取款前检查资金是否充足
- **金额验证**：确保存款/取款金额为正数
- **账户存在性**：在操作前验证账户是否存在
- **输入清理**：验证所有用户输入

### 银行错误类型
- **InsufficientFundsError**：余额问题的特定错误
- **InvalidAmountError**：负数或零金额
- **AccountNotFoundError**：账户查找失败时
- **ValidationError**：输入验证失败

## 银行应用中的线程安全

银行应用需要处理并发访问：

关键概念：
- 使用 `sync.Mutex` 保护账户操作
- 在检查余额和修改余额前加锁
- 使用 `defer` 确保互斥锁始终被解锁
- 对于读取密集型操作考虑读写锁

## 测试错误场景

测试错误处理至关重要：

测试策略：
- 分别测试每个错误条件
- 验证错误类型和消息
- 在处理错误后测试成功操作
- 对多个错误场景使用表驱动测试
- 模拟依赖项以模拟错误条件

## 错误日志记录和报告

正确的错误日志记录是必不可少的：

日志记录最佳实践：
- 记录错误时提供充分的上下文
- 包含相关 ID（账户、交易、用户）
- 在适当级别记录日志（错误、警告、信息）
- 不要在调用栈中多次记录相同错误
- 构建日志以便于解析和监控

## Panic 和 Recover

虽然 Go 偏向显式错误处理，但 panic 和 recover 可用于异常情况：

### 何时使用 panic：
- 表示程序员错误的不可恢复错误
- 阻止程序工作的初始化失败
- 内部一致性违反

### 恢复模式：
- 使用带 `recover()` 的 `defer` 捕获 panic
- 在适当时将 panic 转换为错误
- 记录 panic 以便调试
- 仅在适当边界处恢复

## 进一步阅读

- [Go 中的错误处理](https://blog.golang.org/error-handling-and-go)
- [Go 1.13+ 中的错误处理](https://blog.golang.org/go1.13-errors)
- [Go 中的有效错误处理](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- [Sync 包文档](https://pkg.go.dev/sync)

## 总结

Go 中的错误处理遵循明确和可预测的原则。通过使用自定义错误类型、哨兵错误和错误包装，你可以创建提供清晰错误信息和上下文的健壮银行系统。记住始终验证输入、处理并发访问，并彻底测试错误场景以确保你的银行应用在所有条件下都能可靠运行。

正确的错误处理不仅仅是防止程序崩溃——它是关于创建用户友好的体验和为调试和维护提供有用信息。在银行等关键系统中，强健的错误处理可能意味着成功的交易和灾难性的失败之间的差异。