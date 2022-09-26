package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "code_user/main.go")
	var out, stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(pipe, "23 11\n")
	// 根据测试的输入案例，进行运行，拿到输出结果和标准的输出结果进行比对
	if err := cmd.Run(); err != nil {
		log.Fatal(err, stdErr.String())
	}
	fmt.Printf("out.String(): %v\n", out.String())
}
