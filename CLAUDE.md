# AGENTS.md

## General Behavior

- Think before acting.
- Create a complete implementation plan before making any changes.
- Execute the plan in as few edits as possible.
- Prefer one large edit over many small edits.
- Do not narrate intermediate reasoning or implementation steps.
- Only provide a brief summary after all work is complete.

## File Access

- If I specify one or more file paths, only work with those files.
- Do not search the repository unless I explicitly ask.
- Read each file at most once unless an edit fails or additional context is required.
- Do not re-read a file immediately after editing it.
- Do not inspect unrelated files.

## Searching

- Avoid repository-wide searches.
- Do not grep, find, or scan the project unless necessary.
- Never search the project when the target file has already been provided.

## Editing

- Batch all related changes into a single edit.
- Avoid incremental edits (+5, +20, +50, etc.).
- Do not pause after each modification.
- Finish the complete implementation before responding.

## Tool Usage

- Never call MCP tools (Jira, Slack, Gmail, Drive, Calendar, Figma, Box, Canva, SPAN UI, etc.) unless the current task explicitly requires them.
- Do not use external tools for information that already exists in the repository.

## Token Efficiency

- Minimize tool calls.
- Minimize file reads.
- Minimize repeated searches.
- Prefer using existing context over re-reading files.
- Keep responses concise.