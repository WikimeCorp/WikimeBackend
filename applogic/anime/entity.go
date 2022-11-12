package anime

import (
	"time"

	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
)

type Anime struct {
	ID          types.AnimeID `json:"id"`
	Title       string        `json:"title"`
	OriginTitle string        `json:"originTitle"`
	Genres      []string      `json:"genres"`
	Description string        `json:"description"`
	Images      []string      `json:"images"`
	Director    string        `json:"director"`
	DateAdded   time.Time     `json:"dataAded"`
	ReleaseDate time.Time     `json:"releaseDate"`
	Author      types.UserID  `json:"author"`
	Rating      Rating        `json:"rating"`
}

type Rating struct {
	Five        uint32  `json:"five"`
	Four        uint32  `json:"four"`
	Three       uint32  `json:"three"`
	Two         uint32  `json:"two"`
	One         uint32  `json:"one"`
	InFavorites uint32  `json:"inFavorites"`
	Average     float64 `json:"average"`
}

func (a *Anime) NewDBModel() *dbtypes.Anime {
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
	}
}
