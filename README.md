# nio在golang中实践(目前依然有很多问题,勿用)

## 简介
- 由于Go的服务器开发推荐模型是基于Goroutines,但Goroutines并非没有开销,一个Goroutines大概占用2k~8k的内存,
因此本库想实现一个类似libevent的EventLoop,linux下使用epoll,bsd使用kqueue,其他平台使用select,API上参考java的nio

## 使用限制
- go版本需要1.9以上
- 目前并不支持Accept和Connect,因为go封装到了内部,无法设置NonBlock,除非完全自己实现net.Listener和net.Conn接口
- 由于使用EdgeTriggered模式,读写都要完全处理完
- OP_WRITE在需要的时候注册，不需要的时候取消

## 参考
- https://github.com/mailru/easygo
- https://github.com/tidwall/evio
- https://github.com/npat-efault/poller
- https://github.com/creack/goselect