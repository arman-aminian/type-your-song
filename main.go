package main

import (
	"github.com/arman-aminian/type-your-song/db"
	"github.com/arman-aminian/type-your-song/handler"
	"github.com/arman-aminian/type-your-song/router"
	"github.com/arman-aminian/type-your-song/store"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func main() {
	//
	//r.GET("/swagger/*", echoSwagger.WrapHandler)
	//
	//

	//to := []string{
	//	"arman.aminian78@gmail.com",
	//}
	//content := "https://github.com"
	//err := email.SendEmail(to, content, "confirm your typeasong account")
	//if err != nil {
	//	panic(err)
	//}

	r := router.New()
	v1 := r.Group("/api")
	mongoClient, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	usersDb := db.SetupUsersDb(mongoClient)
	songsDb := db.SetupSongsDb(mongoClient)
	artistsDb := db.SetupArtistsDb(mongoClient)
	us := store.NewUserStore(usersDb)
	ss := store.NewSongStore(songsDb)
	as := store.NewArtistStore(artistsDb)
	h := handler.NewHandler(us, ss, as)

	r.GET("/", Temp)
	h.Register(v1)
	r.Logger.Fatal(r.Start(utils.BaseUrl))
}

func Temp(c echo.Context) error {
	var htmlIndex = `<html>
<body>
	<a href="/api/users/login/google">Google Log In</a>
</body>
</html>`
	return c.HTML(http.StatusTemporaryRedirect, htmlIndex)
}
