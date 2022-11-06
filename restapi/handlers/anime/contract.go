package anime

import (
	"time"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/types"
)

type AnimeCreateRequest struct {
	Title       string   `validate:"required" validateName:"title"`
	OriginTitle string   `validate:"required" validateName:"originTitle"`
	Description string   `validate:"required" validateName:"description"`
	Director    string   `validate:"required" validateName:"director"`
	Genres      []string `validate:"required" validateName:"genres"`
	ReleaseDate int64    `validate:"required" validateName:"releaseDate"`
}

/*
Title       string   `validate:"required" validateName:"title"`
	OriginTitle string   `validate:"required" validateName:"originTitle"`
	Description string   `validate:"required" validateName:"description"`
	Director    string   `validate:"description,required"`
	Genres      []string `validate:"genres,required"`
	ReleaseDate int64    `validate:"releaseDate,required"`
*/

func (a *AnimeCreateRequest) NewAnimeModel() *anime.Anime {
	return &anime.Anime{
		Title:       a.Title,
		OriginTitle: a.OriginTitle,
		Description: a.Description,
		Director:    a.Description,
		Genres:      a.Genres,
		ReleaseDate: time.Unix(a.ReleaseDate, 0),
	}
}

type AnimeResponce struct {
	AnimeID types.AnimeID `json:"animeId"`
}
