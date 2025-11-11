# Don't just check errors, handle them gracefully
# 不要仅仅检查错误，要优雅地处理它们

**Author:** Dave Cheney | **作者：** Dave Cheney  
**Date:** April 27, 2016 | **日期：** 2016年4月27日  
**Source:** https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully  
**来源：** https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

---

This post is an extract from my presentation at the recent GoCon spring conference in Tokyo, Japan.

这篇文章摘录自我在最近于日本东京举办的GoCon春季会议上的演讲。

---

## Errors are just values | 错误只是值

I've spent a lot of time thinking about the best way to handle errors in Go programs. I really wanted there to be a single way to do error handling, something that we could teach all Go programmers by rote, just as we might teach mathematics, or the alphabet.

我花了很多时间思考在Go程序中处理错误的最佳方式。我真的希望有一种处理错误的单一方法，我们可以像教授数学或字母表那样，通过死记硬背的方式教给所有Go程序员。

However, I have concluded that there is no single way to handle errors. Instead, I believe Go's error handling can be classified into the three core strategies.

然而，我得出的结论是没有处理错误的单一方法。相反，我相信Go的错误处理可以分为三种核心策略。

## Sentinel errors | 哨兵错误

The first category of error handling is what I call *sentinel errors*.

第一类错误处理是我称为*哨兵错误*的方式。

```go
if err == ErrSomething { … }
```

The name descends from the practice in computer programming of using a specific value to signify that no further processing is possible. So to with Go, we use specific values to signify an error.

这个名称来源于计算机编程中使用特定值来表示无法进行进一步处理的实践。在Go中也是如此，我们使用特定的值来表示错误。

Examples include values like `io.EOF` or low level errors like the constants in the `syscall` package, like `syscall.ENOENT`.

示例包括像`io.EOF`这样的值，或者像`syscall`包中的常量这样的底层错误，比如`syscall.ENOENT`。

There are even sentinel errors that signify that an error *did not* occur, like `go/build.NoGoError`, and `path/filepath.SkipDir` from `path/filepath.Walk`.

甚至还有表示错误*没有*发生的哨兵错误，比如`go/build.NoGoError`，以及来自`path/filepath.Walk`的`path/filepath.SkipDir`。

Using sentinel values is the least flexible error handling strategy, as the caller must compare the result to predeclared value using the equality operator. This presents a problem when you want to provide more context, as returning a different error would will break the equality check.

使用哨兵值是最不灵活的错误处理策略，因为调用者必须使用等号运算符将结果与预先声明的值进行比较。当你想提供更多上下文时，这会产生问题，因为返回不同的错误会破坏等值检查。

Even something as well meaning as using `fmt.Errorf` to add some context to the error will defeat the caller's equality test. Instead the caller will be forced to look at the output of the `error`'s `Error` method to see if it matches a specific string.

即使是像使用`fmt.Errorf`为错误添加一些上下文这样善意的做法，也会破坏调用者的等值测试。相反，调用者将被迫查看`error`的`Error`方法的输出，来看它是否匹配特定的字符串。

### Never inspect the output of error.Error | 永远不要检查error.Error的输出

As an aside, I believe you should never inspect the output of the `error.Error` method. The `Error` method on the `error` interface exists for humans, not code.

顺便说一下，我认为你永远不应该检查`error.Error`方法的输出。`error`接口上的`Error`方法是为人类而存在的，而不是为代码。

The contents of that string belong in a log file, or displayed on screen. You shouldn't try to change the behaviour of your program by inspecting it.

该字符串的内容属于日志文件，或在屏幕上显示。你不应该试图通过检查它来改变程序的行为。

I know that sometimes this isn't possible, and as someone pointed out on twitter, this advice doesn't apply to writing tests. Never the less, comparing the string form of an error is, in my opinion, a code smell, and you should try to avoid it.

我知道有时这是不可能的，正如有人在Twitter上指出的，这个建议不适用于编写测试。尽管如此，在我看来，比较错误的字符串形式是一种代码异味，你应该尽量避免它。

### Sentinel errors become part of your public API | 哨兵错误成为你公共API的一部分

If your public function or method returns an error of a particular value then that value must be public, and of course documented. This adds to the surface area of your API.

如果你的公共函数或方法返回特定值的错误，那么该值必须是公共的，当然还要有文档记录。这增加了你API的表面积。

If your API defines an interface which returns a specific error, all implementations of that interface will be restricted to returning only that error, even if they could provide a more descriptive error.

如果你的API定义了一个返回特定错误的接口，该接口的所有实现都将被限制为只返回该错误，即使它们可以提供更具描述性的错误。

