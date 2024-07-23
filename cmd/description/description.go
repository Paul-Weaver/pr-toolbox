package description

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// descriptionCmd represents the description command
var DescriptionCmd = &cobra.Command{
	Use:   "description",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		diff, err := getGitDiff()
		if err != nil {
			fmt.Print("Error getting git diff:", err)
			return
		}

		prDescription, err := generatePRDescription(diff)
		if err != nil {
			fmt.Print("Error generating PR description:", err)
			return
		}
		fmt.Println("Generated PR Description (Markdown):")
		fmt.Println(prDescription)
	},
}

func init() {
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "master")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func generatePRDescription(diff string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf("Generate a concise pull request description in Markdown format for the following git diff:\n%s", diff)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
