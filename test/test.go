package main

import (
	"time"
	"tlog"
)

var logT1 = tlog.GetLogger("test1")
var logT2 = tlog.GetLogger("test2")

func main() {

	tlog.SetOption(tlog.DEBUG, "./traceid.log")

	for i := 0; i < 10000; i++ {
		go func(i int) {
			ctx := tlog.NewTraceIdCtx("")
			logT1.DebugF(ctx, "%d", i)
			logT1.InfoF(ctx, "%d", i)
			logT1.ErrorF(ctx, "%d", i)
		}(i)
	}

	go func() {
		for i := 10000; true; i++ {
			ctx := tlog.NewTraceIdCtx("")
			logT2.DebugF(ctx, "%d", i)
			logT2.InfoF(ctx, "%d", i)
		}
	}()

	time.Sleep(time.Second * 60)
}
