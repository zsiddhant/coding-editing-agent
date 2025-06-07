package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ToolInputSchemaParam struct {
	Properties  any            `json:"properties,omitzero"`
	Required    []string       `json:"required,omitzero"`
	ExtraFields map[string]any `json:"-"`
}

type ToolDefinition struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	InputSchema ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

func NewAgent(modelName string, getUserMessage func() (string, bool), tools []ToolDefinition) *Agent {
	return &Agent{
		model:          modelName,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

type Agent struct {
	model          string
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
}

type MessageParam struct {
	Content []string
	Role    string `json:"role,omitempty"`
}

func (a *Agent) Run() error {
	conversation := []MessageParam{}
	fmt.Println("Chat with " + a.model + " (use 'ctrl-c' to quit)")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		conversation = append(conversation, MessageParam{
			Content: []string{userInput},
			Role:    "user",
		})

		message, err := a.runInference(conversation, a.tools)
		if err != nil {
			return err
		}
		conversation = append(conversation, message)

		for _, content := range message.Content {
			fmt.Printf("\u001b[93m"+a.model+"\u001b[0m: %s\n", content)
		}
	}

	return nil
}

func (a *Agent) runInference(conversation []MessageParam, _ []ToolDefinition) (MessageParam, error) {
	// Call local Ollama model via HTTP API
	// Ollama expects messages as [{"role": "...", "content": "..."}]
	var ollamaMessages []map[string]string
	for _, msg := range conversation {
		content := ""
		if len(msg.Content) > 0 {
			content = msg.Content[0]
		}
		ollamaMessages = append(ollamaMessages, map[string]string{
			"role":    msg.Role,
			"content": content,
		})
	}
	// Build tools array as per Ollama API sample
	ollamaTools := []map[string]any{}
	for _, tool := range a.tools {
		toolMap := map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        tool.Name,
				"description": tool.Description,
				"parameters":  tool.InputSchema.Properties,
			},
		}
		ollamaTools = append(ollamaTools, toolMap)
	}
	reqBody := map[string]any{
		"model":    a.model, // or your local model name
		"messages": ollamaMessages,
		"tools":    ollamaTools,
		"stream":   false,
	}
	req, err := json.Marshal(reqBody)
	if err != nil {
		return MessageParam{}, err
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(req))
	if err != nil {
		return MessageParam{}, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	fmt.Printf("Response from Ollama: %s\n", respBody)
	defer resp.Body.Close()

	// Sample response from Ollama:
	// {
	//   "message": {
	//     "role": "assistant",
	//     "content": "Hello! How can I help you today?"
	//   }
	// }
	var ollamaResp struct {
		Message struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []struct {
				Function struct {
					Name      string          `json:"name"`
					Arguments json.RawMessage `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return MessageParam{}, err
	}

	if ollamaResp.Message.ToolCalls != nil {
		// Handle tool calls if any
		for _, toolCall := range ollamaResp.Message.ToolCalls {
			for _, tool := range a.tools {
				if tool.Name == toolCall.Function.Name {
					result, err := tool.Function(toolCall.Function.Arguments)
					if err != nil {
						return MessageParam{}, fmt.Errorf("error calling tool %s: %w", tool.Name, err)
					}
					//fmt.Printf("Tool %s returned: %s\n", tool.Name, result)

					// Append tool result to conversation
					conversation = append(conversation, MessageParam{
						Content: []string{result},
						Role:    "tool",
					})

					// Call the model again with the updated conversation
					return a.runInference(conversation, a.tools)
				}
			}
		}
	}

	return MessageParam{
		Content: []string{ollamaResp.Message.Content},
		Role:    "assistant",
	}, nil
}

func main() {
	tools := []ToolDefinition{ReadFileDefinition}
	model := "llama3.2" // or your local model name

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	agent := NewAgent(model, getUserMessage, tools)
	err := agent.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
