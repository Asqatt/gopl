package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main()  {
	cmd := &exec.Cmd{
		Path:"/home/asqad/Documents/proxy/Qv2ray-2.6.3",
		Args:[]string{},
		Dir:"/home/asqad/Documents/proxy/",
		Stdout:os.Stdout,
		Stdin:os.Stdin,
	}

	err :=cmd.Start()

	fmt.Println(err)
	

	cmd.Wait()
}