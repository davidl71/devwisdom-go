# MCP Integration Examples

This guide demonstrates how to integrate with the devwisdom-go MCP server using the Model Context Protocol (MCP) over JSON-RPC 2.0.

## Overview

The devwisdom-go MCP server implements:
- **Transport**: stdio (standard input/output)
- **Protocol**: JSON-RPC 2.0
- **Tools**: 5 tools for wisdom access
- **Resources**: 4 resources for metadata

## MCP Server Setup

### Configuration

Add to your MCP client configuration (e.g., `.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "devwisdom": {
      "command": "/path/to/devwisdom",
      "args": []
    }
  }
}
```

### Running the Server

The server runs in stdio mode and processes JSON-RPC 2.0 messages:

```bash
# Server reads from stdin, writes to stdout
./devwisdom < requests.jsonl
```

## Initialize Handshake

Every MCP session starts with an initialize request:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "example-client",
      "version": "1.0.0"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": {},
      "resources": {}
    },
    "serverInfo": {
      "name": "devwisdom-go",
      "version": "0.1.0"
    }
  }
}
```

## Tools

### 1. consult_advisor

Consult a wisdom advisor based on metric, tool, or stage.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "consult_advisor",
    "arguments": {
      "metric": "security",
      "score": 40.0,
      "context": "Working on improving project security"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "timestamp": "2026-01-06T20:00:00Z",
    "consultation_type": "advisor",
    "advisor": "pistis_sophia",
    "advisor_icon": "📜",
    "advisor_name": "Pistis Sophia",
    "rationale": "Your security score indicates foundational work is needed.",
    "metric": "security",
    "score_at_time": 40.0,
    "consultation_mode": "building",
    "mode_icon": "🌱",
    "mode_frequency": "milestones",
    "quote": "Security is not a destination, but a journey.",
    "quote_source": "pistis_sophia",
    "encouragement": "Every security improvement builds a stronger foundation."
  }
}
```

**Example with tool:**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "consult_advisor",
    "arguments": {
      "tool": "project_scorecard",
      "score": 75.0
    }
  }
}
```

**Example with stage:**
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "consult_advisor",
    "arguments": {
      "stage": "daily_checkin",
      "score": 65.0
    }
  }
}
```

### 2. get_wisdom

Get a wisdom quote based on project health score and optional source.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "get_wisdom",
    "arguments": {
      "score": 75.0,
      "source": "stoic"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "quote": "The impediment to action advances action.",
    "source": "stoic",
    "encouragement": "What stands in the way becomes the way.",
    "wisdom_source": "Stoic Philosophy",
    "wisdom_icon": "🏛️"
  }
}
```

**Random source:**
```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "get_wisdom",
    "arguments": {
      "score": 50.0,
      "source": "random"
    }
  }
}
```

### 3. get_daily_briefing

Get a daily wisdom briefing with quotes and guidance.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "get_daily_briefing",
    "arguments": {
      "score": 75.0
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "result": {
    "date": "2026-01-06",
    "score": 75.0,
    "source": "stoic",
    "quote": {
      "quote": "The impediment to action advances action.",
      "source": "stoic",
      "encouragement": "What stands in the way becomes the way."
    },
    "sources_count": 21,
    "advisors_count": 36
  }
}
```

### 4. get_consultation_log

