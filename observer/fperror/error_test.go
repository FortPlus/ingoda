package fperror


import (
	"testing"
	"regexp"
)

const SAMPLE = "message sample"

func TestWarn(t *testing.T) {
     msg := Warning(SAMPLE, nil).Error()
     res, _ := regexp.MatchString("Warning.*"+SAMPLE+".*",msg)
     if !res {
     	t.Errorf("error message %s doesn't match pattern", msg)
     }
}

func TestAlarm(t *testing.T) {
     msg := Alarm(SAMPLE, nil).Error()
     res, _ := regexp.MatchString("Alarm.*"+SAMPLE+".*",msg)
     if !res {
     	t.Errorf("error message %s doesn't match pattern", msg)
     }
}

func TestCritical(t *testing.T) {
     msg := Critical(SAMPLE, nil).Error()
     res, _ := regexp.MatchString("Critical.*"+SAMPLE+".*",msg)
     if !res {
     	t.Errorf("error message %s doesn't match pattern", msg)
     }
}