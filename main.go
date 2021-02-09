package main

import (
	"github.com/arman-aminian/type-your-song/db"
	"github.com/arman-aminian/type-your-song/handler"
	"github.com/arman-aminian/type-your-song/router"
	"github.com/arman-aminian/type-your-song/store"
	"log"
)

func main() {
	//
	//r.GET("/swagger/*", echoSwagger.WrapHandler)
	//
	//

	r := router.New()
	v1 := r.Group("/api")
	mongoClient, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	usersDb := db.SetupUsersDb(mongoClient)
	us := store.NewUserStore(usersDb)
	h := handler.NewHandler(us)

	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))

}