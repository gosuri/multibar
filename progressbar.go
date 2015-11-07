package multibar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gosuri/multibar/util/curse"
)

// ProgressBar represents a progress bar
type ProgressBar struct {
	// Width is the width of the progress bar
	Width int

	// Total is the total value of the progress bar
	Total int

	// LeftEnd is character in the left most part of the progress indicator. Defaults to '['
	LeftEnd byte

	// RightEnd is character in the right most part of the progress indicator. Defaults to ']'
	RightEnd byte

	// Fill is the character representing completed progress. Defaults to '='
	Fill byte

	// Head is the character that moves when progress is updated.  Defaults to '>'
	Head byte

	// Empty is the character that represents the empty progress. Default is '-'
	Empty byte

	// ShowPercent is the flag used to determine if completion percentage should be displayed
	ShowPercent bool

	// ShowTimeElapsed is the flag used to determine is the elapsed time should be shown
	ShowTimeElapsed bool

	// StartTime is the time progress has started
	StartTime time.Time

	// Prepend is the string prepended before progress bar
	Prepend string

	// Line is the line in which the progress bar is rendered when multiple progress bars are present.
	// Defaults to current cursor position
	Line int

	progressChan chan int
}

// NewProgressBar returns a new progress bar
func NewProgressBar(total int) *ProgressBar {
	width, _, pos := getDimensions()
	return &ProgressBar{
		Width:           width - 20,
		Line:            pos,
		Total:           total,
		LeftEnd:         '[',
		RightEnd:        ']',
		Fill:            '=',
		Head:            '>',
		Empty:           '-',
		ShowPercent:     true,
		ShowTimeElapsed: true,
		StartTime:       time.Now(),
	}
}

// Increment increments the progress bar and renders
func (p *ProgressBar) Increment(n int) {
	bar := make([]string, p.Width)

	// avoid division by zero errors on non-properly constructed progressbars
	if p.Width == 0 {
		p.Width = 1
	}
	if p.Total == 0 {
		p.Total = 1
	}
	justGotToFirstEmptySpace := true
	for i, _ := range bar {
		if float32(n)/float32(p.Total) > float32(i)/float32(p.Width) {
			bar[i] = string(p.Fill)
		} else {
			bar[i] = string(p.Empty)
			if justGotToFirstEmptySpace {
				bar[i] = string(p.Head)
				justGotToFirstEmptySpace = false
			}
		}
	}

	percent := ""
	if p.ShowPercent {
		asInt := int(100 * (float32(n) / float32(p.Total)))
		padding := ""
		if asInt < 10 {
			padding = "  "
		} else if asInt < 99 {
			padding = " "
		}
		percent = padding + strconv.Itoa(asInt) + "% "
	}

	timeElapsed := ""
	if p.ShowTimeElapsed {
		timeElapsed = " " + prettyTime(time.Since(p.StartTime))
	}

	// record where we are, jump to the progress bar, update it, jump back
	c, _ := curse.New()
	c.Move(1, p.Line)
	c.EraseCurrentLine()
	fmt.Printf("\r%s %s%c%s%c%s", p.Prepend, percent, p.LeftEnd, strings.Join(bar, ""), p.RightEnd, timeElapsed)
	c.Move(c.StartingPosition.X, c.StartingPosition.Y)
}
