package fmte

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"sync"
)

var p *message.Printer

var mxStdOut, mxStdErr sync.Mutex

func init() {
	p = message.NewPrinter(language.English)
}

// Printf is goroutine-safe fmt.Printf for English
func Printf(format string, a ...interface{}) {
	mxStdOut.Lock()
	_, _ = p.Printf(format, a...)
	mxStdOut.Unlock()
}

// PrintfErr is goroutine-safe fmt.Printf to StdErr for English
func PrintfErr(format string, a ...interface{}) {
	mxStdErr.Lock()
	_, _ = p.Fprintf(os.Stderr, format, a...)
	mxStdErr.Unlock()
}
