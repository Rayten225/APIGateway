package handlers

import (
	"micronews/news-service/store"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const pageSize = 10

func init() {
	store.Init()
}

func ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	var total int
	store.DB.QueryRow("SELECT count(*) FROM news").Scan(&total)
	pages := (total + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize
	rows, err := store.DB.Query(
		"SELECT id, title, content, published FROM news ORDER BY published DESC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var list []store.News
	for rows.Next() {
		var n store.News
		rows.Scan(&n.ID, &n.Title, &n.Content, &n.Published)
		list = append(list, n)
	}
	c.JSON(http.StatusOK, gin.H{
		"items":      list,
		"pagination": gin.H{"total_pages": pages, "current_page": page, "per_page": pageSize},
	})
}

func GetNews(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var n store.News
	err := store.DB.QueryRow(
		"SELECT id, title, content, published FROM news WHERE id = $1",
		id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.Published)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		return
	}
	c.JSON(http.StatusOK, n)
}

func CreateNews(c *gin.Context) {
	var news store.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	news.Published = time.Now()
	err := store.DB.QueryRow(
		"INSERT INTO news (title, content, published) VALUES ($1, $2, $3) RETURNING id",
		news.Title, news.Content, news.Published,
	).Scan(&news.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}