Retrieve consultation log entries (requires Phase 5 logging).

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "tools/call",
  "params": {
    "name": "get_consultation_log",
    "arguments": {
      "days": 7
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 8,
  "result": [
    {
      "timestamp": "2026-01-05T10:00:00Z",
      "consultation_type": "advisor",
      "advisor": "pistis_sophia",
      "score_at_time": 40.0,
      "quote": "Security is not a destination, but a journey."
    }
  ]
}
```

### 5. export_for_podcast

Export consultations as podcast episodes (requires Phase 5 logging).

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "tools/call",
  "params": {
    "name": "export_for_podcast",
    "arguments": {
      "days": 7
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 9,
  "result": {
    "episodes": [
      {
        "date": "2026-01-05",
        "consultations": [...]
      }
    ],
    "days": 7
  }
}
```

## Resources

### 1. wisdom://sources

List all available wisdom sources.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 10,
  "method": "resources/read",
  "params": {
    "uri": "wisdom://sources"
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 10,
  "result": {
    "contents": [
      {
        "uri": "wisdom://sources",
        "mimeType": "application/json",
        "text": "[{\"id\":\"pistis_sophia\",\"name\":\"Pistis Sophia\",\"icon\":\"📜\",\"description\":\"Gnostic wisdom\"},...]"
      }
    ]
  }
}
```

### 2. wisdom://advisors

List all available advisors.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 11,
  "method": "resources/read",
  "params": {
    "uri": "wisdom://advisors"
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 11,
  "result": {
    "contents": [
      {
        "uri": "wisdom://advisors",
        "mimeType": "application/json",
        "text": "{\"metric_advisors\":[...],\"tool_advisors\":[...],\"stage_advisors\":[...]}"
      }
    ]
  }
}
```

### 3. wisdom://advisor/{id}

Get details for a specific advisor.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 12,
  "method": "resources/read",
  "params": {
    "uri": "wisdom://advisor/security"
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 12,
  "result": {
    "contents": [
      {
        "uri": "wisdom://advisor/security",
        "mimeType": "application/json",
        "text": "{\"id\":\"security\",\"type\":\"metric\",\"advisor\":\"pistis_sophia\",\"rationale\":\"Security requires foundational wisdom\",\"icon\":\"🔒\",\"helps_with\":\"Security practices and reviews\"}"
      }
    ]
  }
}
```

### 4. wisdom://consultations/{days}

Get consultation log entries (requires Phase 5 logging).

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 13,
  "method": "resources/read",
  "params": {
    "uri": "wisdom://consultations/7"
  }
}
```

## Error Handling

### Invalid Method

```json
{
  "jsonrpc": "2.0",
  "id": 99,
  "method": "invalid/method"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 99,
  "error": {
    "code": -32601,
    "message": "Method not found",
    "data": "invalid/method"
  }
}
```

### Invalid Parameters

```json
{
  "jsonrpc": "2.0",
  "id": 100,
  "method": "tools/call",
  "params": {
    "name": "consult_advisor",
    "arguments": {
      "metric": "security"
      // Missing required "score" parameter
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 100,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": "score is required for metric consultations"
  }
}
```

## Integration Examples

### Python Client Example

```python
import json
import subprocess
import sys

def call_mcp_tool(tool_name, arguments):
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
            "name": tool_name,
            "arguments": arguments
        }
    }
    
    process = subprocess.Popen(
        ["./devwisdom"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        text=True
    )
    
    process.stdin.write(json.dumps(request) + "\n")
    process.stdin.close()
    
    response = json.loads(process.stdout.readline())
    return response["result"]

# Get wisdom quote
quote = call_mcp_tool("get_wisdom", {
    "score": 75.0,
    "source": "stoic"
})
print(quote["quote"])
```

### Node.js Client Example

```javascript
const { spawn } = require('child_process');

function callMCPTool(toolName, arguments) {
  return new Promise((resolve, reject) => {
    const server = spawn('./devwisdom', []);
    
    const request = {
      jsonrpc: '2.0',
      id: 1,
      method: 'tools/call',
      params: {
        name: toolName,
        arguments: arguments
      }
    };
    
    server.stdin.write(JSON.stringify(request) + '\n');
    server.stdin.end();
    
    let response = '';
    server.stdout.on('data', (data) => {
      response += data.toString();
    });
    
    server.on('close', (code) => {
      if (code === 0) {
        const result = JSON.parse(response);
        resolve(result.result);
      } else {
        reject(new Error(`Server exited with code ${code}`));
      }
    });
  });
}

// Get wisdom quote
callMCPTool('get_wisdom', { score: 75.0, source: 'stoic' })
  .then(quote => console.log(quote.quote))
  .catch(err => console.error(err));
```

## Best Practices

1. **Initialize First**: Always send initialize request before other operations
2. **Error Handling**: Check for error responses and handle appropriately
3. **Connection Management**: Keep server process alive for multiple requests
4. **Resource Cleanup**: Properly close stdin/stdout when done
5. **JSON-RPC Compliance**: Follow JSON-RPC 2.0 specification strictly

## See Also

- [CLI Usage Examples](./cli_usage.md)
- [MCP Specification](https://modelcontextprotocol.io/)
- [Main README](../README.md)

