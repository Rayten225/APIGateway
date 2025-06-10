package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"micronews/comment-service/store"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	store.Init()
}

func CreateComment(c *gin.Context) {
	var req store.Comment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comment text cannot be empty"})
		return
	}

	newsResp, err := http.Get(fmt.Sprintf("http://news-service:8005/news/%d", req.NewsID))
	if err != nil || newsResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "news item not found"})
		return
	}

	censorReqBody, _ := json.Marshal(map[string]string{"text": req.Text})
	censorReq, err := http.NewRequest("POST", "http://censor-service:8003/censor", bytes.NewBuffer(censorReqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create censor request"})
		return
	}
	censorReq.Header.Set("Content-Type", "application/json")
	censorReq.Header.Set("X-Request-ID", c.GetString("request_id"))
	client := &http.Client{}
	censorResp, err := client.Do(censorReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "censor service unavailable"})
		return
	}
	defer censorResp.Body.Close()

	if censorResp.StatusCode != http.StatusOK {
		var errorResp map[string]string
		if err := json.NewDecoder(censorResp.Body).Decode(&errorResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse censor error"})
			return
		}
		c.JSON(censorResp.StatusCode, gin.H{"error": errorResp["error"]})
		return
	}

	var id int
	var created time.Time
	err = store.DB.QueryRow(
		"INSERT INTO comments (news_id, parent_id, text, created) VALUES ($1, $2, $3, NOW()) RETURNING id, created",
		req.NewsID, req.ParentID, req.Text,
	).Scan(&id, &created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	req.Created = created
	c.JSON(http.StatusOK, req)
}

func ListComments(c *gin.Context) {
	newsID, _ := strconv.Atoi(c.Query("news_id"))
	rows, err := store.DB.Query(
		"SELECT id, news_id, parent_id, text, created FROM comments WHERE news_id = $1 ORDER BY created DESC",
		newsID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var list []store.Comment
	for rows.Next() {
		var cmt store.Comment
		if err := rows.Scan(&cmt.ID, &cmt.NewsID, &cmt.ParentID, &cmt.Text, &cmt.Created); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		list = append(list, cmt)
	}
	c.JSON(http.StatusOK, list)
}
