package gemini

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type geminiMessage struct {
	genService *genai.Client
}

func NewGeminiService() *geminiMessage {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GENAI_API_KEY")))
	if err != nil {
		panic(err)
	}
	return &geminiMessage{
		genService: client,
	}
}

func (g *geminiMessage) GenerateMessage(ctx context.Context, prompt string) (string, error) {
	model := g.genService.GenerativeModel("gemini-pro")
	content, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", content.Candidates[0].Content.Parts[0]), nil
}

func (g *geminiMessage) GenerateStreamMessage(ctx context.Context, prompt string, MessageStream chan string) error {
	model := g.genService.GenerativeModel("gemini-pro")
	iter := model.GenerateContentStream(ctx, genai.Text(prompt))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			close(MessageStream)
			break
		}
		if err != nil {
			return err
		}
		// Manage Will Assign Data when signal sender is successful
		MessageStream <- fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
	}
	return nil
}
