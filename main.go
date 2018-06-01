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

const (
  VERSION = "v0.1.2"
  USAGE   = "notifyme (%s)\nUsage: %s <command>\n"
)

var (
  SLACK_CHANNEL string = ""
  SLACK_TOKEN   string = ""
)

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
  SLACK_CHANNEL = lookup_env("NOTIFYME_SLACK_CHANNEL", "alerts")
  SLACK_TOKEN   = lookup_env("NOTIFYME_SLACK_TOKEN", "")
}

func main() {
  command := strings.Join(os.Args[1:], " ")

  if len(command) > 0 {
    color := "good"

    fmt.Printf("==> Run notifyme...\n")
    fmt.Printf("--> Command: %s\n", command)
    start := current_timestamp()
    fmt.Printf("--> Start at: %s\n", start)
    stdout, exitcode := exec_command(command)
    stdout = clear_stdout(stdout)
    fmt.Printf("--> Stdout: %s", stdout)
    fmt.Printf("--> Exit code: %d\n", exitcode)
    end  := current_timestamp()
    diff := duration(start, end)
    fmt.Printf("--> End at: %s\n", end)
    fmt.Printf("--> Duration: %d\n", diff)

    if exitcode != 0 {
      color = "danger"
    }

    msg := &Message{
      Text: fmt.Sprintf("*From*: %s (%s)\nFinish executing the command on the server", hostname(), ip_address()),
      Channel: SLACK_CHANNEL,
    }
    msg.AddAttachment(&Attachment{
      Color: color,
      Text: fmt.Sprintf("*Command:* `%s`\n*Start at:* %s\n*End at:* %s\n*Duration:* %d seconds\n*Exit code:* %d", command, start, end, diff, exitcode),
    })

    slack_hook(msg)
  } else {
    fmt.Printf(USAGE, VERSION, os.Args[0])
  }
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

func current_timestamp() string {
  t := time.Now().UTC()
  return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func duration(start string, end string) int {
  parsed_start , _ := time.Parse("2006-01-02 15:04:05", start);
  parsed_end   , _ := time.Parse("2006-01-02 15:04:05", end);

  return int(parsed_end.Sub(parsed_start).Seconds())
}

func clear_stdout(stdout string) string {
  if strings.HasPrefix(stdout, "\n") == false {
    stdout = "\n" + stdout
  }
  if strings.HasSuffix(stdout, "\n") == false {
    stdout = stdout + "\n"
  }
  return stdout
}

func lookup_env(key string, default_value string) string {
  val, ok := os.LookupEnv(key)
  if !ok {
    if len(default_value) == 0 {
      fmt.Printf("Environment variable not defined: %s\n", key)
      os.Exit(1)
    }
    return default_value
  }

  return val
}
