package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Runnable struct {
	Language string   `json:"language" binding:"required"`
	Code     string   `json:"code" binding:"required"`
	Stdin    []string `json:"stdin" binding:"required"`
}

type Result struct {
	Stdout string `json:"stdout" binding:"required"`
	Status string `json:"status" binding:"required"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	exec_queue := make(chan int, 100) //buffer to only execute one 100 jobs at a time

	r.GET("/helloworld", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	//language, code, []stdin
	r.POST("/run", func(c *gin.Context) {
		var data Runnable
		err := c.BindJSON(&data)
		if err != nil {
			c.AbortWithError(400, err)
		}

		var wg sync.WaitGroup
		results := make([]Result, len(data.Stdin))
		for index, input := range data.Stdin {
			wg.Add(1)
			exec_queue <- 1
			go func() {
				defer wg.Done()
				results[index].Stdout, results[index].Status = compile(data.Code, input)
				<-exec_queue
			}()
		}
		wg.Wait()
		c.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	})

	r.Run()
}

func compile(code string, stdin string) (string, string) {
	codeVar := fmt.Sprintf("CODE=%v", code)
	stdinVar := fmt.Sprintf("STDIN=%v", stdin)
	status := "successfully compiled and run"
	cmd := exec.Command("docker", "run", "--rm", "--env", codeVar, "--env", stdinVar, "--network", "none", "--memory=50m", "--memory-swap=50m", "compile-job")
	runOutput, err := cmd.CombinedOutput()
	if err != nil {
		switch fmt.Sprintf("%v", err) {
		case "exit status 4":
			status = "runtime error"
		case "exit status 3":
			status = "compile error"
		case "exit status 2":
			status = "time limit exceeded"
		case "exit status 1":
			status = "code cant be written to file"
		case "exit status 137":
			status = "out of memory error"
		}
	}
	return string(runOutput), status
}
