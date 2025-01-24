/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
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

var SmallStopWatchDiff = [11]int{7, 6, 6, 7, 6, 6, 7, 6, 4, 7, 7}

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

	Run: func(cmd *cobra.Command, args []string) {
		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorTomato)

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
		go func() {
			for {
				s.Clear()

				color := tcell.StyleDefault.Background(tcell.ColorTomato).Foreground(tcell.Color107)

				// Print time on terminal
				nowTime := time.Now()
				hours := nowTime.Hour()
				minutes := nowTime.Minute()
				seconds := nowTime.Second()
				var formattedTime string

				diff := -26
				totalString := 8

				// toggle second flags
				if Seconds {
					formattedTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
				} else {
					formattedTime = fmt.Sprintf("%02d:%02d", hours, minutes)
				}

				drawString(s, termWidth/2, termHeight/2+5, formattedTime, defStyle)

				for i := 0; i < totalString; i++ {
					if i != 0 {
						diff = diff + SmallStopWatchDiff[i-1]
					}
					if i == 2 || i == 5 {
						drawNumber(s, termWidth, termHeight, diff, formattedTime[i], color)
					} else {
						drawNumber(s, termWidth, termHeight, diff, formattedTime[i], color)
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
	// for _, r := range []rune(text) {
	for _, r := range text {
		s.SetContent(x, y, r, nil, style)
		x++
	}
}

func bigColon(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff+2, height/2-1, "  ", style)
	drawString(s, width/2+diff+2, height/2+1, "  ", style)
}

func bigOne(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff+4, height/2-2, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+2, "  ", style)
}
func bigTwo(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)
	drawString(s, width/2+diff, height/2, "      ", style)
	drawString(s, width/2+diff, height/2+1, "  ", style)
	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func bigThree(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)
	drawString(s, width/2+diff, height/2, "      ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func bigFour(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "  ", style)
	drawString(s, width/2+diff+4, height/2-2, "  ", style)

	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)

	drawString(s, width/2+diff, height/2, "      ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+2, "  ", style)
}
func bigFive(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)
	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff, height/2, "      ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func bigSix(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)
	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff, height/2, "      ", style)

	drawString(s, width/2+diff, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)

	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func bigSeven(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+2, "  ", style)
}
func bigEight(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)

	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)

	drawString(s, width/2+diff, height/2, "      ", style)

	drawString(s, width/2+diff, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)

	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func bigNine(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)

	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)

	drawString(s, width/2+diff, height/2, "      ", style)

	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func bigZero(s tcell.Screen, width, height, diff int, style tcell.Style) {
	drawString(s, width/2+diff, height/2-2, "      ", style)
	drawString(s, width/2+diff, height/2-1, "  ", style)
	drawString(s, width/2+diff+4, height/2-1, "  ", style)
	drawString(s, width/2+diff, height/2, "  ", style)
	drawString(s, width/2+diff+4, height/2, "  ", style)
	drawString(s, width/2+diff, height/2+1, "  ", style)
	drawString(s, width/2+diff+4, height/2+1, "  ", style)
	drawString(s, width/2+diff, height/2+2, "      ", style)
}
func drawNumber(s tcell.Screen, termWidth, termHeight, diff int, nowTime byte, style tcell.Style) {
	switch nowTime {
	case '0':
		bigZero(s, termWidth, termHeight, diff, style)
	case '1':
		bigOne(s, termWidth, termHeight, diff, style)
	case '2':
		bigTwo(s, termWidth, termHeight, diff, style)
	case '3':
		bigThree(s, termWidth, termHeight, diff, style)
	case '4':
		bigFour(s, termWidth, termHeight, diff, style)
	case '5':
		bigFive(s, termWidth, termHeight, diff, style)
	case '6':
		bigSix(s, termWidth, termHeight, diff, style)
	case '7':
		bigSeven(s, termWidth, termHeight, diff, style)
	case '8':
		bigEight(s, termWidth, termHeight, diff, style)
	case '9':
		bigNine(s, termWidth, termHeight, diff, style)
	case ':':
		bigColon(s, termWidth, termHeight, diff, style)
	}
}
