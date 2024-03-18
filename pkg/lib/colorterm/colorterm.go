package colorterm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Kwynto/colorable"
	"github.com/Kwynto/isatty"
	"golang.org/x/sys/windows"
)

var (
	withoutColor = colorIsntSet() || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	Output = colorable.NewColorableStdout()
	Error  = colorable.NewColorableStderr()

	cache = make(map[Attribute]*ColorTerm)
	block sync.Mutex
)

var mapResetAttributes = map[Attribute]Attribute{
	Bold:         ResetBold,
	Faint:        ResetBold,
	Italic:       ResetItalic,
	Underline:    ResetUnderline,
	BlinkSlow:    ResetBlinking,
	BlinkRapid:   ResetBlinking,
	ReverseVideo: ResetReversed,
	Concealed:    ResetConcealed,
	CrossedOut:   ResetCrossedOut,
}

func colorIsntSet() bool {
	return os.Getenv("NO_COLOR") != ""
}

func boolLink(v bool) *bool {
	return &v
}

func (ct *ColorTerm) seq() string {
	format := make([]string, len(ct.characteristics))
	for i, v := range ct.characteristics {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

func (c *ColorTerm) formating() string {
	return fmt.Sprintf("%s[%sm", escSimbol, c.seq())
}

func (c *ColorTerm) unformating() string {
	format := make([]string, len(c.characteristics))
	for i, v := range c.characteristics {
		format[i] = strconv.Itoa(int(Reset))
		ra, ok := mapResetAttributes[v]
		if ok {
			format[i] = strconv.Itoa(int(ra))
		}
	}

	return fmt.Sprintf("%s[%sm", escSimbol, strings.Join(format, ";"))
}

func (c *ColorTerm) dontSetColor() bool {
	if c.notColor != nil {
		return *c.notColor
	}

	return withoutColor
}

func (ct *ColorTerm) roll(str string) string {
	if ct.dontSetColor() {
		return str
	}

	return ct.formating() + str + ct.unformating()
}

func (ct *ColorTerm) FSprint() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return ct.roll(fmt.Sprint(a...))
	}
}

func (ct *ColorTerm) FSprintf() func(format string, a ...interface{}) string {
	return func(format string, a ...interface{}) string {
		return ct.roll(fmt.Sprintf(format, a...))
	}
}

func (ct *ColorTerm) AddAttr(value ...Attribute) *ColorTerm {
	ct.characteristics = append(ct.characteristics, value...)
	return ct
}

func NewCT(value ...Attribute) *ColorTerm {
	ct := &ColorTerm{
		characteristics: make([]Attribute, 0),
	}

	if colorIsntSet() {
		ct.notColor = boolLink(true)
	}

	ct.AddAttr(value...)
	return ct
}

func colorStr(format string, p Attribute, a ...interface{}) string {
	block.Lock()
	defer block.Unlock()

	c, ok := cache[p]
	if !ok {
		c = NewCT(p)
		cache[p] = c
	}

	if len(a) == 0 {
		return c.FSprint()(format)
	}

	return c.FSprintf()(format, a...)
}

func StringBlack(format string, a ...interface{}) string {
	return colorStr(format, Black, a...)
}

func StringRed(format string, a ...interface{}) string {
	return colorStr(format, Red, a...)
}

func StringGreen(format string, a ...interface{}) string {
	return colorStr(format, Green, a...)
}

func StringYellow(format string, a ...interface{}) string {
	return colorStr(format, Yellow, a...)
}

func StringBlue(format string, a ...interface{}) string {
	return colorStr(format, Blue, a...)
}

func StringMagenta(format string, a ...interface{}) string {
	return colorStr(format, Magenta, a...)
}

func StringCyan(format string, a ...interface{}) string {
	return colorStr(format, Cyan, a...)
}

func StringWhite(format string, a ...interface{}) string {
	return colorStr(format, White, a...)
}

func StringBlackH(format string, a ...interface{}) string {
	return colorStr(format, BlackH, a...)
}

func StringRedH(format string, a ...interface{}) string {
	return colorStr(format, RedH, a...)
}

func StringGreenH(format string, a ...interface{}) string {
	return colorStr(format, GreenH, a...)
}

func StringYellowH(format string, a ...interface{}) string {
	return colorStr(format, YellowH, a...)
}

func StringBlueH(format string, a ...interface{}) string {
	return colorStr(format, BlueH, a...)
}

func StringMagentaH(format string, a ...interface{}) string {
	return colorStr(format, MagentaH, a...)
}

func StringCyanH(format string, a ...interface{}) string {
	return colorStr(format, CyanH, a...)
}

func StringWhiteH(format string, a ...interface{}) string {
	return colorStr(format, WhiteH, a...)
}

func init() {
	// https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences
	var outMode uint32
	out := windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(out, &outMode); err != nil {
		return
	}
	outMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(out, outMode)
}
