package assetsapi

import (
	"bytes"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/praveennagaraj97/shopee/pkg/utils"
)

func (a *AssetsAPI) saveBlurImage(buffer *[]byte, originalFilePath string, fileName string, fileType *string, ch *chan blueReturn) {
	blurPath := strings.Replace(originalFilePath, "original", "blur", 1)

	blurredBuffer, err := utils.CreateBlurDataForImages(*buffer, 1, 20, 20)
	if err != nil {
		*ch <- blueReturn{s3: nil, err: err}
		return
	}

	res, err := a.conf.AWSUtils.UploadAsset(bytes.NewBuffer(blurredBuffer), blurPath, fileName, fileType)

	if err != nil {
		*ch <- blueReturn{s3: nil, err: err}
		return
	}

	*ch <- blueReturn{s3: res, err: nil}

}

func (a *AssetsAPI) deleteAsset(key string, blueExist bool) (*s3.DeleteObjectOutput, error) {

	blurKey := strings.Replace(key, "original", "blur", 1)

	if blueExist {
		_, _ = a.conf.AWSUtils.DeleteAsset(&blurKey)
	}

	return a.conf.AWSUtils.DeleteAsset(&key)
}
