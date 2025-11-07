**Go 中的 SQL 数据库操作**

Go 提供了出色的支持用于 SQL 数据库的功能，通过标准的 `database/sql` 包。这个挑战专注于实现 SQLite 中的 CRUD 操作，但这些概念同样适用于其他 SQL 数据库，如 MySQL、PostgreSQL 等。

**database/sql 包**

`database/sql` 包提供了一个 SQL（或类似 SQL）数据库的通用接口。它的功能包括：

- 管理连接池
- 处理事务
- 提供预编译语句
- 为各种数据库提供驱动程序

``` go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // 注意使用了下划线导入
)
```

下划线导入（`_`）用于仅为了副作用导入一个包（在此情况下是注册数据库驱动程序）。

**打开数据库连接**

``` go
db, err := sql.Open("sqlite3", "path/to/database.db")
if err != nil {
    return nil, err
}

// 测试连接
if err = db.Ping(); err != nil {
    return nil, err
}

return db, nil
```

`sql.Open()` 函数并不会立即建立与数据库的连接，只有在调用 `Ping()` 方法或执行查询时，连接才会被建立。

**执行简单查询**

可以使用 `db.Exec()` 执行不返回结果集的简单 SQL 语句：

``` go
result, err := db.Exec(
    "CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)"
)
if err != nil {
    return err
}
```

对于 `INSERT` 操作，可以通过 `LastInsertId()` 获取插入的最后一个 ID：

``` go
result, err := db.Exec(
    "INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
    product.Name, product.Price, product.Quantity, product.Category
)
if err != nil {
    return err
}

// 获取插入行的 ID
id, err := result.LastInsertId()
if err != nil {
    return err
}
product.ID = id
```

**查询数据**

使用 `db.Query()` 或 `db.QueryRow()` 来查询并检索数据：

``` go
// 查询多行
rows, err := db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
if err != nil {
    return nil, err
}
defer rows.Close() // 使用完 rows 后总是要关闭

var products []*Product
for rows.Next() {
    p := &Product{}
    err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
    if err != nil {
        return nil, err
    }
    products = append(products, p)
}

// 检查迭代 rows 时的错误
if err = rows.Err(); err != nil {
    return nil, err
}

return products, nil
```

对于单行数据，使用 `QueryRow()`：

``` go
row := db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)

p := &Product{}
err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("product with ID %d not found", id)
    }
    return nil, err
}

return p, nil
```

**预编译语句**

对于多次执行的查询，可以使用预编译语句来提高性能：

``` go
stmt, err := db.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
if err != nil {
    return err
}
defer stmt.Close()

for id, quantity := range updates {
    _, err := stmt.Exec(quantity, id)
    if err != nil {
        return err
    }
}
```

**事务**

事务确保一组操作要么全部成功，要么全部失败：

``` go
// 开始一个事务
tx, err := db.Begin()
if err != nil {
    return err
}
defer func() {
    if err != nil {
        tx.Rollback() // 出现错误时回滚事务
    }
}()

stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
if err != nil {
    return err
}
defer stmt.Close()

for id, quantity := range updates {
    result, err := stmt.Exec(quantity, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product with ID %d not found", id)
    }
}

// 提交事务，在 defer 处已经写了出现错误回滚事务
return tx.Commit()
```

**参数绑定与 SQL 注入预防**

始终使用参数绑定，而不是字符串拼接，以防止 SQL 注入：

``` go
// 不要这样做 - 容易受到 SQL 注入攻击
query := fmt.Sprintf("SELECT * FROM products WHERE category = '%s'", category)

// 应该这样做 - 使用参数绑定
rows, err := db.Query("SELECT * FROM products WHERE category = ?", category)
```

不同的数据库驱动程序使用不同的占位符样式：

- SQLite, MySQL: `?`
- PostgreSQL: `$1`, `$2` 等
- Oracle: `:name`

**处理 NULL 值**

SQL 数据库可能包含 `NULL` 值，Go 的 `database/sql` 包提供了特殊类型来处理这些值：

``` go
import (
    "database/sql"
)

type Product struct {
    ID       int64
    Name     string
    Price    float64
    Quantity int
    Category sql.NullString // 可能为 NULL
}

// 扫描时
var category sql.NullString
err := row.Scan(&id, &name, &price, &quantity, &category)

// 使用时
if category.Valid {
    fmt.Println(category.String)
} else {
    fmt.Println("Category is NULL")
}
```

**错误处理**

有几种常见的错误类型需要检查：

``` go
if err == sql.ErrNoRows {
    // 没有返回行（不一定是错误）
    return nil, fmt.Errorf("product not found")
}

// 检查唯一约束冲突
if strings.Contains(err.Error(), "UNIQUE constraint failed") {
    return nil, fmt.Errorf("product with that name already exists")
}
```

**连接池**

`database/sql` 包自动处理连接池。你可以控制连接池的行为：

``` go
db.SetMaxOpenConns(25)  // 最大打开连接数
db.SetMaxIdleConns(25)  // 最大空闲连接数
db.SetConnMaxLifetime(5 * time.Minute) // 连接最大可复用时间
```

**最佳实践**

- 始终关闭像 `rows` 和 `stmt` 这样的资源
- 对必须作为一个组成功的操作使用事务
- 切勿通过字符串拼接构建 SQL 查询
- 检查特定的错误，例如 `sql.ErrNoRows`
- 保持数据库连接在应用程序的生命周期内开启
- 对于复杂的应用程序，考虑使用 ORM 或查询构建器