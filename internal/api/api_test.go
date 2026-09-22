package api
import("context";"net/http/httptest";"strings";"testing";"github.com/jimlucas/zbx-ssh-gateway/internal/gateway")
type fake struct{}
func(fake)Execute(ctx context.Context,r gateway.Request)(gateway.Result,error){return gateway.Result{Value:"42",Stdout:"42"},nil}
func TestUnauthorized(t *testing.T){s:=&Server{Manager:gateway.New(1,2,fake{}),Token:"secret",MaxBody:1024};r:=httptest.NewRequest("POST","/api/v1/execute",strings.NewReader(`{"target":"x","operation":"y"}`));w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=401{t.Fatalf("got %d",w.Code)}}
func TestUnknownJSONFieldRejected(t *testing.T){s:=&Server{Manager:gateway.New(1,2,fake{}),MaxBody:1024};r:=httptest.NewRequest("POST","/api/v1/execute",strings.NewReader(`{"target":"x","operation":"y","command":"bad"}`));w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=400{t.Fatalf("got %d",w.Code)}}
