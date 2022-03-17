package assetsapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
	"github.com/praveennagaraj97/shoppers-gocommerce/models"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/serialize"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/color"
	logger "github.com/praveennagaraj97/shoppers-gocommerce/pkg/log"
	assetrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/asset"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetsAPI struct {
	repo *assetrepository.AssetRepository
	conf *conf.GlobalConfiguration
}

type blueReturn struct {
	s3  *s3.PutObjectOutput
	err error
}

func (a *AssetsAPI) Initialize(rp *assetrepository.AssetRepository, cfg *conf.GlobalConfiguration) {
	a.conf = cfg
	a.repo = rp
}

func (a *AssetsAPI) UploadSingleAsset() gin.HandlerFunc {
	return func(c *gin.Context) {

		var payload dto.NewAssetDTO
		var ch chan blueReturn

		err := c.ShouldBind(&payload)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		file, err := c.FormFile(payload.FileFieldName)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		fileType := file.Header.Get("Content-Type")

		if payload.BlurDataRequired && !strings.Contains(fileType, "image") {
			api.SendErrorResponse(a.conf.Localize, c, "blur_data_can_only_be_created_for_images", http.StatusUnprocessableEntity, nil)
			return
		}

		// Read the file buffer.
		multiPartFile, err := file.Open()
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, "something_went_wrong", http.StatusBadRequest, nil)
			return
		}
		buffer, err := io.ReadAll(multiPartFile)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, "something_went_wrong", http.StatusBadRequest, nil)
			return
		}

		defer multiPartFile.Close()

		fileId := primitive.NewObjectID()
		fileExtensionClips := strings.Split(file.Filename, ".")
		fileExtension := fileExtensionClips[len(fileExtensionClips)-1]
		fileName := fmt.Sprintf("%s.%s", fileId.Hex(), fileExtension)

		originalFilePath := fmt.Sprintf("%s/original", payload.ContainerName)

		// Upload blur
		if payload.BlurDataRequired {
			ch = make(chan blueReturn, 1)
			go a.saveBlurImage(&buffer, originalFilePath, fileName, &fileType, &ch)
		}

		// Original
		_, err = a.conf.AWSUtils.UploadAsset(bytes.NewBuffer(buffer), originalFilePath, fileName, &fileType)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c,
				"failed_to_upload",
				http.StatusInternalServerError,
				map[string]string{"reason": err.Error()})
			return
		}

		select {
		case value, ok := <-ch:
			if ok {
				if value.err != nil {
					api.SendErrorResponse(a.conf.Localize, c,
						"failed_to_upload",
						http.StatusInternalServerError,
						map[string]string{"reason": value.err.Error()})
					return
				}
			}
		default:
		}

		// Save to database
		var lowQualityUrl string
		if payload.BlurDataRequired {
			blurPath := strings.Replace(fileName, "original", "blur", 1)
			lowQualityUrl = fmt.Sprintf("%s/%s/%s", a.conf.AWSUtils.S3PUBLIC_DOMAIN, blurPath, fileName)
		}

		dataModel := &models.AssetModel{
			ID:          primitive.NewObjectID(),
			OriginalURL: fmt.Sprintf("%s/%s/%s", a.conf.AWSUtils.S3PUBLIC_DOMAIN, fileName, fileName),
			BlurDataURL: lowQualityUrl,
			ContentType: fileType,
			MetaData: &models.FileMetaModel{
				Key:         fileName + "/" + fileName,
				Title:       payload.Title,
				Description: payload.Description,
			},
			CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
			Published:    payload.Published,
			LinkedEntity: make([]*models.AssetLinkedWithInfo, 0),
		}

		res, err := a.repo.AddNewAsset(dataModel)

		if err != nil {
			// Delete assets if db failes to save.
			_, uploadErr := a.deleteAsset(originalFilePath+"/"+fileName, payload.BlurDataRequired)
			if uploadErr != nil {
				logger.PrintLog("Failed to delete asset", color.Red)
			}

			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusCreated, &serialize.DataResponse{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    a.conf.Localize.GetMessage("asset_uploaded_successfully", c),
			},
		})
	}
}

func (a *AssetsAPI) GetAllAssets() gin.HandlerFunc {
	return func(c *gin.Context) {
		paginateOptions := api.ParsePaginationOptions(c)
		sortOptions := api.ParseSortByOptions(c)
		filterOptions := api.ParseFilterByOptions(c)
		keySetSortby := "$gt"

		// Default options | sort by latest
		if len(*sortOptions) == 0 {
			sortOptions = &map[string]int8{"created_at": -1}
		}
		// Key Set fix for created_at desc
		if paginateOptions.PaginateId != nil {
			for key, value := range *sortOptions {
				if value == -1 && key == "created_at" {
					keySetSortby = "$lt"
				}
			}
		}

		res, err := a.repo.FindAll(paginateOptions, sortOptions, filterOptions, keySetSortby)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID
		if paginateOptions.PaginateId == nil {
			docCount, err = a.repo.GetDocumentsCount(filterOptions)
			if err != nil {
				api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, paginateOptions, int64(resLen), lastResId)

		c.JSON(http.StatusOK, &serialize.PaginatedDataResponse{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			DataResponse: serialize.DataResponse{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "list of assets retrieved",
				},
			},
		})
	}
}

func (a *AssetsAPI) GetAssetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objectID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.repo.FindByID(&objectID)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "asset_retrieved_successfully",
			},
		})

	}
}
