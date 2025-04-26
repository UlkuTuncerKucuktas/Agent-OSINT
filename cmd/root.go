// cmd/root.go
package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/UlkuTuncerKucuktas/Agent-OSINT/internal"
  "github.com/pontus-devoteam/agent-sdk-go/pkg/runner"
)

var rootCmd = &cobra.Command{
  Use:   "agent_osint [name]",
  Short: "Run an OSINT agent on a person’s name",
  Args:  cobra.ExactArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    name := args[0]
    key := os.Getenv("OPENAI_API_KEY")
    if key == "" {
      return fmt.Errorf("please set OPENAI_API_KEY")
    }

   
    ag, provider := internal.NewOSINTAgent(key)


    r := runner.NewRunner()
    r.WithDefaultProvider(provider)


    result, err := r.RunSync(ag, &runner.RunOptions{
      Input: name,
    })
    if err != nil {
      return err
    }

    fmt.Println("\n===== OSINT REPORT =====")
    fmt.Println(result.FinalOutput)
    return nil
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    os.Exit(1)
  }
}
