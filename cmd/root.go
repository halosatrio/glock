/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glock",
	Short: "a digital clock in your terminal  ",
	Long: `a digital clock in your terminal, inspired by tty-clock.

Default to 12-hour local time, no seconds

Usage: glock [OPTIONS]

Options:
  -s, --second
        Display seconds
  
  -d, --date
        Display date with format YYYY-MM-DD
  
  -c, --color <COLOR>

  -h, --help
        Print help
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		// boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

		// Initialize screen
		s, err := tcell.NewScreen()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if err := s.Init(); err != nil {
			log.Fatalf("%+v", err)
		}
		s.SetStyle(defStyle)
		s.EnableMouse()
		s.EnablePaste()
		s.Clear()
		termWidth, termHeight := s.Size()

		// Draw initial boxes
		// drawBox(s, 1, 1, 42, 7, boxStyle, "Click and drag to draw a box")
		// drawBox(s, 5, 9, 32, 14, boxStyle, "Press C to reset")

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

		// Here's how to get the screen size when you need it.
		// xmax, ymax := s.Size()

		// Here's an example of how to inject a keystroke where it will
		// be picked up by the next PollEvent call.  Note that the
		// queue is LIFO, it has a limited length, and PostEvent() can
		// return an error.
		// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

		// Event loop
		// ox, oy := -1, -1
		for {
			// Update screen
			s.Show()
			drawString(s, termWidth/2, termHeight/2, "satrio", defStyle)

			// Poll event
			ev := s.PollEvent()

			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				} else if ev.Key() == tcell.KeyCtrlL {
					s.Sync()
				} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.Clear()
				}
				// case *tcell.EventMouse:
				// 	x, y := ev.Position()

				// 	switch ev.Buttons() {
				// 	case tcell.Button1, tcell.Button2:
				// 		if ox < 0 {
				// 			ox, oy = x, y // record location when click started
				// 		}

				// 	case tcell.ButtonNone:
				// 		if ox >= 0 {
				// 			label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
				// 			drawBox(s, ox, oy, x, y, boxStyle, label)
				// 			ox, oy = -1, -1
				// 		}
				// 	}
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.Flags().BoolVarP(&Seconds, "second", "s", false, "display seconds")

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.glock.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func drawString(s tcell.Screen, x, y int, text string, style tcell.Style) {
	for _, r := range []rune(text) {
		s.SetContent(x, y, r, nil, style)
		x++
	}
}

// func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
// 	row := y1
// 	col := x1
// 	for _, r := range []rune(text) {
// 		s.SetContent(col, row, r, nil, style)
// 		col++
// 		if col >= x2 {
// 			row++
// 			col = x1
// 		}
// 		if row > y2 {
// 			break
// 		}
// 	}
// }

// func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
// 	if y2 < y1 {
// 		y1, y2 = y2, y1
// 	}
// 	if x2 < x1 {
// 		x1, x2 = x2, x1
// 	}

// 	// Fill background
// 	for row := y1; row <= y2; row++ {
// 		for col := x1; col <= x2; col++ {
// 			s.SetContent(col, row, ' ', nil, style)
// 		}
// 	}

// 	// Draw borders
// 	for col := x1; col <= x2; col++ {
// 		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
// 		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
// 	}
// 	for row := y1 + 1; row < y2; row++ {
// 		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
// 		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
// 	}

// 	// Only draw corners if necessary
// 	if y1 != y2 && x1 != x2 {
// 		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
// 		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
// 		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
// 		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
// 	}

// 	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
// }
