package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type movie struct {
	EntryNo       int64   `json:"entryNo"`
	MovieName     string  `json:"movieName"`
	ReviewScore   float64 `json:"reviewScore"`
	Synopsis      string  `json:"synopsis"`
	Director      string  `json:"Director"`
	YearOfRelease int64   `json:"yearOfRelease"`
}

var movies = []movie{
	{EntryNo: 1, MovieName: "Robo", ReviewScore: 8.5, Synopsis: "A film about Robos", Director: "Shankar", YearOfRelease: 2012},
	{EntryNo: 2, MovieName: "Robo 2", ReviewScore: 7.5, Synopsis: "A film about Robos - Part 2", Director: "Shankar", YearOfRelease: 2018},
	{EntryNo: 3, MovieName: "RRR", ReviewScore: 9.5, Synopsis: "A fiction film about two freedom fighters", Director: "Rajamouli", YearOfRelease: 2022},
}

func getAllMoviesData(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, movies)
}

func addMovieRecord(context *gin.Context) {
	var newMovie movie
	if err := context.BindJSON(&newMovie); err != nil {
		return
	}
	movies = append(movies, newMovie)
	context.IndentedJSON(http.StatusCreated, newMovie)
}

func getMovieRecordByYear(yearOfRelease int64) (*[]movie, error) {
	var filteredMovies = []movie{}
	for i := 0; i < len(movies); i++ {
		if movies[i].YearOfRelease == yearOfRelease {
			filteredMovies = append(filteredMovies, movies[i])
		}
	}
	if len(filteredMovies) != 0 {
		return &filteredMovies, nil
	}
	errorMsg := "No Movie Data Found in the Year " + string(yearOfRelease)
	return nil, errors.New(errorMsg)
}

func getByYearOfRelease(context *gin.Context) {
	yearOfRelease, _ := strconv.Atoi(context.Param("yearOfRelease"))
	filteredMovies, err := getMovieRecordByYear(int64(yearOfRelease))

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "No Movies Found"})
		return
	}
	context.IndentedJSON(http.StatusOK, filteredMovies)

}

func getMovieRecordByDirector(director string) (*[]movie, error) {
	var filteredMovies = []movie{}
	for i := 0; i < len(movies); i++ {
		if strings.Contains(strings.ToLower(movies[i].Director), strings.ToLower(director)) {
			filteredMovies = append(filteredMovies, movies[i])
		}
	}
	if len(filteredMovies) != 0 {
		return &filteredMovies, nil
	}
	return nil, errors.New("No movie data found with the provided director name")
}

func getByDirectorName(context *gin.Context) {
	directorName := context.Param("director")
	filteredMovies, err := getMovieRecordByDirector(directorName)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "No Movies Found"})
		return
	}
	context.IndentedJSON(http.StatusOK, filteredMovies)

}

func deleteMovieByEntryNo(entryNo int64) (string, error) {
	flag := false
	for i := 0; i < len(movies); i++ {
		if movies[i].EntryNo == entryNo {
			movies[i] = movies[len(movies)-1]
			movies = movies[:len(movies)-1]
			flag = true
		}
	}
	if flag {
		return "Movie record Deleted Successfully", nil
	}
	return "", errors.New("No movie data found with the provided director name")
}

func deleteByEntryNo(context *gin.Context) {
	entryNo, _ := strconv.Atoi(context.Param("entryNo"))
	msg, err := deleteMovieByEntryNo(int64(entryNo))

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "No Movies Found to Delete"})
		return
	}
	context.IndentedJSON(http.StatusOK, msg)

}

func main() {
	router := gin.Default()
	router.GET("/getAllMovies", getAllMoviesData)
	router.POST("/addMovie", addMovieRecord)
	router.GET("/getMovieByYearOfRelease/:yearOfRelease", getByYearOfRelease)
	router.GET("/getMovieByDirectorName/:director", getByDirectorName)
	router.DELETE("/deleteMovieByEntryNo/:entryNo", deleteByEntryNo)
	router.Run("localhost:8082")
}
