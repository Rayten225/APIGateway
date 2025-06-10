package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func CensorText(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text cannot be empty"})
		return
	}
	censorReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s/censor", getEnv("CENSOR_SERVICE_ADDR", "censor-service:8003")), bytes.NewReader(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
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
	body, err := io.ReadAll(censorResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}
	c.Data(censorResp.StatusCode, "application/json", body)
}

func CreateNews(c *gin.Context) {
	var news struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if news.Title == "" || news.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and content cannot be empty"})
		return
	}

	censorData, _ := json.Marshal(map[string]string{"text": news.Content})
	censorReq, err := http.NewRequest("POST", "http://censor-service:8003/censor", bytes.NewBuffer(censorData))
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

	newsData, _ := json.Marshal(news)
	newsReq, err := http.NewRequest("POST", "http://news-service:8005/news", bytes.NewBuffer(newsData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create news request"})
		return
	}
	newsReq.Header.Set("Content-Type", "application/json")
	newsReq.Header.Set("X-Request-ID", c.GetString("request_id"))
	newsResp, err := client.Do(newsReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "news service unavailable"})
		return
	}
	defer newsResp.Body.Close()

	body, err := io.ReadAll(newsResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}
	c.Data(newsResp.StatusCode, "application/json", body)
}

func ListNews(c *gin.Context) {
	req, err := http.NewRequest("GET", "http://news-service:8005/news", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	query := c.Request.URL.RawQuery // Передаем query-параметры
	req.URL.RawQuery = query
	req.Header.Set("X-Request-ID", c.GetString("request_id"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "news service unavailable"})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}
	c.Data(resp.StatusCode, "application/json", body)
}

func GetNewsDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing news id"})
		return
	}

	// Запрос к News Service
	newsReq, err := http.NewRequest("GET", fmt.Sprintf("http://news-service:8005/news/%s", id), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create news request"})
		return
	}
	newsReq.Header.Set("X-Request-ID", c.GetString("request_id"))
	client := &http.Client{}
	newsResp, err := client.Do(newsReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "news service unavailable"})
		return
	}
	defer newsResp.Body.Close()
	newsBody, err := io.ReadAll(newsResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read news response"})
		return
	}
	if newsResp.StatusCode != http.StatusOK {
		c.JSON(newsResp.StatusCode, gin.H{"error": string(newsBody)})
		return
	}

	// Запрос к Comment Service
	commentReq, err := http.NewRequest("GET", fmt.Sprintf("http://comment-service:8004/comments?news_id=%s", id), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment request"})
		return
	}
	commentReq.Header.Set("X-Request-ID", c.GetString("request_id"))
	commentResp, err := client.Do(commentReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "comment service unavailable"})
		return
	}
	defer commentResp.Body.Close()
	commentBody, err := io.ReadAll(commentResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read comment response"})
		return
	}
	if commentResp.StatusCode != http.StatusOK {
		c.JSON(commentResp.StatusCode, gin.H{"error": string(commentBody)})
		return
	}

	// Формируем ответ
	var newsResponse interface{}
	if err := json.Unmarshal(newsBody, &newsResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse news response"})
		return
	}
	var comments []interface{}
	if err := json.Unmarshal(commentBody, &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse comment response"})
		return
	}

	response := map[string]interface{}{
		"news":     newsResponse,
		"comments": comments,
	}
	if len(comments) == 0 {
		response["comments"] = nil
	}
	c.JSON(http.StatusOK, response)
}

func CreateComment(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	var comment struct {
		NewsID int    `json:"news_id"`
		Text   string `json:"text"`
	}
	if err := json.Unmarshal(data, &comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if comment.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comment text cannot be empty"})
		return
	}
	// Проверка цензуры
	censorReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s/censor", getEnv("CENSOR_SERVICE_ADDR", "censor-service:8003")), bytes.NewReader(data))
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
	// Создание комментария
	req, err := http.NewRequest("POST", "http://comment-service:8004/comments", bytes.NewReader(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", c.GetString("request_id"))
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "comment service unavailable"})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}
	c.Data(resp.StatusCode, "application/json", body)
}
