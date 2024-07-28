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
		// Debug statement: Start of command
		fmt.Println("Starting description command")

		baseBranch, err := gitutils.GetBaseBranch()
		if err != nil {
			fmt.Println("Error determining base branch:", err)
			return
		}
		// Debug statement: Base branch
		fmt.Printf("Determined base branch: %s\n", baseBranch)

		// Generate the diff
		diff, err := gitutils.GetGitDiff(baseBranch)
		if err != nil {
			fmt.Println("Error generating git diff:", err)
			return
		}

		// Check if the diff is empty
		if diff == "" {
			fmt.Println("No changes detected. Please commit changes before generating a PR description.")
			return
		}

		// Debug statement: Git diff
		fmt.Printf("Generated git diff:\n%s\n", diff)

		prDescription, err := generatePRDescription(diff)
		if err != nil {
			fmt.Println("Error generating PR description:", err)
			return
		}
		// Debug statement: PR description
		fmt.Printf("Generated PR description:\n%s\n", prDescription)

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
