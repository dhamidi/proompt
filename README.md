# Proompt

A CLI tool for managing LLM prompts with placeholder substitution.

## What is Proompt?

Proompt helps you organize and reuse prompts for Large Language Models (LLMs) with support for:

- **4-level prompt hierarchy**: directory, project, project-local, and user level prompts
- **Placeholder substitution**: Use `${VAR}` and `${VAR:-default}` syntax in your prompts
- **Interactive workflow**: Pick prompts, edit variables and template content via markdown with YAML frontmatter
- **Multiple storage locations**: Store prompts at different scopes for different use cases

## Installation

### From Source

```bash
go install github.com/dhamidi/proompt/cmd/proompt@latest
```

### Build Locally

```bash
git clone https://github.com/dhamidi/proompt.git
cd proompt
go build -o proompt ./cmd/proompt
```

## Quick Start

1. Create a prompts directory and add a prompt:
   ```bash
   mkdir prompts
   echo "Hello ${NAME:-World}!" > prompts/greeting.md
   ```

2. Use the interactive picker to select and process the prompt:
   ```bash
   proompt pick
   ```

3. Fill in the placeholder values in your editor and save to get the processed prompt.

## Commands

- `proompt list` - List all available prompts with their sources
- `proompt show <name>` - Display a specific prompt's content
- `proompt edit [name]` - Edit a prompt (uses picker if no name provided)
- `proompt rm [name]` - Remove a prompt (uses picker if no name provided)
- `proompt pick` - Interactive workflow: select prompt, fill placeholders, output result

## Prompt Hierarchy

Prompts are discovered in the following order (higher priority first):

1. **Directory level**: `./prompts/` (current directory)
2. **Project level**: `<project-root>/prompts/` (shared, committed to git)
3. **Project-local level**: `<project-root>/.git/info/prompts/` (local, git-ignored)
4. **User level**: `$XDG_CONFIG_HOME/proompt/prompts/`

## Environment Variables

- `EDITOR` - Text editor for prompt editing (default: `nano`)
- `PROOMPT_PICKER` - Selection picker command (default: `fzf`)
- `PROOMPT_COPY_COMMAND` - Copy to clipboard command (default: `pbcopy`)

## Placeholder Syntax

- `${VAR}` - Simple placeholder
- `${VAR:-default}` - Placeholder with default value
- `$$` - Escape sequence for literal `$`

## Example

Create a prompt file `prompts/code-review.md`:

```markdown
Review the following ${LANGUAGE:-JavaScript} code for:
- Code quality
- Security issues
- Performance concerns

Focus on ${FOCUS:-general best practices}.

Code:
${CODE}
```

When you run `proompt pick` and select this prompt, you'll get an editor with a markdown file containing:

```markdown
---
CODE: ""
FOCUS: general best practices
LANGUAGE: JavaScript
---
Review the following ${LANGUAGE:-JavaScript} code for:
- Code quality
- Security issues
- Performance concerns

Focus on ${FOCUS:-general best practices}.

Code:
${CODE}
```

Edit the variables in the YAML frontmatter (between the `---` lines) and/or modify the template content below. Save to get the processed prompt output.
