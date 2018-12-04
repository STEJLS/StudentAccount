package main

import (
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

	u.InitDB(config.Admin.Login, config.Admin.Password)

	http.ListenAndServe("0.0.0.0:3000", ctx.GetRoots())
}
