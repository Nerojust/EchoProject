package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {
	fmt.Println("Starting")
	//start echo
	e := echo.New()
	//for grouping routes to one category.
	g := e.Group("/admin")
	//using middleware
	//to format the logs to view on the console
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          nil,
		Format:           `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
		CustomTimeFormat: "",
		Output:           nil,
	}))
	g.Use(middleware.CORS())

	g.GET("/main", mainAdmin)

	e.GET("/", serverStart)
	e.GET("/cats", getCats)
	e.GET("/cats/:data", getDataType)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)

	_ = e.Start(":8000")
}

func mainAdmin(context echo.Context) error {
	return context.String(http.StatusOK, "Secret admin")
}

func addHamster(context echo.Context) error {
	hamster := Hamster{}
	err := context.Bind(&hamster)
	if err != nil {
		log.Printf("failed processing add hamster request %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, )
	}
	log.Printf("this is ur hamster details %#v", hamster)
	return context.String(http.StatusOK, "Hamster request successfully sent")
}

func addDog(context echo.Context) error {
	dog := Dog{}
	defer context.Request().Body.Close()
	err := json.NewDecoder(context.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("failed processing add dog request %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, )
	}
	log.Printf("this is ur dog details %#v", dog)
	return context.String(http.StatusOK, "Request successfully sent")
}

func addCat(context echo.Context) error {
	cat := Cat{}
	defer context.Request().Body.Close()
	body, err := ioutil.ReadAll(context.Request().Body)
	if err != nil {
		fmt.Sprintf("failed to read the body of the request:: %s", err)
		return context.String(http.StatusInternalServerError, err.Error())
	}
	err = json.Unmarshal(body, &cat)
	if err != nil {
		fmt.Sprintf("failed to unmarshall in add cat:: %s", err)
		return context.String(http.StatusInternalServerError, err.Error())
	}
	log.Printf("this is ur cat details %#v", cat)
	return context.String(http.StatusOK, "Request successfully sent")
}

/**
start the server and test
*/
func serverStart(context echo.Context) error {
	return context.String(http.StatusOK, "Server good and running.")
}

/**
get params from url
*/
func getCats(context echo.Context) error {
	catName := context.QueryParam("name")
	catType := context.QueryParam("type")

	return context.String(http.StatusOK,
		fmt.Sprintf("Your cat name is %s \n and type is %s\n", catName, catType))
}
func getDataType(context echo.Context) error {
	catName := context.QueryParam("name")
	catType := context.QueryParam("type")
	dataType := context.Param("data")
	if dataType == "string" {
		return context.String(http.StatusOK,
			fmt.Sprintf("Your cat name is %s \n and type is %s\n", catName, catType))
	}
	if dataType == "json" {
		return context.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return context.JSON(http.StatusOK, map[string]string{
		"error": "Please specify either json or string data",
	})
}
