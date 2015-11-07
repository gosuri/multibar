# MultiBar

MultiBar is Go library to render multiple progress bars in the terminal. 

This is a fork of [sethgrid/multibar](https://github.com/sethgrid/multibar) with API and documentation improvements.

![example](docs/example.gif)

## Example

The below example renders the progress bar like in the intro. 

```go
bars := multibar.New()
count1, count2 := 2000, 2500
bar1 := bars.MakeBar(count1, "bar1")
bar2 := bars.MakeBar(count2, "bar2")

bars.Start()

var wg sync.WaitGroup
wg.Add(1)
go func() {
  defer wg.Done()
  for i := 0; i <= count1; i++ {
    bars.Increment(bar1, i)
    time.Sleep(time.Millisecond)
  }
}()

wg.Add(1)
go func() {
  defer wg.Done()
  for i := 0; i <= count2; i++ {
    bars.Increment(bar2, i)
    time.Sleep(time.Millisecond)
  }
}()

wg.Wait()
```

Full source for the below example is in [example/main.go](example/main.go). To run the example, get the source code and run:

```
$ go run example/main.go
```

## Installation

Install the package using `go get`

```
go get -u github.com/gosuri/multibar
```
