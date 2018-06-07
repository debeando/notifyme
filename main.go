package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net"
  "net/http"
  "os"
  "os/signal"
  "os/exec"
  "strings"
  "syscall"
  "time"
)

const (
  VERSION = "v0.1.3"
  USAGE = "notifyme (%s)\nUsage: %s <command>\n"
  MESSAGE_TEXT = "*From*: %s (%s)\nFinish executing the command on the server"
  MESSAGE_ATTACHMENT_TEXT = "*Command:* `%s`\n*Start at:* %s\n*End at:* %s\n*Duration:* %d seconds\n*Exit code:* %d"
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
  wait_for_interrupt()

  command := get_command()

  if len(command) > 0 {
    color := "good"

    fmt.Printf("==> Run notifyme...\n")
    fmt.Printf("--> Press Ctrl+C to end.\n")
    start := current_timestamp()
    fmt.Printf("--> Start at: %s\n", start)
    fmt.Printf("--> Wait to finish command: %s\n", command)
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
      Text: fmt.Sprintf(MESSAGE_TEXT, hostname(), ip_address()),
      Channel: SLACK_CHANNEL,
    }
    msg.AddAttachment(&Attachment{
      Color: color,
      Text: fmt.Sprintf(MESSAGE_ATTACHMENT_TEXT, command, start, end, diff, exitcode),
    })

    response_code := slack_hook(msg)
    fmt.Printf("--> Slack response code: %d\n", response_code)
  } else {
    fmt.Printf(USAGE, VERSION, os.Args[0])
  }
}

func wait_for_interrupt() {
  shutdown_signals := make(chan os.Signal, 1)
  signal.Notify(shutdown_signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

  go func() {
    <-shutdown_signals
    fmt.Print("\r")
    fmt.Printf("--> Interrupted!.\n")
    os.Exit(1)
  }()
}

func get_command() string {
  return strings.Join(os.Args[1:], " ")
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

func slack_hook(msg *Message) int {
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
  defer resp.Body.Close()
  if err != nil {
    fmt.Print(err)
  }

  return resp.StatusCode
}

func hostname() string {
  host, err := os.Hostname()
  if err != nil {
    fmt.Print(err)
  }

  return host
}

func ip_address() string {
  addrs, _ := net.InterfaceAddrs()

  for _, a := range addrs {
    if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
      if ipnet.IP.To4() != nil {
        return ipnet.IP.String()
      }
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
