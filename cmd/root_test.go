/*
Copyright © 2026 SATRIO BAYU AJI <halosatrio@gmail.com>
MIT License
*/
package cmd

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestFlagColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected tcell.Color
	}{
		{"Black color", "black", tcell.ColorBlack},
		{"White color", "white", tcell.ColorWhite},
		{"Blue color", "blue", tcell.ColorBlue},
		{"Cyan color", "cyan", tcell.ColorDarkCyan},
		{"Green color", "green", tcell.ColorGreen},
		{"Magenta color", "magenta", tcell.ColorDarkMagenta},
		{"Red color", "red", tcell.ColorRed},
		{"Yellow color", "yellow", tcell.ColorYellow},
		{"Gray color", "gray", tcell.ColorDarkGray},
		{"Empty string", "", tcell.ColorGreen},
		{"Invalid color", "invalid", tcell.ColorGreen},
		{"Random color", "purple", tcell.ColorNames["purple"]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flagColor(tt.input)
			if result != tt.expected {
				t.Errorf("flagColor(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBigDigitMap(t *testing.T) {
	digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':'}

	for _, digit := range digits {
		t.Run(string(digit), func(t *testing.T) {
			if _, ok := bigDigitMap[digit]; !ok {
				t.Errorf("digit %q not found in bigDigitMap", digit)
			}
		})
	}
}

func TestDrawNumber(t *testing.T) {
	s := tcell.NewSimulationScreen("")
	if s == nil {
		t.Fatal("failed to create simulation screen")
	}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Fini()

	style := tcell.StyleDefault

	digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':'}

	for _, digit := range digits {
		t.Run(string(digit), func(t *testing.T) {
			s.Clear()
			drawNumber(s, 80, 24, 0, digit, style)

			w, h := s.Size()
			hasContent := false
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					c, _, _, _ := s.GetContent(x, y)
					if c != 0 {
						hasContent = true
						break
					}
				}
				if hasContent {
					break
				}
			}
			if !hasContent {
				t.Errorf("drawNumber did not draw anything for digit %q", digit)
			}
		})
	}
}

func TestDrawString(t *testing.T) {
	s := tcell.NewSimulationScreen("")
	if s == nil {
		t.Fatal("failed to create simulation screen")
	}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Fini()

	style := tcell.StyleDefault
	testCases := []string{"test", "123", "", "a b c"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			s.Clear()
			drawString(s, 0, 0, tc, style)

			for i, r := range tc {
				c, _, _, _ := s.GetContent(i, 0)
				if c != r {
					t.Errorf("expected %q at position %d, got %q", r, i, c)
				}
			}
		})
	}
}

func TestDrawMeridiem(t *testing.T) {
	s := tcell.NewSimulationScreen("")
	if s == nil {
		t.Fatal("failed to create simulation screen")
	}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Fini()

	style := tcell.StyleDefault

	testCases := []string{"AM", "PM", "am", "pm"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			s.Clear()
			drawMeridiem(s, tc, 80, 24, 0, style)

			w, h := s.Size()
			hasContent := false
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					c, _, _, _ := s.GetContent(x, y)
					if c != 0 {
						hasContent = true
						break
					}
				}
				if hasContent {
					break
				}
			}
			if !hasContent {
				t.Errorf("drawMeridiem did not draw anything for %q", tc)
			}
		})
	}
}

func TestBigDigitFunctions(t *testing.T) {
	s := tcell.NewSimulationScreen("")
	if s == nil {
		t.Fatal("failed to create simulation screen")
	}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Fini()

	ctx := drawContext{
		s:      s,
		width:  80,
		height: 24,
		diff:   0,
		style:  tcell.StyleDefault,
	}

	tests := []struct {
		name string
		fn   func(drawContext)
	}{
		{"bigZero", bigZero},
		{"bigOne", bigOne},
		{"bigTwo", bigTwo},
		{"bigThree", bigThree},
		{"bigFour", bigFour},
		{"bigFive", bigFive},
		{"bigSix", bigSix},
		{"bigSeven", bigSeven},
		{"bigEight", bigEight},
		{"bigNine", bigNine},
		{"bigColon", bigColon},
		{"bigA", func(ctx drawContext) { bigA(ctx.s, ctx.width, ctx.height, ctx.diff, ctx.style) }},
		{"bigP", func(ctx drawContext) { bigP(ctx.s, ctx.width, ctx.height, ctx.diff, ctx.style) }},
		{"bigM", func(ctx drawContext) { bigM(ctx.s, ctx.width, ctx.height, ctx.diff, ctx.style) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Clear()
			tt.fn(ctx)

			w, h := s.Size()
			hasContent := false
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					c, _, _, _ := s.GetContent(x, y)
					if c != 0 {
						hasContent = true
						break
					}
				}
				if hasContent {
					break
				}
			}
			if !hasContent {
				t.Errorf("%s did not draw anything", tt.name)
			}
		})
	}
}

func TestClockDiff(t *testing.T) {
	expected := [8]int{0, 7, 6, 6, 7, 6, 6, 7}
	if ClockDiff != expected {
		t.Errorf("ClockDiff = %v, want %v", ClockDiff, expected)
	}
}
