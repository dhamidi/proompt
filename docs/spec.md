# Overview

Proompt manages LLM prompts with placeholders.

Prompts can exist at three levels:

- the directory level (current directory),
- the project level (identified by .git or a prompts/ folder),
  - shared prompts are committed to .git,
  - local prompts are automatically added to .git/info/exclude
- the user level ($XDG_CONFIG_HOME/proompt/)

Prompts are stored in markdown files, in a directory called prompts.

They can contain placeholders matching the syntax of shell variable substitution, a double $$ escapes variable substitution.

# CLI

The following commands are supported:

- `proompt list` lists all prompts together with their source (directory, project (shared or local), user),
- `proompt show <name>` shows a prompt by name (filename)
- `proompt pick` invokes the picker (default fzf, can be configured through `PROOMPT_PICKER`),
- `proompt edit <name>` invokes `$EDITOR` on the given prompt (`--directory`, `--project`, `--project-local`, `--user` specify where to create it)
  - if no name is provided, the picker is invoked to pick an existing prompt
- `proompt rm <name>` removes the given prompt
  - if no name is provided, the picker is invoked to pick an existing prompt
  - attempting to remove a non-existent prompt is not an error but expected behavior

## proompt pick

This is the heart of the application:

1. It invokes the picker,
2. Let's the user select a given prompt,
3. Then parses the prompt to find all placeholders,
4. Then prepares a temporary file with the placeholders and their default values, invoking `$EDITOR` on that file.
5. Once the editor exits and the file is not empty, replaces all placeholders in the prompt and writes it out to stdout.

Given the following prompt:

```
Study ${STEPS:-docs/status.md} and pick the highest-priority step,

then ${ACTION} it.
```

The temporary file should look like this:

```

STEPS=docs/status.md
ACTION=
### Lines starting with # are ignored
# Save empty file to abort
### Full prompt:
# Study ${STEPS:-docs/status.md} and pick the highest-priority step,
# 
# then ${ACTION} it.
```
