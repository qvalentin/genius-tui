package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rramiachraf/dumb/genius"
)

type model struct {
	textInput              textinput.Model
	is_search_open         bool
	is_search_results_open bool
	is_lyrics_page_open    bool
	search_results         genius.SearchRender
	choices                []string
	cursor                 int
	selected               map[int]struct{}
	searchPageModel        SearchPageModel
	selectedSong           genius.Song
	err                    error
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "search"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m := model{
		textInput:      ti,
		is_search_open: true,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		searchPageModel: SearchPageModel{
			list: list.New((make([]list.Item, 1)), list.NewDefaultDelegate(), 3, 4),
		},
	}
	m.searchPageModel.list.Title = "Search results"
	m.searchPageModel.list.SetShowTitle(true)
	return m

}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