We see this with `io.Reader`. Functions like `io.Copy` require a reader implementation to return *exactly* `io.EOF` to signal to the caller *no more data, but that isn't an error*.

我们在`io.Reader`中看到了这一点。像`io.Copy`这样的函数要求reader实现返回*确切的*`io.EOF`来向调用者发出信号：*没有更多数据，但这不是错误*。

### Sentinel errors create a dependency between two packages | 哨兵错误在两个包之间创建依赖关系

By far the worst problem with sentinel error values is they create a source code dependency between two packages. As an example, to check if an error is equal to `io.EOF`, your code must import the `io` package.

到目前为止，哨兵错误值最严重的问题是它们在两个包之间创建了源代码依赖关系。例如，要检查错误是否等于`io.EOF`，你的代码必须导入`io`包。

This specific example does not sound so bad, because it is quite common, but imagine the coupling that exists when many packages in your project export error values, which other packages in your project must import to check for specific error conditions.

这个特定的例子听起来不那么糟糕，因为它很常见，但想象一下当项目中的许多包导出错误值时存在的耦合，项目中的其他包必须导入这些包来检查特定的错误条件。

Having worked in a large project that toyed with this pattern, I can tell you that the spectre of bad design–in the form of an import loop–was never far from our minds.

在一个使用这种模式的大型项目中工作过，我可以告诉你，糟糕设计的阴影——以导入循环的形式——从未远离我们的思维。

### Conclusion: avoid sentinel errors | 结论：避免哨兵错误

So, my advice is to avoid using sentinel error values in the code you write. There are a few cases where they are used in the standard library, but this is not a pattern that you should emulate.

因此，我的建议是在你编写的代码中避免使用哨兵错误值。在标准库中有一些使用它们的情况，但这不是你应该模仿的模式。

If someone asks you to export an error value from your package, you should politely decline and instead suggest an alternative method, such as the ones I will discuss next.

如果有人要求你从包中导出错误值，你应该礼貌地拒绝，并建议其他方法，比如我接下来要讨论的方法。

## Error types | 错误类型

Error types are the second form of Go error handling I want to discuss.

错误类型是我想要讨论的Go错误处理的第二种形式。

```go
if err, ok := err.(SomeType); ok { … }
```

An error type is a type that you create that implements the error interface. In this example, the `MyError` type tracks the file and line, as well as a message explaining what happened.

错误类型是你创建的实现error接口的类型。在这个例子中，`MyError`类型跟踪文件和行，以及解释发生了什么的消息。

```go
type MyError struct {
        Msg string
        File string
        Line int
}

func (e *MyError) Error() string { 
        return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Msg)
}

return &MyError{"Something happened", "server.go", 42}
```

Because `MyError error` is a type, callers can use type assertion to extract the extra context from the error.

因为`MyError error`是一个类型，调用者可以使用类型断言从错误中提取额外的上下文。

```go
err := something()
switch err := err.(type) {
case nil:
        // call succeeded, nothing to do
        // 调用成功，无需处理
case *MyError:
        fmt.Println("error occurred on line:", err.Line)
default:
        // unknown error
        // 未知错误
}
```

A big improvement of error types over error values is their ability to wrap an underlying error to provide more context.

错误类型相对于错误值的一大改进是它们能够包装底层错误以提供更多上下文。

An excellent example of this is the `os.PathError` type which annotates the underlying error with the operation it was trying to perform, and the file it was trying to use.

一个很好的例子是`os.PathError`类型，它用试图执行的操作和试图使用的文件来注释底层错误。

```go
// PathError records an error and the operation
// and file path that caused it.
// PathError记录错误以及导致错误的操作和文件路径。
type PathError struct {
        Op   string
        Path string
        Err  error // the cause | 原因
}

func (e *PathError) Error() string
```

### Problems with error types | 错误类型的问题

So the caller can use a type assertion or type switch, error types must be made public.

为了使调用者能够使用类型断言或类型switch，错误类型必须是公共的。

If your code implements an interface whose contract requires a specific error type, all implementors of that interface need to depend on the package that defines the error type.

如果你的代码实现了一个其契约需要特定错误类型的接口，该接口的所有实现者都需要依赖定义错误类型的包。

This intimate knowledge of a package's types creates a strong coupling with the caller, making for a brittle API.

这种对包类型的深入了解与调用者创建了强耦合，使API变得脆弱。

### Conclusion: avoid error types | 结论：避免错误类型

While error types are better than sentinel error values, because they can capture more context about what went wrong, error types share many of the problems of error values.

虽然错误类型比哨兵错误值更好，因为它们可以捕获更多关于出错内容的上下文，但错误类型与错误值有很多相同的问题。

