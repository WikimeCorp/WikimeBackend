package anime

import (
	"context"

	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
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
		Rating:      Rating(*anime.Rating),
		Poster:      anime.Poster,
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

func SetAverage(anime types.AnimeID, average float64) error {
	err := db.SetAverage(anime, average)
	return err
}

func GetAnimeSortedByRating(genres []string, order int8) ([]types.AnimeID, error) {
	ans, badGenres := db.CheckGenres(genres)
	if ans == false {
		return nil, &myerrors.ErrWrongGenres{badGenres}
	}

	animeIDs, err := db.GetAnimeIDsSortedByRating(genres, order)
	if err != nil {
		return nil, err
	}

	return animeIDs, nil
}

func GetAnimeSortedByFavorites(genres []string, order int8) ([]types.AnimeID, error) {
	ans, badGenres := db.CheckGenres(genres)
	if ans == false {
		return nil, &myerrors.ErrWrongGenres{badGenres}
	}

	animeIDs, err := db.GetAnimeIDsSortedByFavorites(genres, order)
	if err != nil {
		return nil, err
	}

	return animeIDs, nil
}

func GetAnimeSortedByAddingDate(genres []string, order int8) ([]types.AnimeID, error) {
	ans, badGenres := db.CheckGenres(genres)
	if ans == false {
		return nil, &myerrors.ErrWrongGenres{badGenres}
	}

	animeIDs, err := db.GetAnimeIDsSortedByAddingDate(genres, order)
	if err != nil {
		return nil, err
	}

	return animeIDs, nil
}

func GetAnimeSortedByReleaseDate(genres []string, order int8) ([]types.AnimeID, error) {
	ans, badGenres := db.CheckGenres(genres)
	if ans == false {
		return nil, &myerrors.ErrWrongGenres{badGenres}
	}

	animeIDs, err := db.GetAnimeIDsSortedByReleaseDate(genres, order)
	if err != nil {
		return nil, err
	}

	return animeIDs, nil
}

func GetMostPopular(limit int) ([]types.AnimeID, error) {
	return db.GetAnimeIDsSortedByFavorites([]string{}, -1, limit)
}

// CheckAnime returns whether the anime exists. Can be true or false
func CheckAnime(animeID types.AnimeID) (bool, error) {
	ans, err := db.CheckAnime(animeID)
	return ans, err
}
