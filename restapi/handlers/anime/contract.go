package anime

import (
	"time"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/types"
)

type AnimeCreateRequest struct {
	Title       string   `validate:"required" json:"title"`
	OriginTitle string   `validate:"required" json:"originTitle"`
	Description string   `validate:"required" json:"description"`
	Director    string   `validate:"required" json:"director"`
	Genres      []string `validate:"required" json:"genres"`
	ReleaseDate int64    `validate:"required" json:"releaseDate"`
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

type AnimeByListIDRequest struct {
	IDs []types.AnimeID `json:"ids" validate:"required"`
}

type AnimeByListIDResponce struct {
	Animes []*anime.Anime `json:"animes"`
}