So again my advice is to avoid error types, or at least, avoid making them part of your public API.

所以我的建议再次是避免错误类型，或者至少，避免让它们成为你公共API的一部分。

## Opaque errors | 不透明错误

Now we come to the third category of error handling. In my opinion this is the most flexible error handling strategy as it requires the least coupling between your code and caller.

现在我们来到错误处理的第三类。在我看来，这是最灵活的错误处理策略，因为它要求你的代码和调用者之间的耦合最少。

I call this style opaque error handling, because while you know an error occurred, you don't have the ability to see inside the error. As the caller, all you know about the result of the operation is that it worked, or it didn't.

我称这种风格为不透明错误处理，因为虽然你知道发生了错误，但你无法看到错误内部。作为调用者，你对操作结果所知道的就是它成功了，或者失败了。

This is all there is to opaque error handling–just return the error without assuming anything about its contents. If you adopt this position, then error handling can become significantly more useful as a debugging aid.

这就是不透明错误处理的全部内容——只是返回错误而不对其内容做任何假设。如果你采用这种立场，那么错误处理可以作为调试辅助工具变得更加有用。

```go
import "github.com/quux/bar"

func fn() error {
        x, err := bar.Foo()
        if err != nil {
                return err
        }
        // use x
        // 使用x
}
```

For example, `Foo`'s contract makes no guarantees about what it will return in the context of an error. The author of `Foo` is now free to annotate errors that pass through it with additional context without breaking its contract with the caller.

例如，`Foo`的契约不保证在错误上下文中它会返回什么。`Foo`的作者现在可以自由地为通过它的错误添加额外上下文的注释，而不会破坏与调用者的契约。

### Assert errors for behaviour, not type | 断言错误的行为，而不是类型

In a small number of cases, this binary approach to error handling is not sufficient.

在少数情况下，这种二进制的错误处理方法是不够的。

For example, interactions with the world outside your process, like network activity, require that the caller investigate the nature of the error to decide if it is reasonable to retry the operation.

例如，与进程外部世界的交互（如网络活动）需要调用者调查错误的性质，以决定重试操作是否合理。

In this case rather than asserting the error is a specific type or value, we can assert that the error implements a particular behaviour. Consider this example:

在这种情况下，我们可以断言错误实现了特定的行为，而不是断言错误是特定的类型或值。考虑这个例子：

```go
type temporary interface {
        Temporary() bool
}
 
// IsTemporary returns true if err is temporary.
// IsTemporary如果err是临时的则返回true。
func IsTemporary(err error) bool {
        te, ok := err.(temporary)
        return ok && te.Temporary()
}
```

We can pass any error to `IsTemporary` to determine if the error could be retried.

我们可以将任何错误传递给`IsTemporary`来确定该错误是否可以重试。

If the error does not implement the `temporary` interface; that is, it does not have a `Temporary` method, then then error is not temporary.

如果错误没有实现`temporary`接口，也就是说，它没有`Temporary`方法，那么这个错误就不是临时的。

If the error does implement `Temporary`, then perhaps the caller can retry the operation if `Temporary` returns `true`.

如果错误确实实现了`Temporary`，那么如果`Temporary`返回`true`，调用者可能可以重试操作。

The key here is this logic can be implemented without importing the package that defines the error or indeed knowing anything about `err`'s underlying type–we're simply interested in its behaviour.

这里的关键是这个逻辑可以在不导入定义错误的包的情况下实现，实际上不需要了解`err`的底层类型——我们只是对其行为感兴趣。

## Don't just check errors, handle them gracefully | 不要仅仅检查错误，要优雅地处理它们

This brings me to a second Go proverb that I want to talk about; don't just check errors, handle them gracefully. Can you suggest some problems with the following piece of code?

这让我想到了我要谈论的第二个Go谚语；不要仅仅检查错误，要优雅地处理它们。你能指出下面这段代码的一些问题吗？

```go
func AuthenticateRequest(r *Request) error {
        err := authenticate(r.User)
        if err != nil {
                return err
        }
        return nil
}
```

An obvious suggestion is that the five lines of the function could be replaced with

一个明显的建议是函数的五行可以替换为

```go
return authenticate(r.User)
```

But this is the simple stuff that everyone should be catching in code review. More fundamentally the problem with this code is I cannot tell where the original error came from.

但这是每个人都应该在代码审查中发现的简单问题。更根本的问题是我无法知道原始错误来自哪里。

If `authenticate` returns an error, then `AuthenticateRequest` will return the error to its caller, who will probably do the same, and so on. At the top of the program the main body of the program will print the error to the screen or a log file, and all that will be printed is: `No such file or directory`.

