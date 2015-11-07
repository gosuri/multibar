package multibar

import (
	"regexp"
	"strings"
	"time"

	"github.com/gosuri/multibar/util/curse"
)

func countAllNewlines(interfaces ...interface{}) int {
	count := 0
	for _, iface := range interfaces {
		switch s := iface.(type) {
		case string:
			count += strings.Count(s, "\n")
		}
	}
	return count
}

func prettyTime(t time.Duration) string {
	re, err := regexp.Compile(`(\d+).(\d+)(\w+)`)
	if err != nil {
		return err.Error()
	}
	parts := re.FindSubmatch([]byte(t.String()))
	if len(parts) != 4 {
		return "---"
	}
	return string(parts[1]) + string(parts[3])
}

func getDimensions() (width, lines, position int) {
	width, lines, _ = curse.GetScreenDimensions()
	_, position, _ = curse.GetCursorPosition()
	return
}
