package avatar

import (
	"bytes"
	"fmt"
	"image"

	"github.com/czhj/ahfs/modules/setting"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

const AvatarSize = 256

func Prepare(data []byte) (*image.Image, error) {
	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Cannot decode image config: %v", err)
	}

	if imgCfg.Width > setting.AvatarMaxWidth {
		return nil, fmt.Errorf("Image width too large: %d > %d", imgCfg.Width, setting.AvatarMaxWidth)
	}

	if imgCfg.Height > setting.AvatarMaxHeight {
		return nil, fmt.Errorf("Image height too large: %d > %d", imgCfg.Height, setting.AvatarMaxHeight)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Cannot decode image: %v", err)
	}

	if imgCfg.Width != imgCfg.Height {
		var ax, ay, newSize int
		if imgCfg.Width > imgCfg.Height {
			newSize = imgCfg.Height
			ax = (imgCfg.Width - imgCfg.Height) / 2
		} else {
			newSize = imgCfg.Width
			ay = (imgCfg.Height - imgCfg.Width) / 2
		}

		img, err = cutter.Crop(img, cutter.Config{
			Width:  newSize,
			Height: newSize,
			Anchor: image.Point{ax, ay},
		})
		if err != nil {
			return nil, err
		}
	}

	img = resize.Resize(AvatarSize, AvatarSize, img, resize.NearestNeighbor)
	return &img, nil
}
