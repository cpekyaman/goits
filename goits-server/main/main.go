package main

import (
	"github.com/cpekyaman/goits/config"
	"github.com/cpekyaman/goits/framework/orm/db"

	"github.com/cpekyaman/goits/cli"
)

func main() {
	config.InitConfig()
	db.NewDB()

	cli.Execute()
}
