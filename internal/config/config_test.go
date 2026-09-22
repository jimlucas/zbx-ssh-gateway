package config

import (
 "os"
 "path/filepath"
 "testing"
)

func TestValidateRequiresPasswords(t *testing.T){c:=Config{Server:ServerConfig{Port:9443},SSH:SSHConfig{Username:"admin"},Limits:LimitsConfig{MaxWorkers:1,MaxQueue:1},Operations:map[string]Operation{"x":{Command:"show"}}};if c.Validate()==nil{t.Fatal("expected error")}}

func TestLoadSeparateOperationsFile(t *testing.T){
 d:=t.TempDir()
 ops:=filepath.Join(d,"operations.local.yaml")
 if err:=os.WriteFile(ops,[]byte("operations:\n  device.version:\n    command: \"show version\"\n    parameters: {}\n"),0600);err!=nil{t.Fatal(err)}
 cfg:=filepath.Join(d,"gateway.local.yaml")
 body:="server:\n  port: 9443\nssh:\n  username: admin\n  port: 22\n  connect_timeout_seconds: 5\n  command_timeout_seconds: 10\n  password_list:\n    name: test\n    passwords: [one]\nlimits:\n  max_workers: 1\n  max_queue: 1\noperations_file: operations.local.yaml\n"
 if err:=os.WriteFile(cfg,[]byte(body),0600);err!=nil{t.Fatal(err)}
 c,err:=Load(cfg);if err!=nil{t.Fatal(err)}
 if _,ok:=c.Operations["device.version"];!ok{t.Fatal("operation not loaded")}
}
