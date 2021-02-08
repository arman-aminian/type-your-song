package handler

type Handler struct {
	//userStore    user.Store
	//articleStore article.Store
}

//func NewHandler(us user.Store, as article.Store) *Handler {
func NewHandler() *Handler {
	return &Handler{
		//userStore:    us,
		//articleStore: as,
	}
}
