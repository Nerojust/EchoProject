package main

import (
	"EchoProject/models"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("Starting")
	//start echo
	e := echo.New()
	e.Use(setUpSErverHeader)
	//for grouping routes to one category.
	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")

	//using middleware
	//to format the logs to view on the console
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          nil,
		Format:           `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
		CustomTimeFormat: "",
		Output:           nil,
	}))
	adminGroup.Use(middleware.BasicAuth(func(username string, password string, context echo.Context) (b bool, e error) {
		//check in the database if password is valid
		if username == "jack" && password == "1234" {
			return true, nil
		}
		return false, nil
	}))
	adminGroup.Use(middleware.CORS())
	adminGroup.GET("/main", mainAdmin)
	cookieGroup.GET("/main", mainCookie)
	cookieGroup.Use(checkCookie)
	e.GET("/", serverStart)
	e.GET("/login", login)
	e.GET("/cats", getCats)
	e.GET("/cats/:data", getDataType)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)

	_ = e.Start(":8000")
}

func mainCookie(context echo.Context) error {
	return context.String(http.StatusOK, "you are on the secret page")
}
func setUpSErverHeader(e echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		context.Response().Header().Set(echo.HeaderServer, "Server 1.0")
		return e(context)
	}
}

func login(context echo.Context) error {
	username := context.QueryParam("username")
	password := context.QueryParam("password")
	//perform check in db if valid after hashing it
	if username == "jack" && password == "1234" {
		cookie := &http.Cookie{}
		//this is the same
		//cookie:= new(http.Cookie)
		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)

		context.SetCookie(cookie)
		return context.String(http.StatusOK, "You are logged in")
	}
	return context.String(http.StatusInternalServerError, "You username or password is invalid")

}
func checkCookie(e echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		cookie, err := context.Cookie("sessionID")
		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return context.String(http.StatusUnauthorized, "You dont have any cookie")
			}
			log.Println(err)
			return err
		}
		if cookie.Value == "some_string" {
			return e(context)
		}
		return context.String(http.StatusUnauthorized, "Wrong credentials")
	}
}

func mainAdmin(context echo.Context) error {
	return context.String(http.StatusOK, "Secret admin")
}

/**
we create 3 different methods of processing requests.
*/
func addHamster(context echo.Context) error {
	hamster := models.Hamster{}
	err := context.Bind(&hamster)
	if err != nil {
		log.Printf("failed processing add hamster request %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("this is ur hamster details %#v", hamster)
	return context.String(http.StatusOK, "Hamster request successfully sent")
}

func addDog(context echo.Context) error {
	dog := models.Dog{}
	defer context.Request().Body.Close()
	err := json.NewDecoder(context.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("failed processing add dog request %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("this is ur dog details %#v", dog)
	return context.String(http.StatusOK, "Request successfully sent")
}

func addCat(context echo.Context) error {
	cat := models.Cat{}
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
