package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getWordEntry(c *gin.Context){
	word := c.Param("word")
	entry := doOnlineGoogleDictionary(word)
	c.IndentedJSON(http.StatusOK, entry)
}