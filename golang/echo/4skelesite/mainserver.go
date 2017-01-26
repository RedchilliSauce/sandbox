package main

import (
	"github.com/RedchilliSauce/sandbox/sandbox/golang/echo/4skelesite/router"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/register/:name", router.RegisterUser)
	e.GET("/getflicks/:name", router.GetUserFlicks)
	e.POST("/addflick/:name", router.SaveFlick)
	e.Start(":10005")
}

/*
curl http://localhost:10005/register/anish
curl http://localhost:10005/getflicks/anish
curl -F "flickname=Hannibal" -F "rating=7.5" http://localhost:10005/addflick/anish
*/
