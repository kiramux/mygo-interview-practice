# 学习指南：Go Context 包

## 什么是 Context？

Go 中的 context 包提供了一种在 API 边界和 goroutine 之间传递取消信号、超时和请求范围值的方式。它是 Go 中构建健壮、生产就绪应用程序最重要的包之一。

## 为什么 Context 很重要

```go
// 没有 context - 无法取消或超时
func fetchData() ([]byte, error) {
    resp, err := http.Get("https://api.example.com/data")
    // 如果这需要 5 分钟怎么办？没有办法取消！
    // What if this takes 5 minutes? No way to cancel!
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

// 使用 context - 可取消且有时间限制
// With context - cancellable and time-bounded
func fetchDataWithContext(ctx context.Context) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com/data", nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err // 可能是 context.DeadlineExceeded 或 context.Canceled
                       // Could be context.DeadlineExceeded or context.Canceled
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

## 核心 Context 类型

### 1. Background Context
所有 context 的根 - 永不取消，没有截止时间，不携带值。

```go
ctx := context.Background()
// 在 main()、测试或初始化中使用它作为顶级 context
// Use this as the top-level context in main(), tests, or initialization
```

### 2. TODO Context
当你不确定使用哪个 context 时的占位符。

```go
ctx := context.TODO()
// 在开发过程中当 context 不明确时使用
// Use this during development when context isn't clear yet
```

### 3. Cancellation Context
可以手动取消以通知 goroutine 停止工作。

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // 始终调用 cancel 以防止内存泄漏
               // Always call cancel to prevent memory leaks

go func() {
    select {
    case <-ctx.Done():
        fmt.Println("Work cancelled:", ctx.Err())
        return
    case <-time.After(5 * time.Second):
        fmt.Println("Work completed")
    }
}()

// 2 秒后取消
// Cancel after 2 seconds
time.Sleep(2 * time.Second)
cancel() // 这会触发 ctx.Done()
         // This triggers ctx.Done()
```

### 4. Deadline/Timeout Context
在特定时间后自动取消。

```go
// WithDeadline - 在特定时间取消
// WithDeadline - cancel at specific time
deadline := time.Now().Add(30 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// WithTimeout - 在持续时间后取消
// WithTimeout - cancel after duration
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### 5. Value Context
在函数调用之间携带请求范围的数据。

```go
ctx := context.WithValue(context.Background(), "userID", "12345")
ctx = context.WithValue(ctx, "requestID", "req-abc-123")

// 检索值
// Retrieve values
userID := ctx.Value("userID").(string)
requestID := ctx.Value("requestID").(string)
```

## Context 模式

### 模式 1：检查取消

```go
func doWork(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // 定期检查取消
        // Check for cancellation periodically
        select {
        case <-ctx.Done():
            return ctx.Err() // context.Canceled 或 context.DeadlineExceeded
                            // context.Canceled or context.DeadlineExceeded
        default:
            // 继续工作
            // Continue work
        }
        
        // 模拟工作
        // Simulate work
        time.Sleep(10 * time.Millisecond)
        fmt.Printf("Processed item %d\n", i)
    }
    return nil
}
```

### 模式 2：Context 与 Goroutines

```go
func processInParallel(ctx context.Context, items []string) error {
    errChan := make(chan error, len(items))
    
    for _, item := range items {
        go func(item string) {
            select {
            case <-ctx.Done():
                errChan <- ctx.Err()
                return
            case errChan <- processItem(item):
                return
            }
        }(item)
    }
    
    // 等待所有 goroutines
    // Wait for all goroutines
    for i := 0; i < len(items); i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }
    
    return nil
}
```

### 模式 3：Context 竞争

```go
func executeWithTimeout(ctx context.Context, task func() error) error {
    done := make(chan error, 1)
    
    go func() {
        done <- task()
    }()
    
    select {
    case err := <-done:
        return err // 任务首先完成
                  // Task completed first
    case <-ctx.Done():
        return ctx.Err() // Context 首先取消/超时
                        // Context cancelled/timeout first
    }
}
```

## 真实世界示例

### 带请求超时的 Web 服务器

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 为此请求创建带超时的 context
    // Create context with timeout for this request
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    
    // 添加请求特定值
    // Add request-specific values
    ctx = context.WithValue(ctx, "requestID", generateRequestID())
    ctx = context.WithValue(ctx, "userID", getUserID(r))
    
    // 使用 context 处理请求
    // Process request with context
    result, err := processRequest(ctx, r)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusRequestTimeout)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(result)
}

func processRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    // 从 context 提取值
    // Extract values from context
    requestID := ctx.Value("requestID").(string)
    userID := ctx.Value("userID").(string)
    
    log.Printf("Processing request %s for user %s", requestID, userID)
    
    // 使用 context 进行数据库调用
    // Make database call with context
    data, err := fetchFromDatabase(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 使用 context 进行外部 API 调用
    // Make external API call with context
    enriched, err := enrichData(ctx, data)
    if err != nil {
        return nil, err
    }
    
    return enriched, nil
}
```

