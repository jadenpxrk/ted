package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	agentSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"command": {
				Type:        genai.TypeString,
				Description: "The executable command that accomplishes the task",
			},
			"explanation": {
				Type:        genai.TypeString,
				Description: "Brief explanation of what the command does",
			},
		},
		Required: []string{"command", "explanation"},
	}

	askSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"commands": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"command": {
							Type:        genai.TypeString,
							Description: "The command to execute",
						},
						"description": {
							Type:        genai.TypeString,
							Description: "Description of what the command does",
						},
					},
					Required: []string{"command", "description"},
				},
			},
		},
		Required: []string{"commands"},
	}
)

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

type AgentResponse struct {
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
}

type AskResponse struct {
	Commands []CommandOption `json:"commands"`
}

type CommandOption struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

func NewClient(apiKey, modelName string, temperature float32) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	model.SetTemperature(temperature)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) GenerateAgentCommand(ctx context.Context, query string) (*AgentResponse, error) {
	c.model.ResponseMIMEType = "application/json"
	c.model.ResponseSchema = agentSchema

	prompt := fmt.Sprintf(`You are a helpful command-line assistant. The user wants to accomplish the following task: "%s"

Please respond with a JSON object containing the command and explanation.`, query)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	var response AgentResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}

func (c *Client) GenerateAskCommands(ctx context.Context, question string) (*AskResponse, error) {
	c.model.ResponseMIMEType = "application/json"
	c.model.ResponseSchema = askSchema

	prompt := fmt.Sprintf(`The user is asking: "%s"

Please provide exactly 3 different command-line commands that help answer this question. Return a JSON object with a "commands" array.`, question)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	var response AskResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}
