// Code generated by Thrift Compiler (0.15.0). DO NOT EDIT.

package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
)

var _ = nrprotocol.GoUnusedProtection__

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  void connect(JMXConfig config)")
  fmt.Fprintln(os.Stderr, "  void disconnect()")
  fmt.Fprintln(os.Stderr, "  string getClientVersion()")
  fmt.Fprintln(os.Stderr, "   queryMBeanNames(string mBeanNamePattern)")
  fmt.Fprintln(os.Stderr, "   getMBeanAttributeNames(string mBeanName)")
  fmt.Fprintln(os.Stderr, "   getMBeanAttributes(string mBeanName,  attributes)")
  fmt.Fprintln(os.Stderr, "   queryMBeanAttributes(string mBeanNamePattern)")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

type httpHeaders map[string]string

func (h httpHeaders) String() string {
  var m map[string]string = h
  return fmt.Sprintf("%s", m)
}

func (h httpHeaders) Set(value string) error {
  parts := strings.Split(value, ": ")
  if len(parts) != 2 {
    return fmt.Errorf("header should be of format 'Key: Value'")
  }
  h[parts[0]] = parts[1]
  return nil
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  headers := make(httpHeaders)
  var parsedUrl *url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Var(headers, "H", "Headers to set on the http(s) request (e.g. -H \"Key: Value\")")
  flag.Parse()
  
  if len(urlString) > 0 {
    var err error
    parsedUrl, err = url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  var cfg *thrift.TConfiguration = nil
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
    if len(headers) > 0 {
      httptrans := trans.(*thrift.THttpClient)
      for key, value := range headers {
        httptrans.SetHeader(key, value)
      }
    }
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans = thrift.NewTSocketConf(net.JoinHostPort(host, portStr), cfg)
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransportConf(trans, cfg)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactoryConf(cfg)
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactoryConf(cfg)
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryConf(cfg)
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  iprot := protocolFactory.GetProtocol(trans)
  oprot := protocolFactory.GetProtocol(trans)
  client := nrprotocol.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot))
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "connect":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "Connect requires 1 args")
      flag.Usage()
    }
    arg28 := flag.Arg(1)
    mbTrans29 := thrift.NewTMemoryBufferLen(len(arg28))
    defer mbTrans29.Close()
    _, err30 := mbTrans29.WriteString(arg28)
    if err30 != nil {
      Usage()
      return
    }
    factory31 := thrift.NewTJSONProtocolFactory()
    jsProt32 := factory31.GetProtocol(mbTrans29)
    argvalue0 := nrprotocol.NewJMXConfig()
    err33 := argvalue0.Read(context.Background(), jsProt32)
    if err33 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.Connect(context.Background(), value0))
    fmt.Print("\n")
    break
  case "disconnect":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "Disconnect requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.Disconnect(context.Background()))
    fmt.Print("\n")
    break
  case "getClientVersion":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetClientVersion requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetClientVersion(context.Background()))
    fmt.Print("\n")
    break
  case "queryMBeanNames":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "QueryMBeanNames requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.QueryMBeanNames(context.Background(), value0))
    fmt.Print("\n")
    break
  case "getMBeanAttributeNames":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetMBeanAttributeNames requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.GetMBeanAttributeNames(context.Background(), value0))
    fmt.Print("\n")
    break
  case "getMBeanAttributes":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetMBeanAttributes requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    arg37 := flag.Arg(2)
    mbTrans38 := thrift.NewTMemoryBufferLen(len(arg37))
    defer mbTrans38.Close()
    _, err39 := mbTrans38.WriteString(arg37)
    if err39 != nil { 
      Usage()
      return
    }
    factory40 := thrift.NewTJSONProtocolFactory()
    jsProt41 := factory40.GetProtocol(mbTrans38)
    containerStruct1 := nrprotocol.NewJMXServiceGetMBeanAttributesArgs()
    err42 := containerStruct1.ReadField2(context.Background(), jsProt41)
    if err42 != nil {
      Usage()
      return
    }
    argvalue1 := containerStruct1.Attributes
    value1 := argvalue1
    fmt.Print(client.GetMBeanAttributes(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "queryMBeanAttributes":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "QueryMBeanAttributes requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.QueryMBeanAttributes(context.Background(), value0))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
