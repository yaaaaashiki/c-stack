package cstack

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yaaaaashiki/cstack/db"
	"github.com/yaaaaashiki/cstack/domain/repository"
	"github.com/yaaaaashiki/cstack/interfaceadapter/controller"
	"github.com/yaaaaashiki/cstack/usecase"
)

// This holds database connection and router settings based on gin.
type Server struct {
	db  *gorm.DB
	gin *gin.Engine
}

// New returns server object.
func New() *Server {
	return &Server{}
}

// Close makes the database connection to close.
func (s *Server) Close() error {
	return s.db.Close()
}

// Init initialize server state. Connecting to database, compiling templates,
// and settings router.
func (s *Server) Init(dbconf, env string, debug bool) {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		log.Fatalf("cannot open database configuration. exit. %s", err)
	}
	db, err := cs.Open(env)
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}

	s.db = db
	s.gin = gin.Default()
	s.Route()
}

// Run starts running http server.
func (s *Server) Run(addr string) {
	log.Printf("start listening on %s", addr)

	s.gin.Run(addr)
}

func (s *Server) Route() {
	r := s.gin

	// cookie manager initialize
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/image", "./assets/image")
	r.Static("/css", "./assets/css")
	r.Static("/js", "./assets/js")

	r.LoadHTMLGlob("view/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	api := r.Group("/api")

	userRepository := repository.NewUserRepository(s.db)
	itemRepository := repository.NewItemRepository(s.db)

	registerUserCase := usecase.NewRegisterUseCase(userRepository)
	loginUseCase := usecase.NewLoginUseCase(userRepository)
	findAllItemsUseCase := usecase.NewFindAllItemsUseCase(itemRepository)
	registerItemUserCase := usecase.NewRegisterItemUseCase(itemRepository)

	registerUserController := controller.NewRegisterController(registerUserCase)
	loginController := controller.NewLoginController(loginUseCase)
	logoutController := controller.NewLogoutController()
	findAllItemsController := controller.NewFindAllItemsController(findAllItemsUseCase)
	registerItemController := controller.NewRegisterItemController(registerItemUserCase)

	//auth
	api.POST("/users", registerUserController.Execute)
	api.POST("/auth", loginController.Execute)
	api.DELETE("/auth", logoutController.Execute)

	//find all items by user id, register item
	api.GET("/users/:userID/items", findAllItemsController.Execute)
	api.POST("/users/:userID/items", registerItemController.Execute)
}
