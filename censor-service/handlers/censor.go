package handlers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

var forbidden = []string{"qwerty", "йцукен", "zxvbnm"}

func Censor(c *gin.Context) {
    var body struct{ Text string `json:"text"` }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
        return
    }
    low := strings.ToLower(body.Text)
    for _, w := range forbidden {
        if strings.Contains(low, w) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "forbidden content"})
            return
        }
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
