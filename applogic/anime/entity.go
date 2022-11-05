package anime

import (
	"time"

	"github.com/WikimeCorp/WikimeBackend/types"
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
