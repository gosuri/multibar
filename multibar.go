// Package multibar is library to render multiple progress bars in the terminal
package multibar

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gosuri/multibar/util/curse"
)

// MultiBar represnts the container holding progress bars
type MultiBar struct {
	// Bars represent the collection of progress bars in the container
	Bars []*ProgressBar

	screenLines   int
	screenWidth   int
	startingLine  int
	totalNewlines int

	history map[int]string
	sync.Mutex
	progFuncs map[*ProgressBar]progressFunc
}

// progressFunc is function used to increment progress bars
type progressFunc func(progress int)

// NewContainer returns a new instance of a container
func New() *MultiBar {
	width, lines, pos := getDimensions()
	return &MultiBar{
		screenWidth:  width,
		screenLines:  lines,
		startingLine: pos,
		history:      make(map[int]string),
		progFuncs:    make(map[*ProgressBar]progressFunc),
	}
}

// Start starts listening to updates without block
func (b *MultiBar) Start() {
	go b.Listen()
}

// Listen blocks the runtime and listens for updates on the progress bars
func (b *MultiBar) Listen() {
	for len(b.Bars) == 0 {
		// wait until we have some bars to work with
		time.Sleep(time.Millisecond * 100)
	}
	cases := make([]reflect.SelectCase, len(b.Bars))
	for i, bar := range b.Bars {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(bar.progressChan)}
	}

	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining -= 1
			continue
		}

		b.Bars[chosen].Increment(int(value.Int()))
	}
	b.Println()
}

// Increment increments the specified progress bar
func (b *MultiBar) Increment(bar *ProgressBar, n int) {
	b.progFuncs[bar](n)
}

// MakeBar makes and returns a progress. It registers the progress bar in the container
func (b *MultiBar) MakeBar(total int, prepend string) *ProgressBar {
	ch := make(chan int)
	bar := &ProgressBar{
		Width:           b.screenWidth - len(prepend) - 20,
		Total:           total,
		Prepend:         prepend,
		LeftEnd:         '[',
		RightEnd:        ']',
		Fill:            '=',
		Head:            '>',
		Empty:           '-',
		ShowPercent:     true,
		ShowTimeElapsed: true,
		StartTime:       time.Now(),
		progressChan:    ch,
	}

	b.Bars = append(b.Bars, bar)
	bar.Line = b.startingLine + b.totalNewlines
	b.history[bar.Line] = ""
	bar.Increment(0)
	b.Println()
	b.progFuncs[bar] = func(progress int) { bar.progressChan <- progress }
	return bar
}

// Print wrappers to capture newlines to adjust line positions on bars
func (b *MultiBar) Print(a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := countAllNewlines(a...)
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprint(a...)
	return fmt.Print(a...)
}

func (b *MultiBar) Printf(format string, a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := strings.Count(format, "\n")
	newlines += countAllNewlines(a...)
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprintf(format, a...)
	return fmt.Printf(format, a...)
}

func (b *MultiBar) Println(a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := countAllNewlines(a...) + 1
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprint(a...)
	return fmt.Println(a...)
}

func (b *MultiBar) addedNewlines(count int) {
	b.totalNewlines += count

	// if we hit the bottom of the screen, we "scroll" our bar displays by pushing
	// them up count lines (closer to line 0)
	if b.startingLine+b.totalNewlines > b.screenLines {
		b.totalNewlines -= count
		for _, bar := range b.Bars {
			bar.Line -= count
		}
		b.redrawAll(count)
	}
}

func (b *MultiBar) redrawAll(moveUp int) {
	c, _ := curse.New()

	newHistory := make(map[int]string)
	for line, printed := range b.history {
		newHistory[line+moveUp] = printed
		c.Move(1, line)
		c.EraseCurrentLine()
		c.Move(1, line+moveUp)
		c.EraseCurrentLine()
		fmt.Print(printed)
	}
	b.history = newHistory
	c.Move(c.StartingPosition.X, c.StartingPosition.Y)
}
