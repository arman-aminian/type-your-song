package type_your_song

import "github.com/arman-aminian/type-your-song/router"

func main() {
	r := router.New()
	//
	//r.GET("/swagger/*", echoSwagger.WrapHandler)
	//
	v1 := r.Group("/api")
	//
	//d := db.New()
	//db.AutoMigrate(d)
	//
	//us := store.NewUserStore(d)
	//as := store.NewArticleStore(d)
	//h := handler.NewHandler(us, as)
	//h.Register(v1)
	//r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}
