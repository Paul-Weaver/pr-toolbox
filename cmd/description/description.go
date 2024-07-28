package description

import (
	"context"
	"fmt"
	"os"

	"github.com/Paul-Weaver/pr-toolbox/cmd/internal/gitutils"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// descriptionCmd represents the description command
var DescriptionCmd = &cobra.Command{
	Use:   "description",
	Short: "Create a PR description from a git diff",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		baseBranch, err := gitutils.GetBaseBranch()
		if err != nil {
			fmt.Print("Error getting base branch:", err)
			return
		}

		diff, err := gitutils.GetGitDiff(baseBranch)
		if err != nil {
			fmt.Print("Error getting git diff:", err)
			return
		}

		prDescription, err := generatePRDescription(diff)
		if err != nil {
			fmt.Print("Error generating PR description:", err)
			return
		}

		output := formatOutput(prDescription)
		fmt.Println(output)
	},
}

func init() {
}

func generatePRDescription(diff string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf(`Generate a concise pull request description in Markdown format for the following git diff:
%s

Please include only the Summary and Changes sections in your response.
Important: Do not include Markdown fencing in your response.`, diff)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 1024,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func formatOutput(prDescription string) string {
	return fmt.Sprintf(`
Generated PR Description:
------------------------------------------------------
%s
------------------------------------------------------
`, prDescription)
}
