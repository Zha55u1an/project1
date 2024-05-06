package router

import "github.com/gin-gonic/gin"

func InitRouter() (*gin.Engine, error) {
	router := gin.Default()

	//err := router.Run()
	//if err != nil {
	//	return nil, err
	//}
	return router, nil
}
