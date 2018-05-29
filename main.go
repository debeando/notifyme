package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net"
  "net/http"
  "os"
  "os/exec"
  "strings"
  "syscall"
  "time"
)

var SLACK_TOKEN   string = ""
var SLACK_CHANNEL string = ""

type Message struct {
  Text        string        `json:"text"`
  Channel     string        `json:"channel,omitempty"`
  UserName    string        `json:"username,omitempty"`
  Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
  Color  string `json:"color,omitempty"`
  Title  string `json:"title,omitempty"`
  Text   string `json:"text,omitempty"`
  Footer string `json:"footer_icon,omitempty"`
}

func init() {
  SLACK_TOKEN   = os.Getenv("NOTIFYME_SLACK_TOKEN")
  SLACK_CHANNEL = os.Getenv("NOTIFYME_SLACK_CHANNEL")
}

func main() {
  loc, _  := time.LoadLocation("UTC")
  start   := time.Now().In(loc)
  command := strings.Join(os.Args[1:], " ")
  color   := "good"

  fmt.Printf("==> Run...\n")
  fmt.Printf("--> Command: %s\n", command)
  fmt.Printf("--> Start at: %s\n", start)
  stdout, exitcode := exec_command(command)
  fmt.Printf("--> Stdout: \n%s\n", stdout)
  fmt.Printf("--> Exit code: %d\n", exitcode)
  end := time.Now().In(loc)
  fmt.Printf("--> End at: %s\n", end)

  if exitcode != 0 {
    color = "danger"
  }

  msg := &Message{
    Text: fmt.Sprintf("%s(%s)\nFinish executing the command on the server", hostname(), ip_address()),
    Channel: SLACK_CHANNEL,
  }
  msg.AddAttachment(&Attachment{
    Color: color,
    Text: fmt.Sprintf("*Command:* %s\n*Start at:* %s\n*End at:* %s\n*Exit code:* %d", command, start, end, exitcode),
  })

  slack_hook(msg)
}

func exec_command(cmd string) (stdout string, exitcode int) {
  out, err := exec.Command("bash", "-c", cmd).Output()
  if err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      ws := exitError.Sys().(syscall.WaitStatus)
      exitcode = ws.ExitStatus()
    }
  }
  stdout = string(out[:])
  return
}

func (m *Message) AddAttachment(a *Attachment) {
  m.Attachments = append(m.Attachments, a)
}

func slack_hook(msg *Message) {
  jsonValues, _ := json.Marshal(msg)

  req, err := http.NewRequest(
    "POST",
    "https://hooks.slack.com/services/" + SLACK_TOKEN,
    bytes.NewReader(jsonValues),
  )

  if err != nil {
    fmt.Print(err)
  }

  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Print(err)
  }

  fmt.Printf("--> Slack POST status code: %d\n", resp.StatusCode)
  defer resp.Body.Close()
}

func hostname() string {
  host, err := os.Hostname()
  if err != nil {
    fmt.Print(err)
  }

  return host
}

func ip_address() string {
  addrs, err := net.LookupIP(hostname())
  if err != nil {
    fmt.Print(err)
  }

  for _, addr := range addrs {
    if ipv4 := addr.To4(); ipv4 != nil {
      return ipv4.String()
    }
  }
  return ""
}
