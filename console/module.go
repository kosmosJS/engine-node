package console

import (
	"log"

	"github.com/kosmosJS/engine"
	"github.com/kosmosJS/engine-node/require"
	_ "github.com/kosmosJS/engine-node/util"
)

type Console struct {
	runtime *engine.Runtime
	util    *engine.Object
	printer Printer
}

type Printer interface {
	Log(string)
	Warn(string)
	Error(string)
}

type PrinterFunc func(s string)

func (p PrinterFunc) Log(s string) { p(s) }

func (p PrinterFunc) Warn(s string) { p(s) }

func (p PrinterFunc) Error(s string) { p(s) }

var defaultPrinter Printer = PrinterFunc(func(s string) { log.Print(s) })

func (c *Console) log(p func(string)) func(engine.FunctionCall) engine.Value {
	return func(call engine.FunctionCall) engine.Value {
		if format, ok := engine.AssertFunction(c.util.Get("format")); ok {
			ret, err := format(c.util, call.Arguments...)
			if err != nil {
				panic(err)
			}

			p(ret.String())
		} else {
			panic(c.runtime.NewTypeError("util.format is not a function"))
		}

		return nil
	}
}

func Require(runtime *engine.Runtime, module *engine.Object) {
	requireWithPrinter(defaultPrinter)(runtime, module)
}

func RequireWithPrinter(printer Printer) require.ModuleLoader {
	return requireWithPrinter(printer)
}

func requireWithPrinter(printer Printer) require.ModuleLoader {
	return func(runtime *engine.Runtime, module *engine.Object) {
		c := &Console{
			runtime: runtime,
			printer: printer,
		}

		c.util = require.Require(runtime, "util").(*engine.Object)

		o := module.Get("exports").(*engine.Object)
		o.Set("log", c.log(c.printer.Log))
		o.Set("error", c.log(c.printer.Error))
		o.Set("warn", c.log(c.printer.Warn))
	}
}

func Enable(runtime *engine.Runtime) {
	runtime.Set("console", require.Require(runtime, "console"))
}

func init() {
	require.RegisterNativeModule("console", Require)
}
