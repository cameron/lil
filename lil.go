package main

import (
  "github.com/codegangsta/cli"
  "os"
  "net"
  "errors"
  "net/http"
  "fmt"
  "strings"
)

func main() {
  usage := "Usage: lil <name>"
  app := cli.NewApp()
  app.Name = "lil"
  app.Usage = usage
  app.Action = func(c *cli.Context){
    if(len(c.Args().First()) == 0){
      println(usage)
      return
    }

    url := fmt.Sprint("https://lil.firebaseio.com/", c.Args()[0], ".json")
    client := &http.Client{}
    ip, err := externalIP()
    if err != nil {
      println("Error getting external IP.", err)
      return
    }

    ip = "\"" + ip + "\""
    request, err := http.NewRequest("PUT", url, strings.NewReader(ip))
    request.ContentLength = int64(len(ip))
    _, err = client.Do(request)
    if err != nil {
      println("Error making request", err.Error())
    }
    return
  }
  app.Run(os.Args)
}

func externalIP() (string, error) {
  ifaces, err := net.Interfaces()
  if err != nil {
    return "", err
  }
  for _, iface := range ifaces {
    if iface.Flags&net.FlagUp == 0 {
    continue // interface down
    }
    if iface.Flags&net.FlagLoopback != 0 {
    continue // loopback interface
    }
    addrs, err := iface.Addrs()
    if err != nil {
      return "", err
    }
    for _, addr := range addrs {
    var ip net.IP
      switch v := addr.(type) {
      case *net.IPNet:
        ip = v.IP
      case *net.IPAddr:
        ip = v.IP
      }
      if ip == nil || ip.IsLoopback() {
      continue
      }
      ip = ip.To4()
      if ip == nil {
      continue // not an ipv4 address
      }
      return ip.String(), nil
    }
  }
  return "", errors.New("are you connected to the network?")
}