如果`authenticate`返回错误，那么`AuthenticateRequest`会将错误返回给调用者，调用者可能会做同样的事情，如此反复。在程序的顶部，程序的主体会将错误打印到屏幕或日志文件，所有打印的就是：`No such file or directory`。

There is no information of file and line where the error was generated. There is no stack trace of the call stack leading up to the error. The author of this code will be forced to a long session of bisecting their code to discover which code path trigged the file not found error.

没有生成错误的文件和行的信息。没有导致错误的调用栈的堆栈跟踪。这段代码的作者将被迫进行长时间的代码二分查找，以发现哪个代码路径触发了文件未找到错误。

Donovan and Kernighan's *The Go Programming Language* recommends that you add context to the error path using `fmt.Errorf`

Donovan和Kernighan的《Go程序设计语言》建议你使用`fmt.Errorf`向错误路径添加上下文

```go
func AuthenticateRequest(r *Request) error {
        err := authenticate(r.User)
        if err != nil {
                return fmt.Errorf("authenticate failed: %v", err)
        }
        return nil
}
```

But as we saw earlier, this pattern is incompatible with the use of sentinel error values or type assertions, because converting the error value to a string, merging it with another string, then converting it back to an error with `fmt.Errorf` breaks equality and destroys any context in the original error.

但正如我们之前看到的，这种模式与使用哨兵错误值或类型断言不兼容，因为将错误值转换为字符串，将其与另一个字符串合并，然后用`fmt.Errorf`将其转换回错误，这会破坏等值性并销毁原始错误中的任何上下文。

### Annotating errors | 注释错误

