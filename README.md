# Agent-OSINT

An OSINT (Open Source Intelligence) tool that performs multi-engine web searches and generates phishing email drafts based on the gathered information.

## Prerequisites

- Go 1.23 or higher
- OpenAI API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/UlkuTuncerKucuktas/Agent-OSINT.git
cd Agent-OSINT
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your OpenAI API key:
```bash
export OPENAI_API_KEY='your-api-key-here'
```

## Usage

Run the tool by providing a person's name as an argument:

```bash
go run main.go "John Doe"
```

The tool will:
1. Perform concurrent searches across multiple search engines
2. Analyze the search results
3. Generate a comprehensive OSINT report
4. Create phishing email drafts based on the gathered information

## Output

The tool will output:
- A detailed OSINT report about the target
- Three different phishing email drafts based on the gathered information

## Security Note

This tool is intended for educational and security testing purposes only. Always ensure you have proper authorization before using this tool against any target.
