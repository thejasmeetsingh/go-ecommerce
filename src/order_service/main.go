package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	println()
	r := gin.Default()

	log.Infoln("Order services started!")
	r.Run(":8000")
}
