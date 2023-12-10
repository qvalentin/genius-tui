package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/rramiachraf/dumb/genius"
)

type SearchHit struct {
	ArtistNames string
	title       string
	Path        string
}

type SearchPageModel struct {
	list list.Model
}

func (i SearchHit) Title() string       { return fmt.Sprintf("[%s] %s %s", i.ArtistNames, i.title, i.Path) }
func (i SearchHit) Description() string { return "song" }
func (i SearchHit) FilterValue() string { return i.title }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) ViewSeachResultsPage() string {
	return docStyle.Render(m.searchPageModel.list.View())
}

func extractSearchHits(input genius.SearchRender) (result []list.Item) {
	result = make([]list.Item, 0)
	for _, section := range input.Sections {
		if section.Type == "song" {
			for _, hit := range section.Hits {
				result = append(result, SearchHit{
					ArtistNames: hit.Result.ArtistNames,
					title:       hit.Result.Title,
					Path:        hit.Result.Path,
				})
			}
		}
	}
	return result
}
