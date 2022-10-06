package db

import (
	. "github.com/WikimeCorp/WikimeBackend/types"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
)

func createCommentsDoc(animeID AnimeID) error {
	anime, err := CheckAnime(animeID)
	if err != nil {
		return err
	}
	if !anime {
		return &inerr.ErrAnimeNotFound{animeID}
	}
	_, err = commentsCollection.InsertOne(ctx, Comments{ID: animeID})
	return err
}
