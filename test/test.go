package main

import (
	"time"
	"tlog"
)

var logT1 *tlog.Logger
var logT2 *tlog.Logger

func init() {
	tlog.SetOption(tlog.DEBUG, "./traceid.log")
	//tlog.SetOption(tlog.DEBUG, "")

	logT1 = tlog.GetLogger("test1")
	logT2 = tlog.GetLogger("test2")
}

func main() {

	for i := 0; i < 10000; i++ {
		go func(i int) {
			ctx := tlog.NewTraceIdCtx()
			logT1.DebugF(ctx, "%d", i)
			logT1.InfoF(ctx, "%d", i)
			logT1.ErrorF(ctx, "%d", i)
		}(i)
	}

	go func() {
		for i := 10000; true; i++ {
			ctx := tlog.NewTraceIdCtx()
			logT2.DebugF(ctx, "%d", i)
			logT2.InfoF(ctx, "%d", i)
		}
	}()

	time.Sleep(time.Second * 60)
}
