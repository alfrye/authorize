package main

import (
	"fmt"

	"github.com/alfrye/authorize/internal/server"
)

func main() {

	fmt.Println("starting point for Authorize")
	s := server.New("9010")

	s.PopulateRoutes(s.AuthorizeServiceRoutes())
	s.Listen()

}
