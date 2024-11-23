package main

import (
	_ "github.com/nhutHao02/social-network-tweet-service/pkg/common"
	"github.com/nhutHao02/social-network-tweet-service/startup"
)

//	@title			Social Network Service
//	@description	This is tweet service of the social network implament using Go
//	@version		1.0
//	@BasePath		/api/v1
func main() {
	startup.Start()
}
