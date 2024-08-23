package main

import (
	pailog "v1"
)

var r1 = pailog.GetLogger("runlog")
var r2 = pailog.GetLogger("tracelog")
var r3 = pailog.GetLogger("tracelog2")

func ShowLog() {
	pailog.SetOption(pailog.TRACE, "./runlog")

	r1.Debugf("debug log")
	r1.Infof("info log")
	r1.Errorf("error log")
}

func ShowTrace() {

	pailog.SetOption(pailog.DEBUG, "./tracelog1")

	ctx123 := pailog.NewTraceIdCtx("123")
	r2.DebugF(ctx123, "debug log")
	r2.InfoF(ctx123, "info log")
	r2.ErrorF(ctx123, "error log")

	ctx456 := pailog.NewTraceIdCtx("456")
	r2.DebugF(ctx456, "debug log")
	r2.InfoF(ctx456, "info log")
	r2.ErrorF(ctx456, "error log")
}

func DoubleTrace() {

	pailog.SetOption(pailog.DEBUG, "./tracelog2")

	ctx123 := pailog.NewTraceIdCtx("123")
	r3.DebugF(ctx123, "debug log")
	r3.InfoF(ctx123, "info log")
	r3.ErrorF(ctx123, "error log")

	ctx123456 := pailog.UpdateTrace(ctx123, "456")
	r3.DebugF(ctx123456, "debug log")
	r3.InfoF(ctx123456, "info log")
	r3.ErrorF(ctx123456, "error log")

	ctx123456789 := pailog.UpdateTrace(ctx123456, "789")
	r3.DebugF(ctx123456789, "debug log")
	r3.InfoF(ctx123456789, "info log")
	r3.ErrorF(ctx123456789, "error log")
}

func main() {
	ShowLog()
	//ShowTrace()
	//DoubleTrace()
}
