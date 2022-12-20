package anime

import (
	"encoding/json"
	"path"
	"time"

	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
)

type Anime struct {
	ID          types.AnimeID `json:"id"`
	Title       string        `json:"title"`
	OriginTitle string        `json:"originTitle"`
	Genres      []string      `json:"genres"`
	Description string        `json:"description"`
	Poster      *string       `json:"poster"`
	Images      []string      `json:"images"`
	Director    string        `json:"director"`
	DateAdded   time.Time     `json:"dataAdded"`
	ReleaseDate time.Time     `json:"releaseDate"`
	Author      types.UserID  `json:"author"`
	Rating      Rating        `json:"rating"`
}

func (d *Anime) MarshalJSON() ([]byte, error) {
	type Alias Anime
	return json.Marshal(&struct {
		*Alias
		DateAdded   string `json:"dataAdded"`
		ReleaseDate string `json:"releaseDate"`
	}{
		Alias:       (*Alias)(d),
		DateAdded:   d.DateAdded.Format("02.01.2006"),
		ReleaseDate: d.ReleaseDate.Format("02.01.2006"),
	})
}

func (a *Anime) NewDBModel() *dbtypes.Anime {
	defaultPosterPath := path.Join(config.Config.ImagesPathURI, config.Config.DefaultAnimePosterPath)
	return &dbtypes.Anime{
		ID:          a.ID,
		Title:       a.Title,
		OriginTitle: a.OriginTitle,
		Genres:      a.Genres,
		Description: a.Description,
		Director:    a.Director,
		ReleaseDate: a.ReleaseDate,
		Author:      a.Author,
		Images:      []string{},
		Poster:      &defaultPosterPath,
	}
}

type Rating struct {
	Five        uint32  `json:"five"`
	Four        uint32  `json:"four"`
	Three       uint32  `json:"three"`
	Two         uint32  `json:"two"`
	One         uint32  `json:"one"`
	InFavorites uint32  `json:"inFavorites"`
	Average     float64 `json:"average"`
	Watched     uint32  `json:"watched"`
}
