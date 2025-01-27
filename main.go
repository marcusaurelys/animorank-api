package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	var chanSize int
	derr := godotenv.Load()
	num, err := strconv.ParseInt(os.Getenv("CHAN_SIZE"), 10, 0)
	if derr != nil || err != nil {
		fmt.Println("Can not detect environment variable for channel size, setting default to 100.")
		chanSize = 100
	} else {
		chanSize = int(num)
	}

	r := gin.Default()
	r.Use(cors.Default())

	exec_buffer := make(chan int, chanSize) //buffer to only execute one 100 jobs at a time
	fmt.Printf("Can handle %v requests concurrently\n", chanSize)
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
			exec_buffer <- 1
			go func() {
				defer wg.Done()
				results[index] = compile(data.Code, input)
				<-exec_buffer
			}()
		}
		wg.Wait()
		c.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	})

	r.Run()
}

func compile(code string, stdin string) Result {
	codeVar := fmt.Sprintf("CODE=%v", code)
	stdinVar := fmt.Sprintf("STDIN=%v", stdin)
	cmd := exec.Command("docker", "run", "--rm", "--env", codeVar, "--env", stdinVar, "--network", "none", "--memory=50m", "--memory-swap=50m", "--stop-timeout", "10", "--security-opt=no-new-privileges:true", "compile-job")
	runOutput, err := cmd.CombinedOutput()
	status := handleCompileError(err)
	return Result{string(runOutput), status}
}

func handleCompileError(err error) string {
	if err != nil {
		switch fmt.Sprintf("%v", err) {
		case "exit status 4":
			return "runtime error"
		case "exit status 3":
			return "compile error"
		case "exit status 2":
			return "time limit exceeded"
		case "exit status 1":
			return "code cant be written to file"
		case "exit status 137":
			return "out of memory error"
		default:
			return fmt.Sprintf("%v", err)
		}
	}
	return "successfully compiled and run"
}
