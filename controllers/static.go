package controllers

import "github.com/gin-gonic/gin"

func staticHandlers(router *gin.Engine) {

	router.Static("/static", "./static")

}
