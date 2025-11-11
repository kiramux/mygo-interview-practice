# Error handling and Go
# Go语言的错误处理

**Author:** Andrew Gerrand | **作者：** Andrew Gerrand  
**Date:** July 12, 2011 | **日期：** 2011年7月12日  
**Source:** https://go.dev/blog/error-handling-and-go  
**来源：** https://go.dev/blog/error-handling-and-go

---

## Introduction | 引言

If you have written any Go code you have probably encountered the built-in `error` type. Go code uses `error` values to indicate an abnormal state. For example, the `os.Open` function returns a non-nil `error` value when it fails to open a file.

如果您编写过任何Go代码，您可能遇到过内置的`error`类型。Go代码使用`error`值来表示异常状态。例如，`os.Open`函数在无法打开文件时会返回非nil的`error`值。

```go
func Open(name string) (file *File, err error)
```

The following code uses `os.Open` to open a file. If an error occurs it calls `log.Fatal` to print the error message and stop.

下面的代码使用`os.Open`来打开文件。如果发生错误，它会调用`log.Fatal`打印错误消息并停止程序。

```go
f, err := os.Open("filename.ext")
if err != nil {
    log.Fatal(err)
}
// do something with the open *File f
// 对打开的*File f进行一些操作
```

You can get a lot done in Go knowing just this about the `error` type, but in this article we'll take a closer look at `error` and discuss some good practices for error handling in Go.

仅仅了解`error`类型的这些知识，您就可以在Go中完成很多工作，但在本文中，我们将更仔细地研究`error`，并讨论Go中错误处理的一些良好实践。

## The error type | error类型

The `error` type is an interface type. An `error` variable represents any value that can describe itself as a string. Here is the interface's declaration:

`error`类型是一个接口类型。`error`变量表示任何可以将自己描述为字符串的值。以下是接口的声明：

```go
type error interface {
    Error() string
}
```

The `error` type, as with all built in types, is predeclared in the universe block.

`error`类型，与所有内置类型一样，在全局块中预先声明。

The most commonly-used `error` implementation is the `errors` package's unexported `errorString` type.

最常用的`error`实现是`errors`包的未导出的`errorString`类型。

```go
// errorString is a trivial implementation of error.
// errorString是error的一个简单实现。
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

You can construct one of these values with the `errors.New` function. It takes a string that it converts to an `errors.errorString` and returns as an `error` value.

您可以使用`errors.New`函数构造这些值之一。它接受一个字符串，将其转换为`errors.errorString`并作为`error`值返回。

```go
// New returns an error that formats as the given text.
// New返回一个格式化为给定文本的错误。
func New(text string) error {
    return &errorString{text}
}
```

Here's how you might use `errors.New`:

以下是您可能使用`errors.New`的方式：

```go
func Sqrt(f float64) (float64, error) {
    if f < 0 {
        return 0, errors.New("math: square root of negative number")
    }
    // implementation
    // 实现
}
```

A caller passing a negative argument to `Sqrt` receives a non-nil `error` value (whose concrete representation is an `errors.errorString` value). The caller can access the error string ("math: square root of...") by calling the `error`'s `Error` method, or by just printing it:

向`Sqrt`传递负参数的调用者会收到一个非nil的`error`值（其具体表示是`errors.errorString`值）。调用者可以通过调用`error`的`Error`方法，或者仅仅打印它来访问错误字符串（"math: square root of..."）：

```go
f, err := Sqrt(-1)
if err != nil {
    fmt.Println(err)
}
```

The `fmt` package formats an `error` value by calling its `Error() string` method.

`fmt`包通过调用其`Error() string`方法来格式化`error`值。

It is the error implementation's responsibility to summarize the context. The error returned by `os.Open` formats as "open /etc/passwd: permission denied," not just "permission denied." The error returned by our `Sqrt` is missing information about the invalid argument.

错误实现负责总结上下文。`os.Open`返回的错误格式为"open /etc/passwd: permission denied"，而不仅仅是"permission denied"。我们的`Sqrt`返回的错误缺少关于无效参数的信息。

To add that information, a useful function is the `fmt` package's `Errorf`. It formats a string according to `Printf`'s rules and returns it as an `error` created by `errors.New`.

要添加这些信息，一个有用的函数是`fmt`包的`Errorf`。它根据`Printf`的规则格式化字符串，并将其作为由`errors.New`创建的`error`返回。

```go
if f < 0 {
    return 0, fmt.Errorf("math: square root of negative number %g", f)
}
```

In many cases `fmt.Errorf` is good enough, but since `error` is an interface, you can use arbitrary data structures as error values, to allow callers to inspect the details of the error.

在许多情况下，`fmt.Errorf`已经足够好了，但由于`error`是一个接口，您可以使用任意数据结构作为错误值，以允许调用者检查错误的详细信息。

For instance, our hypothetical callers might want to recover the invalid argument passed to `Sqrt`. We can enable that by defining a new error implementation instead of using `errors.errorString`:

例如，我们假设的调用者可能想要恢复传递给`Sqrt`的无效参数。我们可以通过定义新的错误实现而不是使用`errors.errorString`来实现这一点：

```go
type NegativeSqrtError float64

