package incolor

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
	isWithoutColor = colorIsntSet() || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	IoOutput = colorable.NewColorableStdout()
	IoError  = colorable.NewColorableStderr()

	mCache  = make(map[TAttribute]*TStColorTerm)
	mxBlock sync.Mutex
)

var mapResetAttributes = map[TAttribute]TAttribute{
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

func (ct *TStColorTerm) seq() string {
	slFormat := make([]string, len(ct.characteristics))
	for i, v := range ct.characteristics {
		slFormat[i] = strconv.Itoa(int(v))
	}

	return strings.Join(slFormat, ";")
}

func (c *TStColorTerm) formating() string {
	return fmt.Sprintf("%s[%sm", escSimbol, c.seq())
}

func (c *TStColorTerm) unformating() string {
	slFormat := make([]string, len(c.characteristics))
	for i, v := range c.characteristics {
		slFormat[i] = strconv.Itoa(int(Reset))
		ra, ok := mapResetAttributes[v]
		if ok {
			slFormat[i] = strconv.Itoa(int(ra))
		}
	}

	return fmt.Sprintf("%s[%sm", escSimbol, strings.Join(slFormat, ";"))
}

func (c *TStColorTerm) dontSetColor() bool {
	if c.notColor != nil {
		return *c.notColor
	}

	return isWithoutColor
}

func (ct *TStColorTerm) roll(str string) string {
	if ct.dontSetColor() {
		return str
	}

	return ct.formating() + str + ct.unformating()
}

func (ct *TStColorTerm) FSprint() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return ct.roll(fmt.Sprint(a...))
	}
}

func (ct *TStColorTerm) FSprintf() func(format string, a ...interface{}) string {
	return func(format string, a ...interface{}) string {
		return ct.roll(fmt.Sprintf(format, a...))
	}
}

func (ct *TStColorTerm) AddAttr(value ...TAttribute) *TStColorTerm {
	ct.characteristics = append(ct.characteristics, value...)
	return ct
}

func NewCT(value ...TAttribute) *TStColorTerm {
	ct := &TStColorTerm{
		characteristics: make([]TAttribute, 0),
	}

	if colorIsntSet() {
		ct.notColor = boolLink(true)
	}

	ct.AddAttr(value...)
	return ct
}

func colorStr(format string, p TAttribute, a ...interface{}) string {
	mxBlock.Lock()
	defer mxBlock.Unlock()

	c, ok := mCache[p]
	if !ok {
		c = NewCT(p)
		mCache[p] = c
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
	var uOutMode uint32
	out := windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(out, &uOutMode); err != nil {
		return
	}
	uOutMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(out, uOutMode)
}
