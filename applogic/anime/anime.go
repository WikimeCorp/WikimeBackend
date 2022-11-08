package anime

import (
	"context"

	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
	"golang.org/x/sync/errgroup"
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

func CreateAnime(anime *Anime) (types.AnimeID, error) {
	animeDb := anime.NewDBModel()
	animeID, err := db.AddAnime(animeDb)
	if err != nil {
		return 0, err
	}

	return animeID, nil
}

func GetAnimesByListID(animeIDList []types.AnimeID) ([]*Anime, error) {
	errg, ctx := errgroup.WithContext(context.Background())

	results := make(chan types.Pair[int, *Anime])

	for idx, id := range animeIDList {
		idx := idx
		id := id
		errg.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				animeAns, err := GetAnimeByID(id)
				if err != nil {
					return err
				}

				select {
				case results <- types.Pair[int, *Anime]{First: idx, Second: animeAns}:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	go func() {
		errg.Wait()
		close(results)
	}()

	ans := make([]*Anime, len(animeIDList))

	for result := range results {
		ans[result.First] = result.Second
	}

	err := errg.Wait()

	return ans, err
}