func (f NegativeSqrtError) Error() string {
    return fmt.Sprintf("math: square root of negative number %g", float64(f))
}
```

A sophisticated caller can then use a type assertion to check for a `NegativeSqrtError` and handle it specially, while callers that just pass the error to `fmt.Println` or `log.Fatal` will see no change in behavior.

然后，复杂的调用者可以使用类型断言来检查`NegativeSqrtError`并特殊处理它，而那些只是将错误传递给`fmt.Println`或`log.Fatal`的调用者在行为上不会看到任何变化。

As another example, the `json` package specifies a `SyntaxError` type that the `json.Decode` function returns when it encounters a syntax error parsing a JSON blob.

作为另一个例子，`json`包指定了一个`SyntaxError`类型，当`json.Decode`函数在解析JSON数据时遇到语法错误时会返回该类型。

```go
type SyntaxError struct {
    msg    string // description of error | 错误描述
    Offset int64  // error occurred after reading Offset bytes | 读取Offset字节后发生错误
}

func (e *SyntaxError) Error() string { return e.msg }
```

The `Offset` field isn't even shown in the default formatting of the error, but callers can use it to add file and line information to their error messages:

`Offset`字段甚至不会在错误的默认格式中显示，但调用者可以使用它来向错误消息添加文件和行信息：

```go
if err := dec.Decode(&val); err != nil {
    if serr, ok := err.(*json.SyntaxError); ok {
        line, col := findLine(f, serr.Offset)
        return fmt.Errorf("%s:%d:%d: %v", f.Name(), line, col, err)
    }
    return err
}
```

(This is a slightly simplified version of some actual code from the Camlistore project.)

（这是来自Camlistore项目的一些实际代码的稍微简化版本。）

The `error` interface requires only a `Error` method; specific error implementations might have additional methods. For instance, the `net` package returns errors of type `error`, following the usual convention, but some of the error implementations have additional methods defined by the `net.Error` interface:

`error`接口只需要一个`Error`方法；特定的错误实现可能有额外的方法。例如，`net`包按照通常的约定返回`error`类型的错误，但一些错误实现有由`net.Error`接口定义的额外方法：

```go
package net

