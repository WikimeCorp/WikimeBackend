package applogic

import (
	"github.com/WikimeCorp/WikimeBackend/db"
	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
)

func AddAnime(anime *dbtypes.Anime) (AnimeID, error) {
	return db.AddAnime(anime)
}
