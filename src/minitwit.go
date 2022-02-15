package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"minitwit/src/controllers"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pongA",
		})
	})

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		controllers.HandleTimeline(c.Writer, c.Request, c)
	})
	router.GET("/public", func(c *gin.Context) {
		controllers.HandlePublicTimeline(c.Writer, c.Request, c)
	})

	// Logout
	router.GET("/logout", func(c *gin.Context) {
		controllers.HandleLogout(c.Writer, c.Request, c)
	})

	// Register
	router.GET("/register", func(c *gin.Context) {
		controllers.HandleRegister(c.Writer, c.Request, c)
	})
	router.POST("/register", func(c *gin.Context) {
		controllers.HandleRegister(c.Writer, c.Request, c)
	})

	// User timeline
	router.GET("/:user", func(c *gin.Context) {
		username := c.Param("user")
		controllers.HandleUserTimeline(c.Writer, c.Request, c, username)
	})

	// Follow
	router.GET("/:user/follow", func(c *gin.Context) {
		username := c.Param("user")
		controllers.HandleFollowUser(c.Writer, c.Request, c, username)
	})

	// Unfollow
	router.GET("/:user/unfollow", func(c *gin.Context) {
		username := c.Param("user")
		controllers.HandleUnFollowUser(c.Writer, c.Request, c, username)
	})

	// Add message
	router.GET("/add_message", func(c *gin.Context) {
		controllers.HandleAddMessage(c.Writer, c.Request, c)
	})

	router.LoadHTMLFiles("./src/test.html")

	/*
	 FOR TESTING GO TOOL 'FRESH': 'go install github.com/pilu/fresh'
	 TRY TO RUN COMMAND: 'fresh -c my_fresh_runner.conf' AND
	 THEN MAKE CHANGES TO THE 'test.html' OR 'minitwit.go' FILES.
	 IF NO ERROR, THEN FRESH SHOULD BUILD AND RUN THE 'minitwit.go' CODE.
	 THE CHANGES SHOULD BE SEEN REFLECTED ON 'http://localhost:8080/test/test.html'.

	 OBS: MAYBE TURN OFF AUTO-SAVING, SO STUFF IS ONLY BUILD AND RAN, WHEN YOU WANT IT TO.
	*/
	router.Static("/test", "./src")

	router.Run(":8080")
	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
