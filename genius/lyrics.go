package genius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"jaytaylor.com/html2text"

	"github.com/PuerkitoBio/goquery"
	"github.com/rramiachraf/dumb/util"
)

type Song struct {
	Artist  string
	Title   string
	Image   string
	Lyrics  string
	Credits map[string]string
	About   [2]string
}

type songResponse struct {
	Response struct {
		Song struct {
			ArtistNames string `json:"artist_names"`
			Image       string `json:"song_art_image_thumbnail_url"`
			Title       string
			Description struct {
				Plain string
			}
			CustomPerformances []customPerformance `json:"custom_performances"`
		}
	}
}

type customPerformance struct {
	Label   string
	Artists []struct {
		Name string
	}
}

func (s *Song) parseLyrics(doc *goquery.Document) {
	doc.Find("[data-lyrics-container='true']").Each(func(i int, ss *goquery.Selection) {
		h, err := ss.Html()
		if err != nil {
			logger.Errorln("unable to parse lyrics", err)
		}
		s.Lyrics += h
	})

	plain, err := html2text.FromString(s.Lyrics, html2text.Options{
		PrettyTables: true,
		OmitLinks:    true,
		TextOnly:     false,
	})
	if err != nil {
		panic(err)
	}
	s.Lyrics = plain

}

func (s *Song) parseSongData(doc *goquery.Document) {
	attr, exists := doc.Find("meta[property='twitter:app:url:iphone']").Attr("content")
	if exists {
		songID := strings.Replace(attr, "genius://songs/", "", 1)
		u := fmt.Sprintf("https://genius.com/api/songs/%s?text_format=plain", songID)

		res, err := util.SendRequest(u)
		if err != nil {
			logger.Errorln(err)
		}

		defer res.Body.Close()

		var data songResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&data)
		if err != nil {
			logger.Errorln(err)
		}

		songData := data.Response.Song
		s.Artist = songData.ArtistNames
		s.Image = songData.Image
		s.Title = songData.Title
		s.About[0] = songData.Description.Plain
		s.About[1] = truncateText(songData.Description.Plain)
		s.Credits = make(map[string]string)

		for _, perf := range songData.CustomPerformances {
			var artists []string
			for _, artist := range perf.Artists {
				artists = append(artists, artist.Name)
			}
			s.Credits[perf.Label] = strings.Join(artists, ", ")
		}
	}
}

func truncateText(text string) string {
	textArr := strings.Split(text, "")

	if len(textArr) > 250 {
		return strings.Join(textArr[0:250], "") + "..."
	}

	return text
}

func (s *Song) parse(doc *goquery.Document) {
	s.parseLyrics(doc)
	s.parseSongData(doc)
}

func Lyrics(id string) (Song, error) {

	url := fmt.Sprintf("https://genius.com%s", id)
	resp, err := util.SendRequest(url)
	if err != nil {
		logger.Errorln(err)
		return Song{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Song{}, fmt.Errorf("Not found")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Errorln(err)
		return Song{}, err
	}

	cf := doc.Find(".cloudflare_content").Length()
	if cf > 0 {
		logger.Errorln("cloudflare got in the way")
		return Song{}, err
	}

	var s Song
	s.parse(doc)

	return s, nil
}
