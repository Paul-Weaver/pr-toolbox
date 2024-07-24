package description

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// descriptionCmd represents the description command
var DescriptionCmd = &cobra.Command{
	Use:   "description",
	Short: "Create a PR description from a git diff",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		baseBranch, err := getBaseBranch()
		if err != nil {
			fmt.Println("Error determining base branch:", err)
			return
		}

		diff, err := getGitDiff(baseBranch)
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

func getBaseBranch() (string, error) {
	cmd := exec.Command("git", "branch", "-l", "master", "main")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	output := strings.TrimSpace(out.String())
	if strings.Contains(output, "master") {
		return "master", nil
	}
	if strings.Contains(output, "main") {
		return "main", nil
	}
	return "", fmt.Errorf("neither 'master' nor 'main' branches exist")
}

func getGitDiff(baseBranch string) (string, error) {
	cmd := exec.Command("git", "diff", baseBranch)
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
