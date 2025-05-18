package helper

import "github.com/gin-gonic/gin"

func RespondWithError(g *gin.Context, code int, message string) {
	g.JSON(500, gin.H{"msg": message,"code": code,"data": nil})
}

func RespondWithSuccess(g *gin.Context, code int, message string, data interface{}) {
	g.JSON(200, gin.H{"msg": message, "data": data,"code": code})
}