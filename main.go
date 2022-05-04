package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

func getMovies(c *fiber.Ctx) error {
	c.JSON(movies)
	return nil
}

func deleteMovie(c *fiber.Ctx) error {
	for index, item := range movies {
		if item.ID == c.Params("id") {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	c.JSON(movies)
	return nil
}

func getMovie(c *fiber.Ctx) error {
	for _, item := range movies {
		if item.ID == c.Params("id") {
			c.JSON(item)
			return nil
		}
	}
	return fiber.ErrBadRequest
}

func createMovie(c *fiber.Ctx) error {
	var movie Movie
	c.BodyParser(&movie)
	movie.ID = strconv.Itoa(rand.Int())
	movies = append(movies, movie)
	c.Status(http.StatusCreated).JSON(movie)
	return nil
}

func updateMovie(c *fiber.Ctx) error {
	for index, movie := range movies {
		if movie.ID == c.Params("id") {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			c.BodyParser(&movie)
			movie.ID = c.Params("id")
			movies = append(movies, movie)
			c.Status(http.StatusAccepted).JSON(movie)
			return nil
		}
	}
	return fiber.ErrBadRequest
}

func main() {
	injectDummyData()

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/movies", getMovies)
	app.Get("/movies/:id", getMovie)
	app.Post("/movies", createMovie)
	app.Put("/movies/:id", updateMovie)
	app.Delete("/movies/:id", deleteMovie)

	port := 8000
	fmt.Printf("Starting server at port %d\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
