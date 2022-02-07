package whois

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"fort.plus/fperror"
	"fort.plus/im"
	"fort.plus/repository"
)

var (
	carrier   im.Carrier
	whoisFile string
)

func Register(t im.Carrier, whoisFileName string) {
	carrier = t
	whoisFile = whoisFileName
	repository.Register("/whois .*", a)
}

var a = func(message repository.RegExComparator) {
	msg := im.Cast(message)
	re := regexp.MustCompile("/whois (.*)")
	match := re.FindStringSubmatch(msg.Text)
	response := parseFile(match[1])
	carrier.Send(msg.From, strings.Join(response, "\n"))
}

func parseFile(pattern string) []string {
	var result []string

	re := regexp.MustCompile(pattern)
	readFile, err := os.Open(whoisFile)
	if err != nil {
		err = fperror.Warning("can't open file", err)
		log.Println(err)
		return []string{}
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if re.Match([]byte(line)) {
			result = append(result, line)
		}
	}
	return result
}
