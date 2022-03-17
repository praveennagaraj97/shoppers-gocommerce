package api

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaginationOptions struct {
	// Cached Count - Utilized by Key Set Pagination
	CachedCount uint64
	// Cached Page Number - Utilized by Key Set Pagination
	CachedPageNum uint64
	// Results per page
	PerPage int
	// Page number count
	PageNum int
	// Key set page id
	PaginateId *primitive.ObjectID
}

func GetPaginateOptions(docCount int64, pgOpts *PaginationOptions, docLen int64, lastResID *primitive.ObjectID) (*uint64, *bool, *bool, *string) {
	// Paginate Options
	var count uint64
	var next bool
	var prev bool
	var paginateId *string
	var pageNum uint64

	if pgOpts.PaginateId == nil {

		count = uint64(docCount)
		next = pgOpts.PageNum < int(count)/pgOpts.PerPage || count > uint64(pgOpts.PageNum*pgOpts.PerPage)
		prev = pgOpts.PageNum > 1

		// First Next paginate ID
		paginateObjectId := lastResID
		paginateId = encodeKeySetPaginationID(count, paginateObjectId, 1)
	} else {
		count = pgOpts.CachedCount
		pageNum = pgOpts.CachedPageNum + 1
		if pageNum < pgOpts.CachedCount/uint64(pgOpts.PerPage) || count > uint64(int(pageNum)*pgOpts.PerPage) {
			paginateObjectId := lastResID
			next = true
			paginateId = encodeKeySetPaginationID(count, paginateObjectId, int64(pageNum))
		}
	}

	return &count, &next, &prev, paginateId

}

// parse pagination options from request URL. accepts per_page & page_num & paginate_id(for infinite scroll and performance).
func ParsePaginationOptions(c *gin.Context) *PaginationOptions {
	// Check for startId
	count, nextID, cachedPageNum, err := decodeKeySetPaginationID(c.Request.URL.Query().Get("paginate_id"))
	if err != nil {
		count = 0
		nextID = nil
	}

	// read the page num if key set pagination is not requested.
	pageNum, err := strconv.Atoi(c.Request.URL.Query().Get("page_num"))
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	// Results per page
	perPage, err := strconv.Atoi(c.Request.URL.Query().Get("per_page"))
	if err != nil || perPage == 0 {
		perPage = constants.DefaultPerPageResults
	}

	return &PaginationOptions{
		PerPage:       perPage,
		PageNum:       pageNum,
		PaginateId:    nextID,
		CachedCount:   count,
		CachedPageNum: cachedPageNum,
	}
}

func encodeKeySetPaginationID(count uint64, paginateId *primitive.ObjectID, pageNum int64) *string {
	if paginateId == nil {
		return nil
	}
	paginateString := fmt.Sprintf("nextId=%s&count=%d&pageNum=%d", paginateId.Hex(), count, pageNum)
	encryptedID := base64.StdEncoding.EncodeToString([]byte(paginateString))

	return &encryptedID
}

func decodeKeySetPaginationID(encryptedID string) (uint64, *primitive.ObjectID, uint64, error) {

	decodedID, err := base64.StdEncoding.DecodeString(encryptedID)
	if err != nil {
		return 0, nil, 1, err
	}

	info, err := url.ParseQuery(string(decodedID))
	if err != nil {
		return 0, nil, 1, err
	}

	count, err := strconv.ParseUint(info.Get("count"), 10, 64)
	if err != nil {
		return 0, nil, 1, err
	}

	nextId, err := primitive.ObjectIDFromHex(info.Get("nextId"))
	if err != nil {
		return 0, nil, 1, err
	}

	pageNum, err := strconv.ParseUint(info.Get("pageNum"), 10, 64)
	if err != nil {
		return 0, nil, 1, err
	}

	return count, &nextId, pageNum, nil

}
