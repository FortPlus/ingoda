package mbroker

import (
    "regexp"
    "fmt"
)

type SyslogMessage struct {
  Host string       `json:"host"`
  Severity string   `json:"facility"`
  Facility string   `json:"severity"`
  Timestamp string  `json:"timestamp"`
  Text string       `json:"msg"`
}

func (m *SyslogMessage) AsText() string {
    return fmt.Sprintf("#syslog message at %s\n%s\n%s", m.Timestamp, m.Host, m.Text)
}

type NatsMessage struct {
    Subject string
    Text string
}
func (m NatsMessage) IsRegExEqual(pattern string) (bool, error) {
    return regexp.MatchString(pattern, m.Text)
}
