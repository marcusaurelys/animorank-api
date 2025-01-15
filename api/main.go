package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Runnable struct {
	Language string   `json:"language" binding:"required"`
	Code     string   `json:"code" binding:"required"`
	Stdin    []string `json:"stdin" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/helloworld", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	//language, code, []stdin, []test cases
	r.POST("/run", func(c *gin.Context) {
		var data Runnable
		err := c.BindJSON(&data)
		if err != nil {
			c.AbortWithError(400, err)
		}
		stdout := make([]string, len(data.Stdin))
		for index, input := range data.Stdin {
			stdout[index] = compile(data.Code, input)
		}
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
