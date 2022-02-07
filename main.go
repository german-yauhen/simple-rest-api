package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", addAlbum)
	router.GET("/albums/:id", getAlbum)
	router.DELETE("albums/:id", deleteAlbum)

	router.Run("localhost:8085")
}

type album struct {
	ID     uint8   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func getAlbums(cont *gin.Context) {
	cont.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(cont *gin.Context) {
	var newAlbum album
	if error := cont.BindJSON(&newAlbum); error != nil {
		return
	}
	newAlbum.ID = getNewId()
	albums = append(albums, newAlbum)
	cont.IndentedJSON(http.StatusCreated, albums)
}

func getAlbum(cont *gin.Context) {
	id := getIdFromRequestParam(cont)
	foundAlbum := getAlbumById(id)
	if (album{}) == foundAlbum {
		cont.IndentedJSON(http.StatusNotFound, gin.H{"message": "illegal id provided"})		
	}
	cont.IndentedJSON(http.StatusOK, foundAlbum)
}

func deleteAlbum(cont *gin.Context) {
	id := getIdFromRequestParam(cont)
	albumToDeleteIndex := 0
	for index, album := range albums {
		if album.ID == id {
			albumToDeleteIndex = index
		}
	}
	if albumToDeleteIndex == 0 {
		cont.IndentedJSON(http.StatusBadRequest, gin.H{"message": "illegal id provided"})
	}
	albums = append(albums[:albumToDeleteIndex], albums[albumToDeleteIndex+1:]...)
	if (album{}) != getAlbumById(id) {
		cont.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an internal error occured"})
	}
	cont.IndentedJSON(http.StatusNoContent, nil)
}

func getAlbumById(id uint8) album {
	for _, album := range albums {
		if album.ID == id {
			return album
		}
	}
	return album{}
}

func getIdFromRequestParam(cont *gin.Context) uint8 {
	parsedId, error := strconv.ParseUint(cont.Param("id"), 10, 64)
	if error != nil {
		cont.IndentedJSON(http.StatusBadRequest, gin.H{"message": "illegal id provided"})
	}
	id := uint8(parsedId)
	return id
}

func getNewId() uint8 {
	latestAlbum := albums[0]
	for _, album := range albums {
		if album.ID > latestAlbum.ID {
			latestAlbum = album
		}
	}
	return latestAlbum.ID + 1
}

var albums = []album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	{ID: 4, Title: "Sinatra and Strings", Artist: "Frank Sinatra", Price: 45.55},
}
