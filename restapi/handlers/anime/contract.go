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
	ReleaseDate *int64   `validate:"required" json:"releaseDate"`
}

type EditAnimeRequest struct {
	AnimeCreateRequest
}

func (a *AnimeCreateRequest) NewAnimeModel() *anime.Anime {
	return &anime.Anime{
		Title:       a.Title,
		OriginTitle: a.OriginTitle,
		Description: a.Description,
		Director:    a.Director,
		Genres:      a.Genres,
		ReleaseDate: time.Unix(*a.ReleaseDate, 0),
	}
}

type CreateAnimeResponce struct {
	AnimeID types.AnimeID `json:"animeId"`
}

type AnimeByListIDRequest struct {
	IDs []types.AnimeID `json:"ids" validate:"required"`
}

type AnimeByListIDResponce struct {
	Animes []*anime.Anime `json:"animes"`
}

type GetAnimesRequest struct {
	SortBy string   `json:"sortBy" validate:"required"`
	Genres []string `json:"genres" validate:"required"`
	Order  *int8    `json:"order" validate:"required"`
}

type MostPopular struct {
	Count *int `json:"count" validate:"required,gte=1,lte=50"`
}