type Error interface {
    error
    Timeout() bool   // Is the error a timeout? | 错误是否为超时？
    Temporary() bool // Is the error temporary? | 错误是否为临时错误？
}
```

Client code can test for a `net.Error` with a type assertion and then distinguish transient network errors from permanent ones. For instance, a web crawler might sleep and retry when it encounters a temporary error and give up otherwise.

客户端代码可以通过类型断言测试`net.Error`，然后区分临时网络错误和永久错误。例如，网页爬虫在遇到临时错误时可能会休眠并重试，否则就放弃。

```go
if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
    time.Sleep(1e9)
    continue
}
if err != nil {
    log.Fatal(err)
}
```

## Simplifying repetitive error handling | 简化重复的错误处理

In Go, error handling is important. The language's design and conventions encourage you to explicitly check for errors where they occur (as distinct from the convention in other languages of throwing exceptions and sometimes catching them). In some cases this makes Go code verbose, but fortunately there are some techniques you can use to minimize repetitive error handling.

在Go中，错误处理很重要。该语言的设计和约定鼓励您在错误发生的地方明确检查错误（这与其他语言抛出异常并有时捕获它们的约定不同）。在某些情况下，这使Go代码变得冗长，但幸运的是，有一些技术可以用来最小化重复的错误处理。

Consider an App Engine application with an HTTP handler that retrieves a record from the datastore and formats it with a template.

考虑一个App Engine应用程序，它有一个HTTP处理程序，从数据存储中检索记录并用模板格式化它。

```go
func init() {
    http.HandleFunc("/view", viewRecord)
}

