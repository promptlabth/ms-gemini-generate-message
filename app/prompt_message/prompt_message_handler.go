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

func (p promptMessageHandler) Generate(c *gin.Context) {
	ctx := c.Request.Context()
	content, err := p.geminiService.GenerateMessage(ctx, "tell me a long story in thai lang")
	if err != nil {
		c.JSON(404, map[string]string{
			"err": err.Error(),
		})
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
