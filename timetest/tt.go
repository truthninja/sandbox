package main
import (
	"fmt"
	"github.com/mailgun/godebug/lib"
	"time"
)
var tt_go_scope = godebug.EnteringNewScope(tt_go_contents)
func main() {
	ctx, ok := godebug.EnterFunc(main)
	if !ok {
		return
	}
	godebug.Line(ctx, tt_go_scope, 10)
	now := time.Now()
	scope := tt_go_scope.EnteringNewChildScope()
	scope.Declare("now", &now)
	godebug.Line(ctx, scope, 11)
	fmt.Printf("====NOw its %+v", now)
	godebug.Line(ctx, scope, 12)
	for time.Now() == now {
		godebug.Line(ctx, scope, 13)
		fmt.Printf("now %+v and time %+vi\n", now, time.Now())
		godebug.Line(ctx, scope, 12)
	}
	godebug.SetTraceGen(ctx)
	godebug.Line(ctx, scope, 16)
	fmt.Printf("Old time %v new time %v\n", now, time.Now())
	godebug.Line(ctx, scope, 17)
	sec := 1 * time.Second
	scope.Declare("sec", &sec)
	godebug.Line(ctx, scope, 18)
	milli := 1 * time.Millisecond
	scope.Declare("milli", &milli)
	godebug.Line(ctx, scope, 19)
	micro := 1 * time.Microsecond
	scope.Declare("micro", &micro)
	godebug.Line(ctx, scope, 20)
	nano := 1 * time.Nanosecond
	scope.Declare("nano", &nano)
	godebug.Line(ctx, scope, 22)
	fmt.Printf("Duration %+v\n", time.Since(now))
	godebug.Line(ctx, scope, 23)
	fmt.Printf("Duration int %+v\n", int64(time.Since(now)))
	godebug.Line(ctx, scope, 25)
	fmt.Printf("milli %+v\n", milli)
	godebug.Line(ctx, scope, 26)
	fmt.Printf("milli int %+v\n", int64(milli))
	godebug.Line(ctx, scope, 27)
	fmt.Printf("micro %+v\n", micro)
	godebug.Line(ctx, scope, 28)
	fmt.Printf("micro int %+v\n", int64(micro))
	godebug.Line(ctx, scope, 29)
	fmt.Printf("sec %+v\n", sec)
	godebug.Line(ctx, scope, 30)
	fmt.Printf("sec int %+v\n", int64(sec))
	godebug.Line(ctx, scope, 31)
	fmt.Printf("nano %+v\n", nano)
	godebug.Line(ctx, scope, 32)
	fmt.Printf("nano int %+v\n", int64(nano))
}

var tt_go_contents = `package main

import (
	"fmt"
	"github.com/mailgun/godebug/lib"
	"time"
)

func main() {
	now := time.Now()
	fmt.Printf("====NOw its %+v", now)
	for time.Now() == now {
		fmt.Printf("now %+v and time %+vi\n", now, time.Now())
	}
	godebug.SetTrace()
	fmt.Printf("Old time %v new time %v\n", now, time.Now())
	sec := 1 * time.Second
	milli := 1 * time.Millisecond
	micro := 1 * time.Microsecond
	nano := 1 * time.Nanosecond

	fmt.Printf("Duration %+v\n", time.Since(now))
	fmt.Printf("Duration int %+v\n", int64(time.Since(now)))

	fmt.Printf("milli %+v\n", milli)
	fmt.Printf("milli int %+v\n", int64(milli))
	fmt.Printf("micro %+v\n", micro)
	fmt.Printf("micro int %+v\n", int64(micro))
	fmt.Printf("sec %+v\n", sec)
	fmt.Printf("sec int %+v\n", int64(sec))
	fmt.Printf("nano %+v\n", nano)
	fmt.Printf("nano int %+v\n", int64(nano))
}
`
