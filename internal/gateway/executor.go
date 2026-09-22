package gateway
import("context";"fmt";"github.com/jimlucas/zbx-ssh-gateway/internal/config";"github.com/jimlucas/zbx-ssh-gateway/internal/operations";"github.com/jimlucas/zbx-ssh-gateway/internal/sshclient")
type SSHExecutor struct{Ops map[string]config.Operation;Client *sshclient.Client}
func(e *SSHExecutor)Execute(ctx context.Context,r Request)(Result,error){op,ok:=e.Ops[r.Operation];if !ok{return Result{},fmt.Errorf("unknown operation")};cmd,err:=operations.Build(op,r.Parameters);if err!=nil{return Result{},err};x,err:=e.Client.Run(ctx,r.Target,cmd);return Result{Value:x.Stdout,Stdout:x.Stdout,Stderr:x.Stderr,ExitCode:x.ExitCode},err}
