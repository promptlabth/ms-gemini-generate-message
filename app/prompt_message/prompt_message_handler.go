package promptMessage

import (
	"context"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

type PromptGeminiService interface {
	GenerateMessage(ctx context.Context, prompt string) (string, error)
	GenerateStreamMessage(ctx context.Context, prompt string, MessageStream chan string, mu *sync.Mutex) error
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
	content, err := p.geminiService.GenerateMessage(ctx, "hellp")
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

	ctx, cancel := context.WithCancel(ginCtx)
	defer cancel()

	// create channel to get data
	content := make(chan string)
	// create mutex to control routine flow
	var mu = new(sync.Mutex)
	go p.geminiService.GenerateStreamMessage(
		ctx,
		"level of mindset of tester",
		content, // to recive data
		mu,      // to routine control
	)

	go c.Stream(func(w io.Writer) bool {
		message := <-content
		if message == "" {
			// case end of content
			cancel()
		}
		mu.Lock()
		for _, v := range message {
			c.SSEvent("message", string(v))
		}
		mu.Unlock()
		return true

	})
	<-ctx.Done()
}
