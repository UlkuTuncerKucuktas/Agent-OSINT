package internal

import (
    "context"
    "fmt"

    "github.com/pontus-devoteam/agent-sdk-go/pkg/tool"
    openaipkg "github.com/pontus-devoteam/agent-sdk-go/pkg/model/providers/openai"
    modelPkg "github.com/pontus-devoteam/agent-sdk-go/pkg/model"
)

// PhishingEmailTool generates phishing emails based on an OSINT summary.
type PhishingEmailTool struct {
    provider *openaipkg.Provider
}

// NewPhishingEmailTool constructs the tool using your OpenAI provider.
func NewPhishingEmailTool(provider *openaipkg.Provider) tool.Tool {
    return &PhishingEmailTool{provider: provider}
}

func (t *PhishingEmailTool) GetName() string {
    return "generate_phishing_emails"
}

func (t *PhishingEmailTool) GetDescription() string {
    return "Generate 3 persuasive phishing email drafts targeting the person described in the given OSINT summary."
}

func (t *PhishingEmailTool) GetParametersSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "summary": map[string]interface{}{
                "type":        "string",
                "description": "OSINT report summary about the target",
            },
        },
        "required": []string{"summary"},
    }
}

func (t *PhishingEmailTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    raw, ok := params["summary"]
    if !ok {
        return nil, fmt.Errorf("missing required parameter: summary")
    }
    summary, ok := raw.(string)
    if !ok {
        return nil, fmt.Errorf("parameter 'summary' must be a string")
    }

    // Build the prompt using the OSINT summary.
    prompt := fmt.Sprintf(`
You are a social-engineering specialist. Given this OSINT report about a target:
%s

Please draft 3 different, realistic-sounding phishing emails aimed at this person. Each email should:
- Appear to come from someone they know or a trusted organization.
- Reference personal details from the summary.
- Use persuasive language to trick them into clicking a link or opening an attachment.

Format your output clearly, separating each email with a header like "Email 1:", "Email 2:", etc.`, summary)

    // 1) Get the default model from the provider (which you set via SetDefaultModel).
    mdl, err := t.provider.GetModel("")
    if err != nil {
        return nil, fmt.Errorf("failed to get model: %w", err)
    }

    // 2) Construct a Request for the model.
    req := &modelPkg.Request{
        SystemInstructions: "",
        Input:              prompt,
    }

    // 3) Send the request and receive a response.
    res, err := mdl.GetResponse(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("model response error: %w", err)
    }

    // 4) Return the generated phishing emails.
    return res.Content, nil
}
