package main

import (
	"os"

	"github.com/m-mdy-m/psx/internal/cli"
	"github.com/m-mdy-m/psx/internal/shared"
)

var Version = "development"

func main(){
	Args:=os.Args
	if(len(Args)<2){
		shared.Help()
		os.Exit(0)
	}
//	cmd:=Args[1]
//	switch cmd {
//	case "version","-v","--version":
//		fmt.Printf("PSX version %s\n",Version)
//	case "help","-h","--help":
//		cli.Exce(Version)
//		shared.Help()
	//case "init":
	//case "fix":
	//cass "check":
	//case "rules":
	//case "config":

//	default:
//		fmt.Printf("unknown commnd: %s \n\n",cmd)
//		shared.Help()
//		os.Exit(1)
//	}
	if err:=cli.Exce(Version); err!=nil{
		os.Exit(1)
	}
}
