package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	tea "github.com/charmbracelet/bubbletea"
	tui "github.com/rramiachraf/dumb/tui"
)

var logger = logrus.New()

func main() {

	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
