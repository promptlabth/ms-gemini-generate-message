package promptMessage

import (
	"context"
	"io"
	"sync"
	"time"

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
	ctx := c.Request.Context()
	// create channel to get data
	content := make(chan string)
	// create mutex to control routine flow
	var mu = new(sync.Mutex)
	go p.geminiService.GenerateStreamMessage(
		ctx,
		"level of mindset of tester",
		content,
		mu,
	)

	c.Stream(func(w io.Writer) bool {
		select {
		case message := <-content:
			if message == "" {
				// case end of content
				return false
			}
			mu.Lock()
			for _, v := range message {
				select {
				case <-ctx.Done():
					mu.Unlock()
					return false
				default:
					c.SSEvent("message", string(v))
					time.Sleep(2 * time.Millisecond)
				}
			}
			mu.Unlock()
			return true
		case <-ctx.Done():
			return false
		}
	})
}
