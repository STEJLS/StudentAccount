package main

import (
	"fmt"
	"net/http"

	xmlc "github.com/STEJLS/StudentAccount/XMLconfig"
	ctx "github.com/STEJLS/StudentAccount/contexts"
	g "github.com/STEJLS/StudentAccount/globals"
	u "github.com/STEJLS/StudentAccount/utils"
)

func main() {
	u.InitFlags()
	u.InitFiles()
	logFile := u.InitLogger()
	defer logFile.Close()

	config := xmlc.Get(g.ConfigSource)

	u.ConnectToDB(config.DB)
	defer g.DB.Close()

	u.InitDB(config)

	http.ListenAndServe(fmt.Sprintf("%v:%v", config.HTTP.Host, config.HTTP.Port), ctx.GetRoots())
}
