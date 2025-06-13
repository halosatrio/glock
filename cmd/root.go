/*
Copyright © 2025 SATRIO BAYU AJI <halosatrio@gmail.com>
MIT License
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

// give spacing on block
// default block width is 6
// block 2,5,8 have extra 1 spacing
var ClockDiff = [8]int{0, 7, 6, 6, 7, 6, 6, 7}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glock",
	Short: "a digital clock in your terminal  ",
	Long: `a digital clock in your terminal, inspired by tty-clock.

Default to 24-hour local time, no seconds

Usage: glock [FLAGS]

Flags:
  -s, --second
        Display seconds
  
  -m, --meridiem
        Display time in 12 hour format
  
  -c, --color <COLOR>

  -h, --help
        Print help
`,

	Run: func(cmd *cobra.Command, args []string) {
		// Initialize screen
		s, err := tcell.NewScreen()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if err := s.Init(); err != nil {
			log.Fatalf("%+v", err)
		}
		s.EnableMouse()
		s.EnablePaste()
		s.Clear()

		quit := func() {
			// You have to catch panics in a defer, clean up, and
			// re-raise them - otherwise your application can
			// die without leaving any diagnostic trace.
			maybePanic := recover()
			s.Fini()
			if maybePanic != nil {
				panic(maybePanic)
			}
		}
		defer quit()

		// Here's an example of how to inject a keystroke where it will
		// be picked up by the next PollEvent call.  Note that the
		// queue is LIFO, it has a limited length, and PostEvent() can
		// return an error.
		// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

		// Event loop
		// goroutine
		// TODO:
		// fix why the program flicker (sometimes)
		// on Clear() function maybe, or the Sync()
		go func() {
			for {
				// terminal needs to be clear every second
				// if not the previous block will keep appearing on screen
				s.Clear()

				// get terminal size every second
				// to update the clock position when terminal size is changed
				// position will be updated on the next second, not instantly
				termWidth, termHeight := s.Size()

				// set color by flag
				var color tcell.Style
				color = color.Background(flagColor(FlagColor)).Foreground(flagColor(FlagColor))

				// get current time
				nowTime := time.Now()

				var formattedTime string
				formattedDate := fmt.Sprintf("%02d-%02d-%02d", nowTime.Year(), nowTime.Month(), nowTime.Day())
				totalString, timeDiff := 0, 0

				// toggle second flag
				// toggle 12 hour format flag
				if Meridiem {
					formattedTime = fmt.Sprint(nowTime.Format("03:04:05PM"))
					ampm := formattedTime[len(formattedTime)-2:]
					if Seconds {
						totalString, timeDiff = 8, -34
						drawMeridiem(s, ampm, termWidth, termHeight, 24, color)
					} else {
						totalString, timeDiff = 5, -26
						drawMeridiem(s, ampm, termWidth, termHeight, 14, color)
					}
				} else {
					if Seconds {
						totalString, timeDiff = 8, -25
						formattedTime = fmt.Sprintf("%02d:%02d:%02d", nowTime.Hour(), nowTime.Minute(), nowTime.Second())
					} else {
						totalString, timeDiff = 5, -17
						formattedTime = fmt.Sprintf("%02d:%02d", nowTime.Hour(), nowTime.Minute())
					}
				}

				// draw date in the center of terminal below the clock
				drawString(s, termWidth/2-5, termHeight/2+5, formattedDate, color.Background(tcell.ColorReset))

				for i := 0; i < totalString; i++ {
					if i != 0 {
						timeDiff = timeDiff + ClockDiff[i]
					}
					if i == 2 || i == 5 {
						drawNumber(s, termWidth, termHeight, timeDiff, formattedTime[i], color)
					} else {
						drawNumber(s, termWidth, termHeight, timeDiff, formattedTime[i], color)
					}
				}

				// update the screen every second
				// Sleep() is used to stop the goroutine
				s.Sync()
				time.Sleep(time.Second)
			}
		}()

	mainloop:
		for {
			// Poll event
			ev := s.PollEvent()

			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					break mainloop
				}
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var Seconds bool
var Meridiem bool
var FlagColor string

func init() {
	// boolean flags only need the flag without value
	rootCmd.Flags().BoolVarP(&Seconds, "second", "s", false, "display seconds")
	rootCmd.Flags().BoolVarP(&Meridiem, "meridiem", "m", false, "display 12 hour format")

	// string flags need the flag and the value. example: --color magenta
	rootCmd.Flags().StringVarP(&FlagColor, "color", "c", "", "choose clock color")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// function to change color of the clock
func flagColor(ColorString string) tcell.Color {
	switch ColorString {
	case "black":
		return tcell.ColorBlack
	case "white":
		return tcell.ColorWhite
	case "blue":
		return tcell.ColorBlue
	case "cyan":
		return tcell.ColorDarkCyan
	case "green":
		return tcell.ColorGreen
	case "magenta":
		return tcell.ColorDarkMagenta
	case "red":
		return tcell.ColorRed
	case "yellow":
		return tcell.ColorYellow
	case "gray":
		return tcell.ColorDarkGray
	default:
		if ColorString == "" {
			return tcell.ColorGreen
		} else {
			// can use int8 value 0 < 255 (total color number supported)
			// num, _ := strconv.Atoi(ColorString)
			// return tcell.PaletteColor(num)
			if tcell.ColorNames[ColorString].Valid() {
				return tcell.ColorNames[ColorString]
			} else {
				// dont know how to exit the program and show error
				// for the time being, invalid color will return default color = green
				return tcell.ColorGreen
			}
		}
	}
}

func drawString(s tcell.Screen, x, y int, text string, style tcell.Style) {
	for _, r := range text {
		s.SetContent(x, y, r, nil, style)
		x++
	}
}

// function print am/pm
func drawMeridiem(s tcell.Screen, val string, width, height, diff int, style tcell.Style) {
	if val == "AM" {
		bigA(s, width, height, diff, style)
		bigM(s, width, height, diff+7, style)
	} else {
		bigP(s, width, height, diff, style)
		bigM(s, width, height, diff+7, style)
	}
}

// Big chunks of code just to draw block numbers on screen
// the block numbers are drawn in 6cols x 5rows grid
//
//	..00..
//	0000..
//	..00..
//	..00..
//	000000

// this is the old function before refactor
func bigA(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff+1, height/2-2, "    ", style)

	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)

	drawString(s, width/2+diff, height/2, "      ", style)

	drawString(s, width/2+diff, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)

	drawString(s, width/2+diff, height/2+2, "  ", style)
	drawString(s, width/2+diff+4, height/2+2, "  ", style)
}
func bigP(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)

	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)

	drawString(s, width/2+diff, height/2, "      ", style)

	drawString(s, width/2+diff, height/2+1, "  ", style)

	drawString(s, width/2+diff, height/2+2, "  ", style)
}
func bigM(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "  ", style)
	drawString(s, width/2+diff+4, height/2-2, "  ", style)

	drawString(s, width/2+diff, height/2-1, "      ", style)

	drawString(s, width/2+diff, height/2, "  ", style)
	drawString(s, width/2+diff+4, height/2, "  ", style)

	drawString(s, width/2+diff, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)

	drawString(s, width/2+diff, height/2+2, "  ", style)
	drawString(s, width/2+diff+4, height/2+2, "  ", style)
}

// this chunk of code bellow has been refactored
type drawContext struct {
	s      tcell.Screen
	width  int
	height int
	diff   int
	style  tcell.Style
}
type charPart struct {
	dx, dy int
	text   string
}

func drawBigChar(ctx drawContext, parts []charPart) {
	for _, p := range parts {
		x := ctx.width/2 + ctx.diff + p.dx
		y := ctx.height/2 + p.dy
		drawString(ctx.s, x, y, p.text, ctx.style)
	}
}

func bigColon(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{2, -1, "  "},
		{2, 1, "  "},
	})
}
func bigOne(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{2, -2, "  "},
		{0, -1, "    "},
		{2, 0, "  "},
		{2, 1, "  "},
		{0, 2, "      "},
	})
}
func bigTwo(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{4, -1, "  "},
		{0, 0, "      "},
		{0, 1, "  "},
		{0, 2, "      "},
	})
}
func bigThree(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{4, -1, "  "},
		{0, 0, "      "},
		{4, 1, "  "},
		{0, 2, "      "},
	})
}
func bigFour(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "  "},
		{4, -2, "  "},
		{0, -1, "  "},
		{4, -1, "  "},
		{0, 0, "      "},
		{4, 1, "  "},
		{4, 2, "  "},
	})
}
func bigFive(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{0, -1, "  "},
		{0, 0, "      "},
		{4, 1, "  "},
		{0, 2, "      "},
	})
}
func bigSix(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{0, -1, "  "},
		{0, 0, "      "},
		{0, 1, "  "},
		{4, 1, "  "},
		{0, 2, "      "},
	})
}
func bigSeven(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{4, -1, "  "},
		{4, 0, "  "},
		{4, 1, "  "},
		{4, 2, "  "},
	})
}
func bigEight(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{0, -1, "  "},
		{4, -1, "  "},
		{0, 0, "      "},
		{0, 1, "  "},
		{4, 1, "  "},
		{0, 2, "      "},
	})
}
func bigNine(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{0, -1, "  "},
		{4, -1, "  "},
		{0, 0, "      "},
		{4, 1, "  "},
		{0, 2, "      "},
	})
}
func bigZero(ctx drawContext) {
	drawBigChar(ctx, []charPart{
		{0, -2, "      "},
		{0, -1, "  "},
		{4, -1, "  "},
		{0, 0, "  "},
		{4, 0, "  "},
		{0, 1, "  "},
		{4, 1, "  "},
		{0, 2, "      "},
	})
}

type drawFunc func(ctx drawContext)

var bigDigitMap = map[byte]drawFunc{
	'0': bigZero,
	'1': bigOne,
	'2': bigTwo,
	'3': bigThree,
	'4': bigFour,
	'5': bigFive,
	'6': bigSix,
	'7': bigSeven,
	'8': bigEight,
	'9': bigNine,
	':': bigColon,
}

func drawNumber(s tcell.Screen, termWidth, termHeight, diff int, nowTime byte, style tcell.Style) {
	ctx := drawContext{
		s:      s,
		width:  termWidth,
		height: termHeight,
		diff:   diff,
		style:  style,
	}

	if fn, ok := bigDigitMap[nowTime]; ok {
		fn(ctx)
	}
}
