# mermaid snippets

## Design
```
graph TD
    C[Client] -->|set task| S(scheduler API)
    C[Client] -->|get task| S
    
    W(worker) -->|claim tasks| WA(worker API)
    W(worker) -->|suceed| WA(worker API)
    W(worker) -->|fail| WA(worker API)
    
    subgraph scheduler
    S --> D[(Storage)]
    WA --> D
    end
```