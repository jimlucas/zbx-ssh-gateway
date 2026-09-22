package gateway
import("context";"testing")
type fake struct{}
func(fake)Execute(ctx context.Context,r Request)(Result,error){return Result{Value:"ok"},nil}
func TestManager(t *testing.T){m:=New(1,2,fake{});r,e:=m.Submit(context.Background(),Request{ID:"1"});if e!=nil||r.Value!="ok"{t.Fatalf("%v %#v",e,r)}}
