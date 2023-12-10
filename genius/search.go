package genius

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rramiachraf/dumb/util"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type response struct {
	Response struct {
		Sections Sections
	}
}

type Result struct {
	ArtistNames string `json:"artist_names"`
	Title       string
	Path        string
	Thumbnail   string `json:"song_art_image_thumbnail_url"`
}

type Hits []struct {
	Result Result
}

type Sections []struct {
	Type string
	Hits Hits
}

type SearchRender struct {
	Query    string
	Sections Sections
}

func Search(query string) SearchRender {
	url := fmt.Sprintf(`https://genius.com/api/search/multi?q=%s`, url.QueryEscape(query))
	res, err := util.SendRequest(url)
	if err != nil {
		logger.Errorln(err)
	}

	defer res.Body.Close()

	var data response

	d := json.NewDecoder(res.Body)
	d.Decode(&data)

	return SearchRender{query, data.Response.Sections}
}
