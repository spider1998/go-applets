package main

import (
	"sdkeji/wechat/pkg/engine"
)

func main() {
	engine.NewStdInstance().Boot().Run()
}
