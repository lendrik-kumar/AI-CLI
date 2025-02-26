package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// Root command for the CLI
var rootCmd = &cobra.Command{
	Use:   "skm",
	Short: "An AI-powered CLI to execute natural language commands",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a command.")
			return
		}

		// Combine all arguments into a single string
		naturalCommand := args[0]
		fmt.Println("Interpreting command:", naturalCommand)

		// Get AI-generated terminal command
		terminalCommand := interpretCommand(naturalCommand)
		fmt.Println("AI-generated command:", terminalCommand)

		// Execute the generated command
		executeCommand(terminalCommand)
	},
}

// Function to call OpenAI GPT API
func getGPTResponse(prompt string) (string, error) {
	ctx := context.Background()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing OPENAI_API_KEY")
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Generate response
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are an AI that converts natural language into terminal commands."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}

	// Return the response text
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "No response from OpenAI API", nil
}

// Function to interpret natural language into a command
func interpretCommand(command string) string {
	response, err := getGPTResponse("Convert this natural language to a terminal command: " + command)
	if err != nil {
		log.Println("Error calling OpenAI API:", err)
		return "Error: Unable to process command"
	}
	return response
}

// Function to execute a system command
func executeCommand(command string) {
	fmt.Println("Executing:", command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}

// Initialize the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}