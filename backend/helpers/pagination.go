package helpers

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-commerce/structs"

)

type PaginationLink struct {
	URL    string `json:"url"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil || i < 1 {
		return 1
	}
	return i
}

func TotalPage(total int64, perPage int) int {
	if perPage == 0 {
		return 1
	}
	pages := int(total) / perPage
	if int(total)%perPage != 0 {
		pages++
	}
	return pages
}

func BuildPaginationLinks(currentPage, lastPage int, baseURL, search string) []PaginationLink {
	links := []PaginationLink{}

	var prevURL string
	if currentPage > 1 {
		prevURL = PageURL(baseURL, currentPage-1, lastPage, search)
	}

	links = append(links, PaginationLink{
		URL:    prevURL,
		Label:  "&laquo; Previous",
		Active: false,
	})

	for i := 1; i <= lastPage; i++ {
		links = append(links, PaginationLink{
			URL:    baseURL + "?page=" + strconv.Itoa(i) + QueryString(search),
			Label:  strconv.Itoa(i),
			Active: i == currentPage,
		})
	}

	var nextURL string
	if currentPage < lastPage {
		nextURL = PageURL(baseURL, currentPage+1, lastPage, search)
	}

	links = append(links, PaginationLink{
		URL:    nextURL,
		Label:  "Next &raquo;",
		Active: false,
	})

	return links
}

func PageURL(baseURL string, page, lastPage int, search string) string {
	if page < 1 || page > lastPage {
		return ""
	}
	return baseURL + "?page=" + strconv.Itoa(page) + QueryString(search)
}

func QueryString(search string) string {
	if search == "" {
		return ""
	}
	return "&search=" + search
}

func GetPaginationParams(c *gin.Context) (search string, page, limit, offset int) {
	search = c.Query("search")
	page = StringToInt(c.DefaultQuery("page", "1"))
	limit = StringToInt(c.DefaultQuery("limit", "10"))
	offset = (page - 1) * limit
	return
}

func BuildBaseURL(c *gin.Context) string {
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return scheme + "://" + c.Request.Host + c.Request.URL.Path
}

func BuildHostURL(c *gin.Context) string {
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return scheme + "://" + c.Request.Host
}

func PaginateResponse(c *gin.Context, data interface{}, total int64, page, limit int, baseURL, search, message string) {
	lastPage := TotalPage(total, limit)

	dataLen := reflect.ValueOf(data).Len()
	var from, to int
	if dataLen > 0 {
		from = (page-1)*limit + 1
		to = from + dataLen - 1
	} else {
		from = 0
		to = 0
	}

	links := BuildPaginationLinks(page, lastPage, baseURL, search)

	var prevPageURL, nextPageURL string
	if page > 1 {
		prevPageURL = PageURL(baseURL, page-1, lastPage, search)
	}
	if page < lastPage {
		nextPageURL = PageURL(baseURL, page+1, lastPage, search)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data: gin.H{
			"current_page":   page,
			"data":           data,
			"first_page_url": baseURL + "?page=1" + QueryString(search),
			"from":           from,
			"last_page":      lastPage,
			"last_page_url":  baseURL + "?page=" + strconv.Itoa(lastPage) + QueryString(search),
			"links":          links,
			"next_page_url":  nextPageURL,
			"path":           baseURL,
			"per_page":       limit,
			"prev_page_url":  prevPageURL,
			"to":             to,
			"total":          total,
		},
	})
}
