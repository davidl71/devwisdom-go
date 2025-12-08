# Testing devwisdom MCP Server

After restarting the server in Cursor, you can test it by:

## 1. Check Server Initialization

The server should print to stderr (visible in Cursor's MCP logs):
```
devwisdom-go MCP server v0.1.0 starting...
```

## 2. Test Tools

Try calling these MCP tools:

### `get_wisdom`
```json
{
  "score": 75.0,
  "source": "stoic"
}
```

### `consult_advisor`
```json
{
  "metric": "security",
  "score": 80.0,
  "context": "Testing server"
}
```

### `get_daily_briefing`
```json
{
  "score": 70.0
}
```

## 3. Test Resources

Try reading these resources:

- `wisdom://sources` - List all wisdom sources
- `wisdom://advisors` - List all advisors
- `wisdom://advisor/pistis_sophia` - Get specific advisor details

## 4. Verify No JSON Errors

After restart, you should NOT see:
- ❌ "} is not valid JSON"
- ❌ "Unexpected non-whitespace character after JSON"
- ❌ Any JSON parsing errors

## 5. Check Version in Initialize Response

The `initialize` response should include:
```json
{
  "serverInfo": {
    "name": "devwisdom",
    "version": "0.1.0"
  }
}
```

## Troubleshooting

If you still see errors:
1. Verify the binary path in `~/.cursor/mcp.json` is correct
2. Check that the binary is executable: `chmod +x /Users/davidl/Projects/devwisdom-go/devwisdom`
3. Check Cursor's MCP server logs for the version message
4. Try rebuilding: `cd /Users/davidl/Projects/devwisdom-go && go build -o devwisdom cmd/server/main.go`

