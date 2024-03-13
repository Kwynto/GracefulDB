package colorterm

// Types

type ColorTerm struct {
	characteristics []Attribute
	notColor        *bool
}

type Attribute int

// Constants

const escSimbol = "\x1b"

// Attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	ResetBold Attribute = iota + 22
	ResetItalic
	ResetUnderline
	ResetBlinking
	_
	ResetReversed
	ResetConcealed
	ResetCrossedOut
)

// Text colors
const (
	Black Attribute = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Hi-Intensity text colors
const (
	BlackH Attribute = iota + 90
	RedH
	GreenH
	YellowH
	BlueH
	MagentaH
	CyanH
	WhiteH
)

// // Background colors
// const (
// 	BlackBack Attribute = iota + 40
// 	RedBack
// 	GreenBack
// 	YellowBack
// 	BlueBack
// 	MagentaBack
// 	CyanBack
// 	WhiteBack
// )

// // Background Hi-Intensity colors
// const (
// 	BlackHBack Attribute = iota + 100
// 	RedHBack
// 	GreenHBack
// 	YellowHBack
// 	BlueHBack
// 	MagentaHBack
// 	CyanHBack
// 	WhiteHBack
// )
