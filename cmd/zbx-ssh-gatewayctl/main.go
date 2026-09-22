package main
import("flag";"fmt";"io";"net/http";"os")
func main(){url:=flag.String("url","http://127.0.0.1:9443/api/v1/health","health URL");flag.Parse();r,e:=http.Get(*url);if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)};defer r.Body.Close();b,_:=io.ReadAll(r.Body);fmt.Print(string(b));if r.StatusCode>=400{os.Exit(1)}}
