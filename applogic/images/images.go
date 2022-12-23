package images

import (
	"io"
	"os"
	"path"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/WikimeCorp/WikimeBackend/utils"
)

func saveFile(imagePath string, image io.Reader) error {
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, image)
	if err != nil {
		return err
	}

	return nil
}

// AddImageToAnime adding image to anime
func AddImageToAnime(animeID types.AnimeID, image io.Reader, imageTail string) error {
	filename := utils.FastRandomString(40) + "." + imageTail
	imagePath := path.Join(config.Config.ImagePathDisk, config.Config.ImagesPathURI, filename)

	err := saveFile(imagePath, image)
	if err != nil {
		return err
	}

	imagePathURI := path.Join("/", config.Config.ImagesPathURI, filename)

	err = db.AddImageToAnime(animeID, imagePathURI)
	if err != nil {
		err1 := os.Remove(imagePath)
		if err1 != nil {
			return err1
		}
		return err
	}

	err = db.UpdateDataAdded(animeID)

	return err
}

func SetPoster(animeID types.AnimeID, image io.Reader, imageTail string) error {
	filename := utils.FastRandomString(40) + "." + imageTail
	imagePath := path.Join(config.Config.ImagePathDisk, config.Config.ImagesPathURI, filename)

	err := saveFile(imagePath, image)
	if err != nil {
		return err
	}

	imagePathURI := path.Join("/", config.Config.ImagesPathURI, filename)

	curAnime, err := anime.GetAnimeByID(animeID)
	if err != nil {
		return err
	}

	oldImage := curAnime.Poster

	err = db.SetPoster(animeID, imagePathURI)
	if err != nil {
		err = os.Remove(imagePath)
		if err != nil {
		}
		return err
	}

	if *oldImage != path.Join("/", config.Config.ImagesPathURI, config.Config.DefaultAnimePosterPath) {
		err = os.Remove(path.Join(config.Config.ImagePathDisk, *oldImage))
		if err != nil {
			return err
		}
	}

	err = db.UpdateDataAdded(animeID)

	return err
}

func SetAvatar(userID types.UserID, image io.Reader, imageTail string) error {
	filename := utils.FastRandomString(40) + "." + imageTail
	imagePath := path.Join(config.Config.ImagePathDisk, config.Config.ImagesPathURI, filename)

	err := saveFile(imagePath, image)
	if err != nil {
		return err
	}

	imagePathURI := path.Join("/", config.Config.ImagesPathURI, filename)
	userObj, err := db.GetUser(userID)
	if err != nil {
		return err
	}

	oldImage := userObj.Avatar

	err = db.SetAvatar(userID, imagePathURI)
	if err != nil {
		err1 := os.Remove(imagePath)
		if err1 != nil {
			return err1
		}
		return err
	}

	if oldImage != path.Join("/", config.Config.ImagesPathURI, config.Config.DefaultUserAvatarPath) {
		err = os.Remove(path.Join(config.Config.ImagePathDisk, oldImage))
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteImageFromAnime(animeID types.AnimeID, img string) error {
	if img == config.Config.DefaultAnimePosterPath {
		return myerrors.ErrForbiden
	}
	err := db.DeleteImageFromAnime(animeID, path.Join(config.Config.ImagesPathURI, img))
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(config.Config.ImagePathDisk, config.Config.ImagesPathURI, img))
	if err != nil {
		return err
	}
	err = db.UpdateDataAdded(animeID)
	return err

}
