package repository

import (
	"testing"
	"fmt"
	"regexp"
)

type TestData struct {
    Text string
}
func (d TestData) IsRegExEqual(pattern string) (bool, error) {
    fmt.Printf("a1 function, message is:%s\n",pattern)
    return regexp.MatchString(pattern, d.Text)
}


var a1  = func(d RegExComparator)  {
    fmt.Printf("function call with data:%s\n", d.(TestData).Text)
}

//TODO: do some real tests here
func TestCheckRepository(t *testing.T) {
    var param TestData
    param.Text = "test-11"

    Register("t.st", a1)
    Register("te[sS]t", a1)

    Call(param)
}
