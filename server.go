package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/otoolep/go-httpd/store"
	"net/http"
)

func main() {
	r := gin.Default()
	m := melody.New()
	users := store.New()

	if err := users.Open(); err != nil {
		/*
			Handle Error with error page
		*/
	}

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "login.html")
	})

	r.GET("/room/:name", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/login/:name", func(c *gin.Context) {
		name := c.Param("name")
		v, err := users.Get(name)
		if v == "" && err == nil {
			users.Set(name, "1")
			c.Redirect(http.StatusMovedPermanently, "/room/"+name)
		}
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	r.Run(":5000")
}