### 带优雅关闭的工作池

```go
type WorkerPool struct {
    workers int
    jobs    chan Job
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, 100),
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    log.Printf("Worker %d started", id)
    defer log.Printf("Worker %d stopped", id)
    
    for {
        select {
        case <-wp.ctx.Done():
            log.Printf("Worker %d shutting down: %v", id, wp.ctx.Err())
            return
        case job := <-wp.jobs:
            wp.processJob(job)
        }
    }
}

func (wp *WorkerPool) processJob(job Job) {
    // 为此作业创建带超时的 context
    // Create context with timeout for this job
    ctx, cancel := context.WithTimeout(wp.ctx, job.Timeout)
    defer cancel()
    
    err := job.Execute(ctx)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            log.Printf("Job %s timed out", job.ID)
        } else {
            log.Printf("Job %s failed: %v", job.ID, err)
        }
    }
}

func (wp *WorkerPool) Shutdown() {
    wp.cancel() // 这将导致所有 worker 停止
               // This will cause all workers to stop
}
```

### 带 Context 的数据库操作

```go
func getUserOrders(ctx context.Context, db *sql.DB, userID string) ([]Order, error) {
    // 使用 context 创建查询
    // Create query with context
    query := `
        SELECT id, user_id, product_name, amount, created_at 
        FROM orders 
        WHERE user_id = $1 
        ORDER BY created_at DESC
    `
    
    // 使用 context 执行查询（如果 context 被取消将被取消）
    // Execute query with context (will be cancelled if context is cancelled)
    rows, err := db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    var orders []Order
    for rows.Next() {
        // 处理行时检查取消
        // Check for cancellation while processing rows
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        
        var order Order
        err := rows.Scan(&order.ID, &order.UserID, &order.ProductName, &order.Amount, &order.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        orders = append(orders, order)
    }
    
    return orders, nil
}
```

## Context 最佳实践

### ✅ 应该做的：

**将 context 作为第一个参数传递**

```go
func ProcessData(ctx context.Context, data []byte) error // ✅ 好的
func ProcessData(data []byte, ctx context.Context) error // ❌ 坏的
```

**始终调用 cancel() 以防止内存泄漏**

```go
ctx, cancel := context.WithTimeout(parent, 30*time.Second)
defer cancel() // ✅ 始终这样做
               // ✅ Always do this
```

**在循环和长操作中检查 ctx.Done()**

```go
for i := 0; i < len(items); i++ {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    processItem(items[i])
}
```

**使用 context 传递请求范围的值**

```go
ctx = context.WithValue(ctx, "traceID", "abc123")
ctx = context.WithValue(ctx, "userID", "user456")
```

**从父 context 派生子 context**

```go
childCtx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
```

### ❌ 不应该做的：

**不要在结构体中存储 context（少数例外情况除外）**

```go
// ❌ 坏的 - context 存储在结构体中
// ❌ Bad - context stored in struct
type Server struct {
    ctx context.Context
}

// ✅ 好的 - context 作为参数传递
// ✅ Good - context passed as parameter
func (s *Server) ProcessRequest(ctx context.Context) error
```

**不要传递 nil context**

```go
ProcessData(nil, data) // ❌ 坏的
ProcessData(context.Background(), data) // ✅ 好的
```

**不要使用 context 传递可选参数**

```go
// ❌ 坏的 - 使用 context 传递配置
// ❌ Bad - using context for config
ctx = context.WithValue(ctx, "retryCount", 3)

// ✅ 好的 - 使用结构体传递配置
// ✅ Good - use struct for config
type Config struct {
    RetryCount int
}
func ProcessData(ctx context.Context, cfg Config) error
```

**不要忽略 context 取消**

```go
// ❌ 坏的 - 忽略 context
// ❌ Bad - ignoring context
func doWork(ctx context.Context) {
    for i := 0; i < 1000; i++ {
        // 没有 context 检查
        // No context checking
        time.Sleep(100 * time.Millisecond)
    }
}

// ✅ 好的 - 尊重 context
// ✅ Good - respecting context
func doWork(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        time.Sleep(100 * time.Millisecond)
    }
    return nil
}
```

