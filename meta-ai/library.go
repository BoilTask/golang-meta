package metaai

import (
	"context"
	"fmt"
	"log"
	metapanic "meta/meta-panic"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

func Test(url string, apiKey string) {
	client := openai.NewClient(
		option.WithBaseURL(url),
		option.WithAPIKey(apiKey),
	)
	ctx := context.Background()

	// 创建流式 Chat Completion
	stream := client.Chat.Completions.NewStreaming(
		ctx, openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("测试流式输出"),
			},
			Model: "THUDM/GLM-4.1V-9B-Thinking",
		},
	)
	defer func(stream *ssestream.Stream[openai.ChatCompletionChunk]) {
		err := stream.Close()
		if err != nil {
			metapanic.ProcessError(err)
		}
	}(stream)

	fmt.Println("==== 开始流式输出 ====")

	// 持续读取流
	for stream.Next() {
		event := stream.Current()

		// 每个 event 里可能包含多个 choice
		for _, choice := range event.Choices {
			if choice.Delta.Content != "" {
				fmt.Print(choice.Delta.Content)
			}
		}
	}

	if err := stream.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n==== 流式输出结束 ====")
}
