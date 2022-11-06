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
