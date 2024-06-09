package service

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type GeminiService struct{}

func NewGeminiService() *GeminiService {
	return &GeminiService{}
}

func (g *GeminiService) AnalyzeChanges(
	ctx context.Context,
	diff string,
	userContext *string,
) (string, error) {
	client, err := genai.NewClient(
		ctx,
		option.WithAPIKey(viper.GetString("api.key")),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(
		ctx,
		genai.Text(
			fmt.Sprintf(
				`You're an automated AI that will only generate a conventional git commit message based on this diff changes:
%s

Follow this format:
<type>(<scope>): <description>

[optional body]

[optional footer(s)]

Types: docs, build(build related), ci, feat, fix(fix potential bug), perf, refactor, style(formatting), test, chore(internal changes), wip
NB:
User input this "%s"
Commits use a type, scope, and description. The type is a noun, scope is optional, and description is required.
Decide the commit type and scope(can be the filename) based on the diff and/or user input.
Description is a short summary of the changes and/or user input.
A longer body message may be provided after the description.
Each line in footer starts with a word token (use '-' instead of spaces), followed by ':' or '#' and a value.
Breaking changes are indicated by a ! in the type/scope prefix or as a footer.
Implementors treat units of information as case insensitive, except for BREAKING CHANGE which must be uppercase.`,
				diff,
				*userContext,
			),
		),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
