package router

//Love this example. Go to BOTTOM to see how to dummy a client
import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RedchilliSauce/sandbox/sandbox/golang/echo/4skelesite/data"
	"github.com/labstack/echo"
)

const avatarLoc string = "/Users/hesh/data/2postavatars/"

func RegisterUser(c echo.Context) error {
	username := c.Param("name")

	_, exists := data.Users[username]

	if exists {
		return c.HTML(http.StatusOK, "<b>Already registered, "+username+"</b>")
	}
	user := data.User{Name: username}

	data.Users[username] = user
	return c.HTML(http.StatusOK, "<b>You have been registered, "+username+"</b>")
	//Below code is to save avatar of user. Can be used later
	/*
		avatar, err := c.FormFile("avatar")
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		src, err := avatar.Open()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		defer src.Close()

		dst, err := os.Create(avatarLoc + avatar.Filename)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			fmt.Println(err.Error())
			return err
		}
	*/
}

//To make a call to this server, use the following example as a hint
//Refer to curl -F -> It's extraordinarily useful for this
//$ curl -F "name=Joe Smith" -F "avatar=@/path/to/your/avatar.png" http://localhost:10001/save

func SaveFlick(c echo.Context) error {
	username := c.Param("name")
	_, exists := data.Users[username]

	if !exists {
		return c.HTML(http.StatusBadRequest, "<b>User does not exist, "+username+"</b>")
	}
	flickname := c.FormValue("flickname")

	//TODO: Some problem with the conversion
	rat := c.FormValue("rating")
	rating, _ := strconv.ParseFloat(rat, 64)

	flick := data.Flick{Name: flickname, Rating: rating}
	fmt.Println(flick)
	flickList, exists := data.UserFlicks[username]

	if exists {
		data.UserFlicks[username] = append(flickList, flick)
	} else {
		flicks := make([]data.Flick, 0, 4) //if you initialize len param, it will create zero value elements
		flicks = append(flicks, flick)
		data.UserFlicks[username] = flicks
	}
	fmt.Println(data.UserFlicks)
	return c.HTML(http.StatusOK, "<b>Your flick has been added "+username+"</b>")
}

func GetUserFlicks(c echo.Context) error {
	username := c.Param("name")
	_, exists := data.Users[username]

	if !exists {
		return c.HTML(http.StatusBadRequest, "<b>User does not exist, "+username+"</b>")
	}
	flicks, exists := data.UserFlicks[username]

	if exists {
		len := len(flicks) - 1
		rating := fmt.Sprintf("%f", flicks[len].Rating)
		//		rating := strconv.FormatFloat(flicks[0].Rating, 'E', 1, 64)

		return c.HTML(http.StatusOK, "<b>Latest flick rated was "+flicks[len].Name+" "+rating+"</b>")
	} else {
		return c.HTML(http.StatusOK, "<b>No flicks registered for this user, "+username+"</b>")
	}
}
