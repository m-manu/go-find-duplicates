package fmte

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"sync"
)

var p *message.Printer

var mx sync.Mutex

func init() {
	p = message.NewPrinter(language.English)
}

var normalPrint = true

func Off() {
	normalPrint = false
}

// Printf is goroutine-safe fmt.Printf for English
func Printf(format string, a ...any) {
	if !normalPrint {
		return
	}
	mx.Lock()
	_, _ = p.Printf(format, a...)
	mx.Unlock()
}

// PrintfErr is goroutine-safe fmt.Printf to StdErr for English
func PrintfErr(format string, a ...any) {
	if !normalPrint {
		return
	}
	mx.Lock()
	_, _ = p.Fprintf(os.Stderr, format, a...)
	mx.Unlock()
}
