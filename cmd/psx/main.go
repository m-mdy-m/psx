package main

import (
	"fmt"
	"os"

	"github.com/m-mdy-m/psx/internal/shared"
)

func main(){
	Args:=os.Args
	if(len(Args)<2){
		shared.Help()
		os.Exit(0)
	}
	cmd:=Args[1]
	fmt.Printf("CMD: %T\n",cmd)
}
