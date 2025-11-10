# 挑战 7：带错误处理的银行账户

相关学习资料: ./bank_account_errors.md

## 问题描述

实现一个带有完善错误处理机制的简单银行系统。你需要创建一个 `BankAccount` 结构体来管理账户余额操作，并实现相应的错误处理。

---

## 要求

实现一个包含以下字段的 `BankAccount` 结构体：

* **ID (string)**：账户唯一标识符
* **Owner (string)**：账户所有者姓名
* **Balance (float64)**：账户当前余额
* **MinBalance (float64)**：账户必须保持的最低余额

实现以下方法：

* `NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error)`：构造函数，用于验证输入参数并创建账户
* `Deposit(amount float64) error`：向账户存款
* `Withdraw(amount float64) error`：从账户取款
* `Transfer(amount float64, target *BankAccount) error`：从一个账户向另一个账户转账

你必须实现自定义错误类型：

* **InsufficientFundsError**：当取款或转账导致余额低于最低余额时触发
* **NegativeAmountError**：当存款、取款或转账金额为负数时触发
* **ExceedsLimitError**：当存取款金额超过设定限额时触发
* **AccountError**：通用账户错误类型，包含合适的子类型

---

## 函数签名

```go
// 构造函数
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error)

// 方法
func (a *BankAccount) Deposit(amount float64) error
func (a *BankAccount) Withdraw(amount float64) error
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error

// 错误类型
type AccountError struct {
    // 实现自定义错误类型及必要字段
}

type InsufficientFundsError struct {
    // 实现自定义错误类型及必要字段
}

type NegativeAmountError struct {
    // 实现自定义错误类型及必要字段
}

type ExceedsLimitError struct {
    // 实现自定义错误类型及必要字段
}

// 每种错误类型都需要实现 Error() 方法
func (e *AccountError) Error() string
func (e *InsufficientFundsError) Error() string
func (e *NegativeAmountError) Error() string
func (e *ExceedsLimitError) Error() string
```

---

## 约束条件

* 所有金额必须为有效的非负值。
* 取款或转账操作不得使账户余额低于最低余额。
* 为存取款操作定义合理的限额（例如：$10,000）。
* 错误消息应具有描述性，并包含相关信息。
* 所有操作应是**线程安全的**（需使用合适的同步机制）。

---

## 示例用法

```go
// 创建新账户
account1, err := NewBankAccount("ACC001", "Alice", 1000.0, 100.0)
if err != nil {
    // 处理错误
}

account2, err := NewBankAccount("ACC002", "Bob", 500.0, 50.0)
if err != nil {
    // 处理错误
}

// 存款操作
if err := account1.Deposit(200.0); err != nil {
    // 处理错误
}

// 取款操作
if err := account1.Withdraw(50.0); err != nil {
    // 处理错误
}

// 转账操作
if err := account1.Transfer(300.0, account2); err != nil {
    // 处理错误
}
```

---

此挑战旨在帮助你深入理解 Go 中的**错误设计模式**、**结构体方法实现**以及**并发安全性处理**。
