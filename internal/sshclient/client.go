package sshclient
import("bytes";"context";"fmt";"net";"strconv";"time";"golang.org/x/crypto/ssh";"golang.org/x/crypto/ssh/knownhosts")
type Config struct{Username string;Port int;Passwords []string;ConnectTimeout,CommandTimeout time.Duration;KnownHosts string;MaxOutputBytes int64}
type Result struct{Stdout,Stderr string;ExitCode int;Attempts int}
type Client struct{cfg Config}
func New(c Config)*Client{return &Client{cfg:c}}
func(c *Client)Run(ctx context.Context,target,command string)(Result,error){var last error;for i,pw:=range c.cfg.Passwords{r,e:=c.runOnce(ctx,target,command,pw);r.Attempts=i+1;if e==nil{return r,nil};last=e;if !isAuth(e){return r,e}};return Result{Attempts:len(c.cfg.Passwords)},fmt.Errorf("authentication failed after %d attempts: %w",len(c.cfg.Passwords),last)}
func(c *Client)runOnce(ctx context.Context,target,command,password string)(Result,error){cb,e:=knownhosts.New(c.cfg.KnownHosts);if e!=nil{return Result{},e};cfg:=&ssh.ClientConfig{User:c.cfg.Username,Auth:[]ssh.AuthMethod{ssh.Password(password)},HostKeyCallback:cb,Timeout:c.cfg.ConnectTimeout};addr:=net.JoinHostPort(target,strconv.Itoa(c.cfg.Port));cl,e:=ssh.Dial("tcp",addr,cfg);if e!=nil{return Result{},e};defer cl.Close();s,e:=cl.NewSession();if e!=nil{return Result{},e};defer s.Close();var out,errb limitedBuffer;out.max=c.cfg.MaxOutputBytes;errb.max=c.cfg.MaxOutputBytes;s.Stdout=&out;s.Stderr=&errb;done:=make(chan error,1);go func(){done<-s.Run(command)}();select{case e:=<-done:r:=Result{Stdout:out.String(),Stderr:errb.String()};if ee,ok:=e.(*ssh.ExitError);ok{r.ExitCode=ee.ExitStatus();return r,nil};return r,e;case<-ctx.Done():return Result{},ctx.Err();case<-time.After(c.cfg.CommandTimeout):return Result{},fmt.Errorf("command timeout")}}
func isAuth(e error)bool{return e!=nil&&bytes.Contains([]byte(e.Error()),[]byte("unable to authenticate"))}
type limitedBuffer struct{b bytes.Buffer;max int64}
func(l *limitedBuffer)Write(p []byte)(int,error){if int64(l.b.Len()+len(p))>l.max{return 0,fmt.Errorf("output limit exceeded")};return l.b.Write(p)}
func(l *limitedBuffer)String()string{return l.b.String()}
