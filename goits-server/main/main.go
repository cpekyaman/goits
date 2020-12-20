package main

import (
	"github.com/cpekyaman/goits/config"
	"github.com/cpekyaman/goits/framework/orm"

	"github.com/cpekyaman/goits/cli"
)

func main() {
	config.InitConfig()
	orm.NewDB()

	cli.Execute()
}
