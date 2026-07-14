# Roadmap

Planned features and improvements for my-agent.

## In Progress

_None yet._

## Planned

### 1. REPL TUI

Build an interactive terminal UI for chatting with agents.

- Use [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- Features: streaming responses, tool call visibility, model switching, conversation history
- Replace current basic REPL in `cmd/my-agent/main.go`

### 2. Configurable Agent

Make the agent configurable via external files (YAML/TOML/JSON).

- Model selection, provider settings, system prompts
- Tool registration and enable/disable
- Config file locations: `~/.config/my-agent/config.yaml` or `./my-agent.yaml`

### 3. Mixture of Agents 

make a custom agent that call multiple llms with one main llm acting as the aggregator. self_moa and multi_moa

#### 3.1 Fallbacks/Dual Mode 

well, I dont know if those features are subsets of mixture of agents idea, but they're similar enough so I made it a subtitle.
fallback would define one llm as fallback, whenever the main cant continue doing their job, the fallback enter the scene. for example, we can use deepseek-v4-flash as main model and set mimo-v2.5 as fallback to image files, whenever the deepseek-v4-flash tries to read a image, it will fallback to mimo, or for some reason deepseek-v4-flash api is failing, so we can rely on mimo until deepseek-v4-flash goes back. 

dual mode is basically the same idea of mixture of agents, but instead of proposer and aggregator, the second model would act more like as a second opinion. even though I dont know exactly how this would work... well, I think better about this later 

### 4. subagents

a tool that can invoke other agents instances. ideally subagents can be configured with a file tool 

### 5. MCP

well, I want mcp as a builtin tool 

### 6. Agent skills 

well, I want to support to agent skills, but I want a better implementation than most agents do. I want to add a syste of skill hinter before each llm call the harness would inject a list of most recommended skills based on context(not only user prompt, because I also want this to happen between agent loop iteractions turn)

### 7. Agent Client Protocol 

### 8. Agent to Agent Protocol 

### 9. LSP builtin 

### 10. a custom DSL for orchestration of the agent itself 

this is more like a dream feature than a confirmed feature. I really want to someday implement this, but I guess that I need to have atleast a mature code agent that I can use reliable(currently I'm using pi agent instead of my own agent)

### 11. better observability/logging system 

why? I want to make the agent being able to analyze user usage to suggest creation of new hooks/skills/prompts 

### 12. prompts templates.

copy of how pi does that 

### 13. multiple provider support. 

### 14. permission system 

the user can create different permission profiles of tool usage 
basically, each tool have a category and the user can create profiles that can configure to always ask, always allow, always disallow, hide 
example(might be very different, but the idea is supposed to be represented): 
```json 
"read_only": {
    "read": "allow",
    "*": "ask"
}

"yolo": {
    "*": "allow"
}

"custom": {
    "*": "ask",
    "tools": [
       "bash": "allow"
    ]
}
```
that's the idea. the user can configure his own profiles, but 4 builtin profiles come builtin: read_only, always_ask, read_write, yolo 

### 15. the tui repl might evolve to a full tui 

well, why? I want to add support to session management panel, this would require more control of the terminal than a simple repl that lives in the normal buffer. only try to do that after we have a mature repl. they're are supposed to coexist for sometime 

### 16. well, if we achieve successfully the acp protocol, we can implement a web ui to our agent with support to artifact rendering(latex, mermaid, html)

### 17. plan mode 

well, I want to implement a plan mode which would basically use the read_only profile. the point here is the agent would be able to create full plan*.md files and actually track tasks with some kind of checkpoint system per turn iteration and keep the checklist up to date. 

currently I tried to do some kind of planning with pi, but I dont like the fact that I need to keep too much power in the hands of the model instead of the harness handling what needs to handle.

btw, after plan phase, we would have a implement mode to the implementation phase. yolo or read_write mode 

### . _TBD_

_Add more items as ideas develop._

## Ideas / Backlog

- Multiple provider support beyond Ollama (OpenAI, Anthropic, etc.)
- Persistent conversation sessions
- Plugin/tool system for extensibility
- Configurable output formats (plain, markdown, JSON)

---

_Last updated: 2026-07-13_
