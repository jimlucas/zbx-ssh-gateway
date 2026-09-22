package config

import (
 "fmt"
 "os"
 "path/filepath"
 "time"

 "gopkg.in/yaml.v3"
)

type Config struct { Server ServerConfig `yaml:"server"`; API APIConfig `yaml:"api"`; SSH SSHConfig `yaml:"ssh"`; Limits LimitsConfig `yaml:"limits"`; OperationsFile string `yaml:"operations_file"`; Operations map[string]Operation `yaml:"-"` }
type ServerConfig struct { Listen string `yaml:"listen"`; Port int `yaml:"port"`; TLS TLSConfig `yaml:"tls"` }
type TLSConfig struct { Enabled bool `yaml:"enabled"`; Certificate string `yaml:"certificate"`; PrivateKey string `yaml:"private_key"` }
type APIConfig struct { BearerTokenFile string `yaml:"bearer_token_file"` }
type SSHConfig struct { Username string `yaml:"username"`; Port int `yaml:"port"`; ConnectTimeout time.Duration `yaml:"-"`; ConnectTimeoutSeconds int `yaml:"connect_timeout_seconds"`; CommandTimeout time.Duration `yaml:"-"`; CommandTimeoutSeconds int `yaml:"command_timeout_seconds"`; KnownHosts string `yaml:"known_hosts"`; PasswordList PasswordList `yaml:"password_list"` }
type PasswordList struct { Name string `yaml:"name"`; Passwords []string `yaml:"passwords"` }
type LimitsConfig struct { MaxWorkers int `yaml:"max_workers"`; MaxQueue int `yaml:"max_queue"`; MaxOutputBytes int64 `yaml:"max_output_bytes"`; MaxRequestBytes int64 `yaml:"max_request_bytes"` }
type Operation struct { Command string `yaml:"command"`; Parameters map[string]Parameter `yaml:"parameters"`; TimeoutSeconds int `yaml:"timeout_seconds"`; MaxOutputBytes int64 `yaml:"max_output_bytes"` }
type Parameter struct { Type string `yaml:"type"`; Values []string `yaml:"values"`; Required bool `yaml:"required"`; Pattern string `yaml:"pattern"` }
type operationsConfig struct { Operations map[string]Operation `yaml:"operations"` }

func Load(path string)(*Config,error){
 b,e:=os.ReadFile(path);if e!=nil{return nil,e}
 var c Config
 if e=yaml.Unmarshal(b,&c);e!=nil{return nil,e}
 if c.OperationsFile==""{return nil,fmt.Errorf("operations_file is required")}
 opath:=c.OperationsFile
 if !filepath.IsAbs(opath){opath=filepath.Join(filepath.Dir(path),opath)}
 ob,e:=os.ReadFile(opath);if e!=nil{return nil,fmt.Errorf("read operations file %q: %w",opath,e)}
 var oc operationsConfig
 if e=yaml.Unmarshal(ob,&oc);e!=nil{return nil,fmt.Errorf("parse operations file %q: %w",opath,e)}
 c.Operations=oc.Operations
 c.SSH.ConnectTimeout=time.Duration(c.SSH.ConnectTimeoutSeconds)*time.Second
 c.SSH.CommandTimeout=time.Duration(c.SSH.CommandTimeoutSeconds)*time.Second
 if e=c.Validate();e!=nil{return nil,e}
 return &c,nil
}
func(c *Config)Validate()error{if c.Server.Port<1||c.Server.Port>65535{return fmt.Errorf("invalid server port")};if c.SSH.Username==""{return fmt.Errorf("ssh username required")};if len(c.SSH.PasswordList.Passwords)==0{return fmt.Errorf("at least one SSH password is required")};if c.Limits.MaxWorkers<1||c.Limits.MaxQueue<1{return fmt.Errorf("worker and queue limits must be positive")};if len(c.Operations)==0{return fmt.Errorf("at least one operation is required")};for n,o:=range c.Operations{if n==""||o.Command==""{return fmt.Errorf("operation %q requires command",n)}};return nil}
