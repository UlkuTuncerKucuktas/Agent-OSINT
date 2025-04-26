package internal

import (
  "github.com/pontus-devoteam/agent-sdk-go/pkg/agent"
  openaipkg "github.com/pontus-devoteam/agent-sdk-go/pkg/model/providers/openai"
)


func NewOSINTAgent(apiKey string) (*agent.Agent, *openaipkg.Provider) {
  provider := openaipkg.NewProvider(apiKey)
  provider.SetDefaultModel("gpt-4.1-mini")

  ag := agent.NewAgent("OSINT Agent")
  ag.SetModelProvider(provider)
  ag.WithModel("gpt-4.1-mini")
  ag.SetSystemInstructions(
    `You are an OSINT agent. 
1. First call the "google_search" tool with user's name to retrieve top search results.
2. Review the list and choose which URLs look most relevant.
3. For each chosen URL, call the "fetch_page" tool to get its page text.
4. Synthesize all data into a concise OSINT report.
5. Automatically call the "generate_phishing_emails" tool to create phishing email drafts based on the OSINT report.
6. Present both the OSINT report and the phishing email drafts in your response.`,
  )

  ag.WithTools(
    NewMultiSearchTool(),  
    NewFetchPageTool(),     
    NewPhishingEmailTool(provider),
  )

  return ag, provider
}
