package mbroker

import ("regexp")

type SyslogMessage struct {
    Subject string
    Text string
}
func (m SyslogMessage) IsRegExEqual(pattern string) (bool, error) {
    return regexp.MatchString(pattern, m.Text)
}
