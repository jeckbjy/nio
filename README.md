# nio参考java nio的golang实现(目前依然有很多问题,误用)

## 简介
- 由于Go的服务器开发推荐模型是基于Goroutines,但Goroutines并非没有开销,一个Goroutines大概占用2k~8k的内存,
因此本库想实现一个类似libevent的EventLoop,linux下使用epoll,bsd使用kqueue,其他平台使用select,API上参考java的nio

## 范例
```$xslt
// 服务器范例
func StartServer() {
	log.Printf("start server")
	l, err := net.Listen("tcp", ":6789")
	if err != nil {
		panic(err)
	}

    // 创建Selector
	selector, err := nio.New()
	if err != nil {
		panic(err)
	}

    // 注册Listner
	if _, err := selector.Add(l, nio.OP_ACCEPT, nil); err != nil {
		panic(err)
	}

	for {
	    // 监听Select事件
		keys, err := selector.Select()
		if err != nil {
			break
		}
		
        // 处理事件
		for _, key := range keys {
			switch {
			case key.Acceptable():
				log.Printf("accept\n")
				conn, err := key.Listener().Accept()
				if err != nil {
					panic(err)
				}
				selector.Add(conn, nio.OP_READ, nil)
			case key.Readable():
				bytes := make([]byte, 1024)
				n, err := key.Conn().Read(bytes)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", bytes[:n])
				key.Conn().Write([]byte("pong"))
			}
		}
	}
}
```

## 使用限制
- 目前不支持Connect,由于Go将socket创建，Connect都封装到Dial函数中,并没有办法在Dial之前获取socket fd，注册到selector中
- Add注册传入的是interface{},需要能支持ifile接口，用于获取os.File，这是因为Go没有提供Socket的概念,而且提供的net.Listener和net.Conn接口
- 由于使用EdgeTriggered模式,读写都要完全处理完
- 写状态在需要的时候注册，不需要的时候取消

## 参考
- https://github.com/mailru/easygo
- https://github.com/tidwall/evio
- https://github.com/npat-efault/poller
- https://github.com/creack/goselect