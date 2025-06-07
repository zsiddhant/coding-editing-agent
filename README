# Coding Editing Agent

A simple command-line agent framework in Go that connects to a local Ollama LLM (Large Language Model) via HTTP, supports tool-calling, and can be extended with custom tool definitions.

## Features

- Interactive chat interface with a local LLM (Ollama).
- Extensible tool system: add new tools by defining their schema and function.
- Example tool: read the contents of a file in the workspace.
- Easy to add more tools for automation or code editing tasks.



## Requirements

Go 1.20+
Ollama running locally (default: http://localhost:11434)
(Optional) Additional Go packages for tool schemas (e.g., github.com/invopop/jsonschema)


## Usage

Start your Ollama server locally.

Build and run the agent:

```
ollama serve
``` 

Type your message at the prompt. The agent will respond using the LLM and can call tools as needed.


## Adding New Tools

- Define a new ToolDefinition in tool-definitions.go.
- Implement the tool’s function and input schema.
- Add the tool to the tools slice in main.go.

## Example : 

```
You: help me summarise the content from sample.txt
llama3.2: Here is a summary of the content from sample.txt:

The research paper "Calvin: Fast Distributed Transactions for Partitioned Database Systems" introduces Calvin, a system that enables high-throughput, ACID-compliant distributed transactions over partitioned and replicated data systems.

Key features of Calvin include:

* Deterministic execution, which reduces contention footprint and allows higher throughput
* Separation of concerns into three layers: sequencing, scheduling, and storage
* Support for asynchronous and Paxos-based synchronous replication
* Scalability on commodity hardware, with demonstrated performance on a 100-node cluster
* Fault tolerance through active replication and checkpointing

However, Calvin also has some trade-offs, including requiring advanced knowledge of read/write sets for transactions and handling dependent transactions using an Optimistic Lock Location Prediction (OLLP) scheme.

Let me know if you'd like me to summarize anything else!
```


##

Let me know if you want to add usage examples, API details, or contribution guidelines!