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

func main() {
	r := gin.Default()
	r.Use(cors.Default())

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
		stdout := make([]string, len(data.Stdin))
		for index, input := range data.Stdin {
			wg.Add(1)
			go func() {
				defer wg.Done()
				stdout[index] = compile(data.Code, input)
			}()
		}
		wg.Wait()
		c.JSON(http.StatusOK, gin.H{
			"stdout": stdout,
		})
	})

	// r.POST("/submit", func(c *gin.Context){})
	//language, code, []stdout, []test cases

	r.Run()
}

func compile(code string, stdin string) string {
	codeVar := fmt.Sprintf("CODE=%v", code)
	stdinVar := fmt.Sprintf("STDIN=%v", stdin)
	cmd := exec.Command("docker", "run", "--rm", "--env", codeVar, "--env", stdinVar, "compile-job")
	runOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%v", err)
	}
	return string(runOutput)

}