## 常见错误和解决方案

### 错误 1：不调用 Cancel 导致的内存泄漏

```go
// ❌ 内存泄漏 - cancel 未被调用
// ❌ Memory leak - cancel not called
func badExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    // cancel() 永远不会被调用 - goroutine 和计时器泄漏！
    // cancel() never called - goroutine and timer leak!
    doWork(ctx)
}

// ✅ 修复 - 始终调用 cancel
// ✅ Fixed - always call cancel
func goodExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // 始终调用 cancel
                  // Always call cancel
    doWork(ctx)
}
```

### 错误 2：Context 值的竞争条件

```go
// ❌ 竞争条件 - 值可能改变
// ❌ Race condition - value might change
func badExample(ctx context.Context) {
    go func() {
        userID := ctx.Value("userID").(string) // 如果为 nil 可能 panic
                                              // Might panic if nil
        processUser(userID)
    }()
}

// ✅ 安全的值提取
// ✅ Safe value extraction
func goodExample(ctx context.Context) {
    userIDValue := ctx.Value("userID")
    if userIDValue == nil {
        return // 处理缺失值
               // Handle missing value
    }
    userID, ok := userIDValue.(string)
    if !ok {
        return // 处理错误类型
               // Handle wrong type
    }
    
    go func() {
        processUser(userID)
    }()
}
```

### 错误 3：Context 继承问题

```go
// ❌ 坏的 - 创建独立的 context
// ❌ Bad - creating independent contexts
func badChain() {
    ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel1()
    
    ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second) // 独立的！
                                                                             // Independent!
    defer cancel2()
    
    doWork(ctx2) // 不会继承 ctx1 的取消
                // Won't inherit ctx1's cancellation
}

// ✅ 好的 - 正确的 context 链
// ✅ Good - proper context chaining
func goodChain() {
    ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel1()
    
    ctx2, cancel2 := context.WithTimeout(ctx1, 5*time.Second) // 从 ctx1 继承
                                                             // Inherits from ctx1
    defer cancel2()
    
    doWork(ctx2) // 当 ctx1 或 ctx2 超时时将被取消
                // Will be cancelled when ctx1 OR ctx2 times out
}
```

## 使用 Context 进行测试

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    err := doSlowWork(ctx)
    if err != context.DeadlineExceeded {
        t.Errorf("Expected timeout, got %v", err)
    }
}

func TestWithCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(50 * time.Millisecond)
        cancel()
    }()
    
    err := doWork(ctx)
    if err != context.Canceled {
        t.Errorf("Expected cancellation, got %v", err)
    }
}
```

## 高级 Context 模式

### 自定义 Context 类型（高级）

```go
type contextKey string

const (
    RequestIDKey contextKey = "requestID"
    UserIDKey   contextKey = "userID"
)

// 类型安全的 context 辅助函数
// Type-safe context helpers
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
    requestID, ok := ctx.Value(RequestIDKey).(string)
    return requestID, ok
}
```

### Context 中间件

```go
func contextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 向 context 添加请求 ID
        // Add request ID to context
        requestID := generateRequestID()
        ctx := WithRequestID(r.Context(), requestID)
        
        // 向 context 添加用户信息
        // Add user info to context
        if userID := getUserFromAuth(r); userID != "" {
            ctx = context.WithValue(ctx, UserIDKey, userID)
        }
        
        // 设置响应头
        // Set response header
        w.Header().Set("X-Request-ID", requestID)
        
        // 使用丰富的 context 调用下一个处理器
        // Call next handler with enriched context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## 性能考虑

- Context 开销很小 - 正常使用不用担心性能
- 避免过度的 context 链 - 每个 WithValue 都会创建新的 context
- 谨慎使用 context 值 - 它们不是为大数据设计的
- 在热路径中小心使用 context - 如果怀疑有问题请进行分析

## 进一步学习资源

- Go Context 包文档
- Go 博客：Go 并发模式：Context
- Effective Go：并发
- Go Wiki：Context

## 总结

context 包对于以下方面至关重要：

- **取消**：当不再需要时停止工作
- **超时**：防止操作运行太长时间
- **请求范围值**：跨函数边界传递数据
- **优雅关闭**：协调 goroutine 间的清理

掌握这些模式，你将编写更健壮、可维护的 Go 应用程序！