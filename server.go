package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/otoolep/go-httpd/store"
	"log"
	"net/http"
	"strings"
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

	r.GET("/ws/:name", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		msgfragments := strings.Split(string(msg[:]), "|")
		from, to := msgfragments[0], msgfragments[1]
		if to == "broadcast" {
			m.Broadcast(msg)
			log.Println("Broadcast msg - ", msg)
		} else {
			log.Println("Private msg - ", msg, " from - ", from, " to - ", to)
			m.BroadcastFilter(msg, func(q *melody.Session) bool {
				path := q.Request.URL.Path
				fragments := strings.Split(path, "/")
				log.Println("Compare fragment - ", fragments[2], " with from - ", from, "and to - ", to)
				return strings.Compare(fragments[2], from) == 0 || strings.Compare(fragments[2], to) == 0
			})
		}
	})

	r.Run(":5000")
}
