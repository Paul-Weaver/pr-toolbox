package description

import (
	"context"
	"fmt"
	"os"

	"github.com/Paul-Weaver/pr-toolbox/cmd/internal/gitutils"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var detailLevel int

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

		prDescription, err := generatePRDescription(diff, detailLevel)
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
	// Define the detail level flag with a default value of 1 (medium detail)
	DescriptionCmd.Flags().IntVarP(&detailLevel, "detail-level", "d", 1, "detail level for the PR description (0, 1, 2)")
}

// GeneratePRDescription generates a pull request description from a git diff
func generatePRDescription(diff string, detailLevel int) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient(apiKey)

	// Adjust max tokens and prompt based on detail level
	var maxTokens int
	var prompt string
	switch detailLevel {
	case 0:
		maxTokens = 256 // low detail
		prompt = fmt.Sprintf(`Generate a brief pull request description in Markdown format for the following git diff:
%s

Please include a very brief Summary and Changes sections in your response.
Important: Do not include Markdown fencing in your response.`, diff)
	case 2:
		maxTokens = 2048 // high detail
		prompt = fmt.Sprintf(`Generate a detailed pull request description in Markdown format for the following git diff:
%s

Please include a comprehensive Summary and detailed Changes sections in your response.
Important: Do not include Markdown fencing in your response.`, diff)
	default:
		maxTokens = 1024 // medium detail
		prompt = fmt.Sprintf(`Generate a pull request description in Markdown format for the following git diff:
%s

Please include a concise Summary and Changes sections in your response.
Important: Do not include Markdown fencing in your response.`, diff)
	}

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
			MaxTokens: maxTokens,
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
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
