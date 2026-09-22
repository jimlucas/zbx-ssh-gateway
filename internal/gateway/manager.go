package gateway
import("context";"errors";"sync";"time")
var ErrQueueFull=errors.New("request queue full")
type Request struct{ID,Target,Operation string;Parameters map[string]string}
type Result struct{Value,Stdout,Stderr string;ExitCode int;Err error;Duration time.Duration}
type Executor interface{Execute(context.Context,Request)(Result,error)}
type job struct{ctx context.Context;req Request;ch chan Result}
type Manager struct{q chan job;ex Executor;mu sync.RWMutex;active map[string]Request}
func New(workers,queue int,ex Executor)*Manager{m:=&Manager{q:make(chan job,queue),ex:ex,active:map[string]Request{}};for i:=0;i<workers;i++{go m.worker()};return m}
func(m *Manager)Submit(ctx context.Context,r Request)(Result,error){j:=job{ctx:ctx,req:r,ch:make(chan Result,1)};m.mu.Lock();m.active[r.ID]=r;m.mu.Unlock();select{case m.q<-j:case<-ctx.Done():m.del(r.ID);return Result{},ctx.Err();default:m.del(r.ID);return Result{},ErrQueueFull};select{case x:=<-j.ch:return x,x.Err;case<-ctx.Done():return Result{},ctx.Err()}}
func(m *Manager)worker(){for j:=range m.q{start:=time.Now();res,err:=m.ex.Execute(j.ctx,j.req);res.Duration=time.Since(start);res.Err=err;m.del(j.req.ID);j.ch<-res}}
func(m *Manager)del(id string){m.mu.Lock();defer m.mu.Unlock();delete(m.active,id)}
func(m *Manager)ActiveCount()int{m.mu.RLock();defer m.mu.RUnlock();return len(m.active)}
