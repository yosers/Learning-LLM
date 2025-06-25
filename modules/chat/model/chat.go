package model

type ChatPayload struct {
	Message   string `json:"message" binding:"required"`
	ChannelID int    `json:"channel_id" binding:"required,min=1"`
}

type ChatMessage struct {
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	Name      *string        `json:"name"`
	ToolCalls *[]interface{} `json:"tool_calls"`
}

type ChatChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason string       `json:"finish_reason"`
	Logprobs     *interface{} `json:"logprobs"`
}

type ChatUsage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	EstimatedCost    float64 `json:"estimated_cost"`
}

type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   ChatUsage    `json:"usage"`
}

type ChatResponse struct {
	Message      string                 `json:"message"`
	FullResponse ChatCompletionResponse `json:"full_response"`
}

type ChatSession struct {
	UserID    int `json:"user_id" binding:"required"`
	ChannelID int `json:"channel_id" binding:"required"`
}

type ChatMessagePayload struct {
	Message   string `json:"message"`
	SessionID int32  `json:"session_id" binding:"required"`
	UserID    int    `json:"user_id" binding:"required"`
	ChannelID int    `json:"channel_id" binding:"required"`
}
