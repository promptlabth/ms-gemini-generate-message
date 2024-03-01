package promptMessage

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
)

type PromptGeminiService interface {
	GenerateMessage(ctx context.Context, prompt string) (string, error)
	GenerateStreamMessage(ctx context.Context, prompt string, MessageStream chan string) error
}

type PromptMessageStorage interface {
	Save(ctx context.Context, promptMessage PromptMessage) (*PromptMessage, error)
}

type promptMessageHandler struct {
	promptMessageStorage PromptMessageStorage
	geminiService        PromptGeminiService
}

func NewPromptMessageHandler(
	promptMessageStorage PromptMessageStorage,
	geminiService PromptGeminiService,
) *promptMessageHandler {
	return &promptMessageHandler{
		promptMessageStorage: promptMessageStorage,
		geminiService:        geminiService,
	}
}

func (p promptMessageHandler) GenerateMessage(c *gin.Context) {
	ctx := c.Request.Context()
	promptMessage := new(PromptMessageRequest)
	// binding a request to json
	if err := c.ShouldBindJSON(promptMessage); err != nil {
		c.JSON(400, map[string]string{
			"error": err.Error(),
		})
		return
	}
	// use content to generate message with gemini
	content, err := p.geminiService.GenerateMessage(ctx, promptMessage.InputMessage)
	if err != nil {
		c.JSON(503, map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]string{
		"data": content,
	})
}

func (p promptMessageHandler) GenerateStream(c *gin.Context) {
	ginCtx := c.Request.Context()

	ctx, cancle := context.WithCancel(ginCtx)
	defer cancle()

	// create channel to get data
	rawContent := make(chan string)
	// content := make(chan rune, 50)
	// create mutex to control routine flow

	go p.geminiService.GenerateStreamMessage(
		ctx,
		"detail of level of tester ",
		rawContent, // to recive data
	)

	// go func() {
	// 	for {
	// 		d, ok := <-rawContent
	// 		if !ok {
	// 			close(content)
	// 			return
	// 		}
	// 		for _, v := range d {
	// 			content <- v
	// 			time.Sleep(100 * time.Millisecond)
	// 		}
	// 	}
	// }()

	go c.Stream(func(w io.Writer) bool {
		message, ok := <-rawContent
		if !ok {
			// case end of content
			c.SSEvent("status", "done")
			cancle()
			return false
		}
		// resp struct

		c.SSEvent("message2", map[string]string{
			"data": string(message),
		})
		return true
	})
	<-ctx.Done()

}
