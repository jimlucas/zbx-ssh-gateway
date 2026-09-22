package config
import "testing"
func TestValidateRequiresPasswords(t *testing.T){c:=Config{Server:ServerConfig{Port:9443},SSH:SSHConfig{Username:"admin"},Limits:LimitsConfig{MaxWorkers:1,MaxQueue:1},Operations:map[string]Operation{"x":{Command:"show"}}};if c.Validate()==nil{t.Fatal("expected error")}}
