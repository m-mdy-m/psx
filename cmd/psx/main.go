package main

import (
	"fmt"
	"os"

	"github.com/m-mdy-m/psx/internal/shared"
)

var Version = "development"

func main(){
	Args:=os.Args
	if(len(Args)<2){
		shared.Help()
		os.Exit(0)
	}
	cmd:=Args[1]
	switch cmd {
	case "version","-v","--version":
		fmt.Printf("PSX version %s\n",Version)

	case "help","-h","--help":
		shared.Help()
	default:
		fmt.Printf("unknown commnd: %s \n\n",cmd)
		shared.Help()
		os.Exit(1)
	}
}
