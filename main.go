package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func injectDummyData() {
	movies = append(movies,
		Movie{ID: "1", Isbn: "45423", Title: "first", Director: &Director{Firstname: "John", Lastname: "Doe"}},
		Movie{ID: "2", Isbn: "45425", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}},
	)
}

func getMovies(c *gin.Context) {
	c.JSON(200, movies)
}

func deleteMovie(c *gin.Context) {
	for index, item := range movies {
		if item.ID == c.Param("id") {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	c.JSON(200, movies)
}

func getMovie(c *gin.Context) {
	for _, item := range movies {
		if item.ID == c.Param("id") {
			c.JSON(http.StatusOK, item)
			break
		}
	}
}

func createMovie(c *gin.Context) {
	var movie Movie
	c.BindJSON(&movie)
	movie.ID = strconv.Itoa(rand.Int())
	movies = append(movies, movie)
	c.JSON(http.StatusCreated, movie)
}

func updateMovie(c *gin.Context) {
	for index, movie := range movies {
		if movie.ID == c.Param("id") {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			c.BindJSON(&movie)
			movie.ID = c.Param("id")
			movies = append(movies, movie)
			c.JSON(200, movie)
			return
		}
	}
}

func main() {
	injectDummyData()

	router := gin.Default()

	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovie)
	router.POST("/movies", createMovie)
	router.PUT("/movies/:id", updateMovie)
	router.DELETE("/movies/:id", deleteMovie)

	port := 8000
	fmt.Printf("Starting server at port %d\n", port)
	log.Fatal(router.Run(fmt.Sprintf(":%d", port)))
}
