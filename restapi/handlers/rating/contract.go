package rating

import "github.com/WikimeCorp/WikimeBackend/types"

type SetRatingRequest struct {
	Rating *types.AnimeRating `json:"rating" validate:"required,gte=1,lte=5"`
}
