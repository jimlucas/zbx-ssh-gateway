package operations
import("testing";"github.com/jimlucas/zbx-ssh-gateway/internal/config")
func TestBuildEnum(t *testing.T){o:=config.Operation{Command:"show ${iface}",Parameters:map[string]config.Parameter{"iface":{Type:"enum",Values:[]string{"wlan0"},Required:true}}};s,e:=Build(o,map[string]string{"iface":"wlan0"});if e!=nil||s!="show wlan0"{t.Fatalf("%q %v",s,e)};if _,e=Build(o,map[string]string{"iface":"wlan0; reboot"});e==nil{t.Fatal("expected rejection")}}
func TestUnknownParameter(t *testing.T){o:=config.Operation{Command:"show",Parameters:map[string]config.Parameter{}};if _,e:=Build(o,map[string]string{"x":"y"});e==nil{t.Fatal("expected rejection")}}
