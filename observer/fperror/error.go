package fperror

import(
    "fmt"
    "runtime"
)


const (
    ALARM    = 0
    WARNING  = 1
    CRITICAL = 2
)
var levelName = []string{"Alarm", "Warning", "Critical"}

type customError struct {
    level int
    err error
    msg string
    trace string
}

func (e *customError) Error() string {
    if e.err != nil {
        return fmt.Sprintf("%s, %s, trace:%s, %s",levelName[e.level], e.msg, e.trace, e.err)
    } else { 
        return fmt.Sprintf("%s, %s, trace:%s",levelName[e.level], e.msg, e.trace)
    }
}

func Alarm(errorText string, err error) *customError {
    return &customError{level: ALARM, msg: errorText, err: err, trace: getTrace()}
}

func Warning(errorText string, err error) *customError {
    return &customError{level: WARNING, msg: errorText, err: err, trace: getTrace()}
}

func Critical(errorText string, err error) *customError {
    return &customError{level: CRITICAL, msg: errorText, err: err, trace: getTrace()}
}

func getTrace() string {
    stackSlice := make([]byte, 1024)
	s := runtime.Stack(stackSlice, false)
    return fmt.Sprintf("\n%s\n", stackSlice[0:s])
}

