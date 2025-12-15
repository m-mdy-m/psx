package main

import (
	"os"

	"github.com/m-mdy-m/psx/internal/commond"
)

var Version = "development"

func main(){
	if err:=commond.Exce(Version); err!=nil{
		os.Exit(1)
	}
}