func viewRecord(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if err := viewTemplate.Execute(w, record); err != nil {
        http.Error(w, err.Error(), 500)
    }
}
```

This function handles errors returned by the `datastore.Get` function and `viewTemplate`'s `Execute` method. In both cases, it presents a simple error message to the user with the HTTP status code 500 ("Internal Server Error"). This looks like a manageable amount of code, but add some more HTTP handlers and you quickly end up with many copies of identical error handling code.

这个函数处理由`datastore.Get`函数和`viewTemplate`的`Execute`方法返回的错误。在两种情况下，它都向用户呈现带有HTTP状态码500（"Internal Server Error"）的简单错误消息。这看起来像是可管理的代码量，但添加更多的HTTP处理程序，您很快就会得到许多相同错误处理代码的副本。

To reduce the repetition we can define our own HTTP `appHandler` type that includes an `error` return value:

为了减少重复，我们可以定义自己的HTTP `appHandler`类型，其中包含`error`返回值：

```go
type appHandler func(http.ResponseWriter, *http.Request) error
```

Then we can change our `viewRecord` function to return errors:

然后我们可以更改我们的`viewRecord`函数来返回错误：

```go
func viewRecord(w http.ResponseWriter, r *http.Request) error {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        return err
    }
    return viewTemplate.Execute(w, record)
}
```

This is simpler than the original version, but the `http` package doesn't understand functions that return `error`. To fix this we can implement the `http.Handler` interface's `ServeHTTP` method on `appHandler`:

这比原始版本更简单，但`http`包不理解返回`error`的函数。为了解决这个问题，我们可以在`appHandler`上实现`http.Handler`接口的`ServeHTTP`方法：

```go
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := fn(w, r); err != nil {
        http.Error(w, err.Error(), 500)
    }
}
```

The `ServeHTTP` method calls the `appHandler` function and displays the returned error (if any) to the user. Notice that the method's receiver, `fn`, is a function. (Go can do that!) The method invokes the function by calling the receiver in the expression `fn(w, r)`.

`ServeHTTP`方法调用`appHandler`函数并向用户显示返回的错误（如果有的话）。注意，该方法的接收器`fn`是一个函数。（Go可以做到这一点！）该方法通过在表达式`fn(w, r)`中调用接收器来调用函数。

Now when registering `viewRecord` with the http package we use the `Handle` function (instead of `HandleFunc`) as `appHandler` is an `http.Handler` (not an `http.HandlerFunc`).

现在，当向http包注册`viewRecord`时，我们使用`Handle`函数（而不是`HandleFunc`），因为`appHandler`是一个`http.Handler`（不是`http.HandlerFunc`）。

```go
func init() {
    http.Handle("/view", appHandler(viewRecord))
}
```

With this basic error handling infrastructure in place, we can make it more user friendly. Rather than just displaying the error string, it would be better to give the user a simple error message with an appropriate HTTP status code, while logging the full error to the App Engine developer console for debugging purposes.

有了这个基本的错误处理基础设施，我们可以让它更用户友好。与其仅仅显示错误字符串，不如给用户一个带有适当HTTP状态码的简单错误消息，同时将完整的错误记录到App Engine开发者控制台以用于调试目的。

To do this we create an `appError` struct containing an `error` and some other fields:

为了做到这一点，我们创建一个包含`error`和一些其他字段的`appError`结构体：

```go
type appError struct {
    Error   error
    Message string
    Code    int
}
```

Next we modify the appHandler type to return `*appError` values:

接下来我们修改appHandler类型以返回`*appError`值：

```go
type appHandler func(http.ResponseWriter, *http.Request) *appError
```

(It's usually a mistake to pass back the concrete type of an error rather than `error`, for reasons discussed in the Go FAQ, but it's the right thing to do here because `ServeHTTP` is the only place that sees the value and uses its contents.)

（通常返回错误的具体类型而不是`error`是错误的，原因在Go FAQ中有讨论，但在这里这样做是正确的，因为`ServeHTTP`是唯一看到该值并使用其内容的地方。）

And make `appHandler`'s `ServeHTTP` method display the `appError`'s `Message` to the user with the correct HTTP status `Code` and log the full `Error` to the developer console:

并使`appHandler`的`ServeHTTP`方法向用户显示`appError`的`Message`，使用正确的HTTP状态`Code`，并将完整的`Error`记录到开发者控制台：

```go
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { // e is *appError, not os.Error. | e是*appError，不是os.Error。
        c := appengine.NewContext(r)
        c.Errorf("%v", e.Error)
        http.Error(w, e.Message, e.Code)
    }
}
```

Finally, we update `viewRecord` to the new function signature and have it return more context when it encounters an error:

最后，我们将`viewRecord`更新为新的函数签名，并在遇到错误时让它返回更多上下文：

```go
func viewRecord(w http.ResponseWriter, r *http.Request) *appError {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        return &appError{err, "Record not found", 404}
    }
    if err := viewTemplate.Execute(w, record); err != nil {
        return &appError{err, "Can't display record", 500}
    }
    return nil
}
```

This version of `viewRecord` is the same length as the original, but now each of those lines has specific meaning and we are providing a friendlier user experience.

这个版本的`viewRecord`与原始版本长度相同，但现在每一行都有特定的含义，我们提供了更友好的用户体验。

It doesn't end there; we can further improve the error handling in our application. Some ideas:

这还没有结束；我们可以进一步改进应用程序中的错误处理。一些想法：

- give the error handler a pretty HTML template,

  为错误处理程序提供一个漂亮的HTML模板，

- make debugging easier by writing the stack trace to the HTTP response when the user is an administrator,

  当用户是管理员时，通过将堆栈跟踪写入HTTP响应来简化调试，

- write a constructor function for `appError` that stores the stack trace for easier debugging,

  为`appError`编写一个构造函数，存储堆栈跟踪以便于调试，

- recover from panics inside the `appHandler`, logging the error to the console as "Critical," while telling the user "a serious error has occurred." This is a nice touch to avoid exposing the user to inscrutable error messages caused by programming errors. See the Defer, Panic, and Recover article for more details.

  从`appHandler`内部的panic中恢复，将错误记录到控制台为"Critical"，同时告诉用户"发生了严重错误"。这是一个很好的做法，可以避免向用户暴露由编程错误引起的难以理解的错误消息。有关更多详细信息，请参阅Defer, Panic, and Recover文章。

## Conclusion | 结论

Proper error handling is an essential requirement of good software. By employing the techniques described in this post you should be able to write more reliable and succinct Go code.

适当的错误处理是好软件的基本要求。通过使用本文中描述的技术，您应该能够编写更可靠和简洁的Go代码。

---

**Related articles | 相关文章：**
- Next article: Go for App Engine is now generally available | 下一篇文章：Go for App Engine现已正式可用
- Previous article: First Class Functions in Go | 上一篇文章：Go中的一等函数
- [Blog Index | 博客索引](https://go.dev/blog/all)