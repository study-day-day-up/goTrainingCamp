# goTrainingCamp

### errgroup 理解
（以下摘自 04期Go交流群 涂超Eric 的总结）
> errgroup 的实现依靠于结构体 Group ，它通过封装 sync.WaitGroup，提供了比原生 sync.WaitGroup 更强大的三个特性：
> > 1.子 goroutine 的错误传递，
> > 一般情况来说，程序里启动的子 goroutine 运行出错，除了记录日志以外，是不好拿到的，而 errgroup 通过闭包的方式把子 goroutine的错误带出来，同时它又使用 sync.Once 保证即便多个 goroutine 同时报大量的错误，它也只会拿第一个 goroutine 的错误，避免了因处理大量且重复的错误数据，而耽误了进程的退出时间。
> > 2.继承了 WaitGroup 的控制能力，
> > 3.context传递机制，
> > 它可以非常方便的处理相关联的 goroutine 同时退出。只需使用 errgroup 提供的 Go方法 开启子goroutine，当传入的任务函数执行返回 error 时，errgroup 会自动帮你调 cancel 函数，你只需要在你的主进程或者需要同时退出的子 goroutine 中获取 context.Done() 信号，就可以做到同时退出。
