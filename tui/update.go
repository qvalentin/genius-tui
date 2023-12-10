package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rramiachraf/dumb/genius"
)

var appStyle = lipgloss.NewStyle().Padding(1, 2)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.searchPageModel.list.SetSize(msg.Width-h, msg.Height-v)

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.is_search_open {
				m.is_search_open = false
				search_results := genius.Search(m.textInput.Value())
				m.is_search_results_open = true

				m.searchPageModel.list.Title = "Results for " + search_results.Query
				return m, m.searchPageModel.list.SetItems(extractSearchHits(search_results))
			}
			if m.is_search_results_open {
				i, ok := m.searchPageModel.list.SelectedItem().(SearchHit)
				if ok {
					song, err := genius.Lyrics(i.Path)

					m.err = err
					m.selectedSong = song
				}
				m.is_search_results_open = false
				m.is_lyrics_page_open = true
			}
		}

		if m.is_search_open {
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	}
	m.searchPageModel.list, cmd = m.searchPageModel.list.Update(msg)
	return m, cmd
}