I'd like to suggest a method to add context to errors, and to do that I'm going to introduce a simple package. The code is online at [`github.com/pkg/errors`](https://godoc.org/github.com/pkg/errors). The errors package has two main functions:

我想建议一种向错误添加上下文的方法，为了做到这一点，我将介绍一个简单的包。代码在[`github.com/pkg/errors`](https://godoc.org/github.com/pkg/errors)在线可用。errors包有两个主要函数：

```go
// Wrap annotates cause with a message.
// Wrap用消息注释原因。
func Wrap(cause error, message string) error
```

The first function is `Wrap`, which takes an error, and a message and produces a new error.

第一个函数是`Wrap`，它接受一个错误和一个消息，并产生一个新错误。

```go
// Cause unwraps an annotated error.
// Cause解包一个注释的错误。
func Cause(err error) error
```

The second function is `Cause`, which takes an error that has possibly been wrapped, and unwraps it to recover the original error.

第二个函数是`Cause`，它接受一个可能已被包装的错误，并将其解包以恢复原始错误。

Using these two functions, we can now annotate any error, and recover the underlying error if we need to inspect it. Consider this example of a function that reads the content of a file into memory.

使用这两个函数，我们现在可以注释任何错误，并在需要检查时恢复底层错误。考虑这个将文件内容读入内存的函数示例。

```go
func ReadFile(path string) ([]byte, error) {
        f, err := os.Open(path)
        if err != nil {
                return nil, errors.Wrap(err, "open failed")
        } 
        defer f.Close()
 
        buf, err := ioutil.ReadAll(f)
        if err != nil {
                return nil, errors.Wrap(err, "read failed")
        }
        return buf, nil
}
```

We'll use this function to write a function to read a config file, then call that from `main`.

我们将使用这个函数来编写一个读取配置文件的函数，然后从`main`中调用它。

```go
func ReadConfig() ([]byte, error) {
        home := os.Getenv("HOME")
        config, err := ReadFile(filepath.Join(home, ".settings.xml"))
        return config, errors.Wrap(err, "could not read config")
}
 
func main() {
        _, err := ReadConfig()
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}
```

If the `ReadConfig` code path fails, because we used `errors.Wrap`, we get a nicely annotated error in the K&D style.

如果`ReadConfig`代码路径失败，因为我们使用了`errors.Wrap`，我们得到了K&D风格的良好注释错误。

```
could not read config: open failed: open /Users/dfc/.settings.xml: no such file or directory
```

Because `errors.Wrap` produces a stack of errors, we can inspect that stack for additional debugging information. This is the same example again, but this time we replace `fmt.Println` with `errors.Print`

因为`errors.Wrap`产生了错误堆栈，我们可以检查该堆栈以获取额外的调试信息。这又是同一个例子，但这次我们用`errors.Print`替换`fmt.Println`

```go
func main() {
        _, err := ReadConfig()
        if err != nil {
                errors.Print(err)
                os.Exit(1)
        }
}
```

We'll get something like this:

我们会得到类似这样的结果：

```
readfile.go:27: could not read config
readfile.go:14: open failed
open /Users/dfc/.settings.xml: no such file or directory
```

The first line comes from `ReadConfig`, the second comes from the `os.Open` part of `ReadFile`, and the remainder comes from the `os` package itself, which does not carry location information.

第一行来自`ReadConfig`，第二行来自`ReadFile`的`os.Open`部分，其余部分来自`os`包本身，它不携带位置信息。

Now we've introduced the concept of wrapping errors to produce a stack, we need to talk about the reverse, unwrapping them. This is the domain of the `errors.Cause` function.

现在我们已经介绍了包装错误以产生堆栈的概念，我们需要讨论相反的操作——解包它们。这是`errors.Cause`函数的领域。

```go
// IsTemporary returns true if err is temporary.
// IsTemporary如果err是临时的则返回true。
func IsTemporary(err error) bool {
        te, ok := errors.Cause(err).(temporary)
        return ok && te.Temporary()
}
```

In operation, whenever you need to check an error matches a specific value or type, you should first recover the original error using the `errors.Cause` function.

在操作中，每当你需要检查错误是否匹配特定值或类型时，你应该首先使用`errors.Cause`函数恢复原始错误。

## Only handle errors once | 只处理一次错误

Lastly, I want to mention that you should only handle errors once. Handling an error means inspecting the error value, and making a decision.

最后，我想提到你应该只处理一次错误。处理错误意味着检查错误值，并做出决定。

```go
func Write(w io.Writer, buf []byte) {
        w.Write(buf)
}
```

If you make less than one decision, you're ignoring the error. As we see here, the error from `w.Write` is being discarded.

如果你做出少于一个决定，你就是在忽略错误。正如我们在这里看到的，来自`w.Write`的错误被丢弃了。

But making more than one decision in response to a single error is also problematic.

但是对单个错误做出多于一个决定也是有问题的。

```go
func Write(w io.Writer, buf []byte) error {
        _, err := w.Write(buf)
        if err != nil {
                // annotated error goes to log file
                // 注释的错误记录到日志文件
                log.Println("unable to write:", err)
 
                // unannotated error returned to caller
                // 未注释的错误返回给调用者
                return err
        }
        return nil
}
```

In this example if an error occurs during `Write`, a line will be written to a log file, noting the file and line that the error occurred, and the error is also returned to the caller, who possibly will log it, and return it, all the way back up to the top of the program.

在这个例子中，如果在`Write`期间发生错误，将向日志文件写入一行，记录发生错误的文件和行，错误也会返回给调用者，调用者可能会记录它，并返回它，一直回到程序的顶部。

So you get a stack of duplicate lines in your log file, but at the top of the program you get the original error without any context. Java anyone?

所以你的日志文件中会有一堆重复的行，但在程序的顶部你得到的是没有任何上下文的原始错误。有人说Java吗？

```go
func Write(w io.Write, buf []byte) error {
        _, err := w.Write(buf)
        return errors.Wrap(err, "write failed")
}
```

Using the `errors` package gives you the ability to add context to error values, in a way that is inspectable by both a human and a machine.

使用`errors`包让你能够向错误值添加上下文，这种方式既可以被人类也可以被机器检查。

## Conclusion | 结论

In conclusion, errors are part of your package's public API, treat them with as much care as you would any other part of your public API.

总之，错误是你包的公共API的一部分，要像对待公共API的任何其他部分一样谨慎地对待它们。

For maximum flexibility I recommend that you try to treat all errors as opaque. In the situations where you cannot do that, assert errors for behaviour, not type or value.

为了获得最大的灵活性，我建议你尝试将所有错误视为不透明的。在无法做到这一点的情况下，断言错误的行为，而不是类型或值。

Minimise the number of sentinel error values in your program and convert errors to opaque errors by wrapping them with `errors.Wrap` as soon as they occur.

最小化程序中哨兵错误值的数量，并通过在错误发生时立即用`errors.Wrap`包装它们来将错误转换为不透明错误。

Finally, use `errors.Cause` to recover the underlying error if you need to inspect it.

最后，如果你需要检查底层错误，使用`errors.Cause`来恢复它。

---

**Related posts | 相关文章：**
- [Constant errors | 常量错误](https://dave.cheney.net/2016/04/07/constant-errors)
- [Stack traces and the errors package | 堆栈跟踪和errors包](https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package)
- [Inspecting errors | 检查错误](https://dave.cheney.net/2014/12/24/inspecting-errors)
- [Errors and Exceptions, redux | 错误和异常，重现](https://dave.cheney.net/2015/01/26/errors-and-exceptions-redux)