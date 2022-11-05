package anime

import (
	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
)

// GetAnimeByID ...
func GetAnimeByID(animeID types.AnimeID) (*Anime, error) {
	anime, err := db.GetAnimeByID(animeID)
	if err != nil {
		return nil, err
	}
	animeAns := &Anime{ID: anime.ID,
		Title:       anime.Title,
		OriginTitle: anime.OriginTitle,
		Genres:      anime.Genres,
		Images:      anime.Images,
		Description: anime.Description,
		Director:    anime.Director,
		DateAdded:   anime.DateAdded,
		ReleaseDate: anime.ReleaseDate,
		Author:      anime.Author,
	}
	return animeAns, nil
}
