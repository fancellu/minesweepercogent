package main

import (
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
	_ "embed"
)

type MyButton struct {
	core.Button

	Flag     bool
	Mine     bool
	ShowMine bool
}

func (mb *MyButton) Init() {
	mb.Button.Init()

	mb.Flag = false
	mb.Mine = false
	mb.ShowMine = false

	// fix for icon size
	tree.AddChildInit(mb, "icon", func(w *core.Icon) {
		w.Styler(func(s *styles.Style) {
			s.Font.Size.Dp(30)
		})
	})

	mb.Styler(func(s *styles.Style) {
		s.Border.Width.Set(units.Dp(1))
		s.Border.Radius.Zero()
		s.Padding.Set(units.Dp(10))
		size := units.Dp(45)
		s.Min.Set(size, size)
		s.Max.Set(size, size)
		s.Font.Size.Dp(25)
		if s.Is(states.Checked) {
			s.Background = colors.Uniform(colors.Orange)
			s.Color = colors.Uniform(colors.Blue)
		} else if mb.Mine && mb.ShowMine {
			s.Background = colors.Uniform(colors.Red)
			s.Color = colors.Uniform(colors.Black)
		}

	})
	mb.SetType(core.ButtonAction)
	mb.SetIcon(icons.Blank)
}

//go:embed mine.svg
var mineSVG string

func (mb *MyButton) ShowMineIcon() {
	mb.SetText("")
	mb.SetIcon(icons.Icon(mineSVG))
	mb.ShowMine = true
	mb.Button.Update()
}

//go:embed flag2.svg
var flagSVG string

func (mb *MyButton) ShowFlagIcon() {
	mb.Flag = true

	mb.SetIcon(icons.Icon(flagSVG))
	mb.Button.Update()
}
