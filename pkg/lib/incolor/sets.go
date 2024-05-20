package incolor

// Types

type TStColorTerm struct {
	characteristics []TAttribute
	notColor        *bool
}

type TAttribute int

// Constants

const escSimbol = "\x1b"

// Attributes
const (
	Reset TAttribute = iota
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
	ResetBold TAttribute = iota + 22
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
	Black TAttribute = iota + 30
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
	BlackH TAttribute = iota + 90
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
// 	BlackBack TAttribute = iota + 40
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
// 	BlackHBack TAttribute = iota + 100
// 	RedHBack
// 	GreenHBack
// 	YellowHBack
// 	BlueHBack
// 	MagentaHBack
// 	CyanHBack
// 	WhiteHBack
// )
