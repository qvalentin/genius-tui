package tui

import (
	"fmt"
)

func (m model) View() string {

	if m.err != nil {
		return fmt.Sprintf("Error %e", m.err)
	}

	// The header
	s := "Dumb cli - use genius with style\n\n"
	if m.is_search_open {

		s += fmt.Sprintf(
			"Enter search\n\n%s\n\n",
			m.textInput.View(),
		)
	}

	if m.is_search_results_open {

		s += m.ViewSeachResultsPage()

	}

	if m.is_lyrics_page_open {
		s += fmt.Sprintf(
			"Lyrics of %s \n\n%s\n\n%s",
			m.selectedSong.Title,
			m.selectedSong.Lyrics,
			"(esc to quit)",
		) + "\n"
	}

	// Iterate over our choices

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
