package service

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AzureOpenAI struct {
	Model  string
	Client *azopenai.Client
}

func NewOpenAI(ctx context.Context) *AzureOpenAI {

	model := "gpt-35-turbo"

	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	if endpoint == "" {
		log.Fatal("AZURE_OPENAI_ENDPOINT is not set")
	}

	secretAPIKey := os.Getenv("AZURE_OPENAI_API_KEY")
	if secretAPIKey == "" {
		log.Fatal("AZURE_OPENAI_API_KEY is not set")
	}

	keyCredential := azcore.NewKeyCredential(secretAPIKey)
	client, err := azopenai.NewClientWithKeyCredential(endpoint, keyCredential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// resp, err := client.GetCompletions(context.TODO(), azopenai.CompletionsOptions{
	// 	Prompt:         []string{"What is Azure OpenAI, in 20 words or less"},
	// 	MaxTokens:      to.Ptr(int32(2048)),
	// 	Temperature:    to.Ptr(float32(0.0)),
	// 	DeploymentName: &model,
	// }, nil)

	return &AzureOpenAI{
		Model:  model,
		Client: client,
	}
}
