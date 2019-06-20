package log

import (
	"fmt"
	l "log"
	"os"
	"runtime"
)

// Level defines severity of log entity
type Level int

const (
	info Level = iota + 1
	warn
	erro
)

func (l Level) String() string {
	switch l {
	case info:
		return "INFO"
	case warn:
		return "WARN"
	case erro:
		return "ERRO"
	}
	return ""
}

// OpenTag returns cli open tag
func (l Level) OpenTag() string {
	switch l {
	case info:
		return "\033[36m"
	case warn:
		return "\033[33m"
	case erro:
		return "\033[31m"
	}
	return "\033[36m"
}

// CloseTag returns cli close tag
func (l Level) CloseTag() string {
	switch l {
	case info:
		return "\033[00m"
	case warn:
		return "\033[00m"
	case erro:
		return "\033[00m"
	}
	return "\033[00m"
}

// Channel for listening logs
var Channel []chan *Entity

// NewChan addes new channel to the list of channels where logs will be released
func NewChan() chan *Entity {
	ch := make(chan *Entity)
	Channel = append(Channel, ch)
	return ch
}

// Entity single log
type Entity struct {
	Level   Level
	Path    string
	Line    int
	Message string
	Objects map[string]interface{}
	Err     error
	Context context
}

// New creates new instance of Entity
func New(f string, l int) *Entity {
	return &Entity{
		Path:    f,
		Line:    l,
		Objects: make(map[string]interface{}),
	}
}

// Info msg
func Info(msg string) {
	_, f, l, _ := runtime.Caller(1)
	New(f, l).Info(msg)
}

type context interface {
	GetCasino() string
}

// WithContext logs context
func WithContext(c context) *Entity {
	_, f, l, _ := runtime.Caller(1)
	return New(f, l).WithContext(c)
}

// Caller sets caller
func Caller(i int) *Entity {
	_, f, l, _ := runtime.Caller(i)
	return New(f, l)
}

// WithError receives error
func WithError(err error) *Entity {
	_, f, l, _ := runtime.Caller(1)
	return New(f, l).WithError(err)
}

// With key and object
func With(key string, o interface{}) *Entity {
	_, f, l, _ := runtime.Caller(1)
	return New(f, l).With(key, o)
}

// Info receives msg and sets it into entity
func (e *Entity) Info(msg string) {
	e.Message = msg
	e.Level = info
	e.print()
}

// Warn receives msg and sets it into entity
func (e *Entity) Warn(msg string) {
	e.Message = msg
	e.Level = warn
	e.print()
}

// Error receives msg and sets it into entity
func (e *Entity) Error(msg string) {
	e.Message = msg
	e.Level = erro
	e.print()
}

// WithError receives error
func (e *Entity) WithError(err error) *Entity {
	e.Err = err
	return e
}

// With key and object
func (e *Entity) With(key string, o interface{}) *Entity {
	e.Objects[key] = o
	return e
}

// WithContext stores ctx for better logging
func (e *Entity) WithContext(c context) *Entity {
	e.Context = c
	return e
}

// Do should be at the end of every log to make log write
func (e *Entity) print() {
	cli := l.New(os.Stdout, fmt.Sprintf("%s%s%s ", e.Level.OpenTag(), e.Level, e.Level.CloseTag()), l.LstdFlags)
	str := fmt.Sprintf("%v:%v \n\t%s>> message: %v%s", e.Path, e.Line, e.Level.OpenTag(), e.Message, e.Level.CloseTag())
	if e.Context != nil {
		e.With("Casino", e.Context.GetCasino())
	}
	if e.Err != nil {
		str += fmt.Sprintf("\n\t%s>> Error:%s %#v", e.Level.OpenTag(), e.Level.CloseTag(), e.Err)
	}
	if len(e.Objects) > 0 {
		for k, v := range e.Objects {
			str += fmt.Sprintf("\n\t%s>> %s:%s %#v", e.Level.OpenTag(), k, e.Level.CloseTag(), v)
		}
	}
	cli.Print(str)
	for _, ch := range Channel {
		go send(ch, e)
	}
}

func send(ch chan *Entity, e *Entity) {
	ch <- e
}
