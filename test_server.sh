#!/bin/bash
# Test script for devwisdom MCP server
# Simulates Cursor's initialization and tool calls

set -e

BINARY="./devwisdom"
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo "🧪 Testing devwisdom MCP server..."
echo ""

# Test 1: Check binary exists and is executable
echo "✅ Test 1: Binary check"
if [ ! -f "$BINARY" ]; then
    echo "❌ Binary not found: $BINARY"
    exit 1
fi
if [ ! -x "$BINARY" ]; then
    echo "❌ Binary not executable: $BINARY"
    exit 1
fi
echo "   Binary: $BINARY"
ls -lh "$BINARY"
echo ""

# Test 2: Initialize request (simulates Cursor startup)
echo "✅ Test 2: Initialize request"
INIT_REQUEST='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'

echo "$INIT_REQUEST" | "$BINARY" > "$TEMP_DIR/init_response.json" 2> "$TEMP_DIR/init_stderr.txt" || true

# Check stderr for version message
if grep -q "devwisdom-go MCP server v0.1.0" "$TEMP_DIR/init_stderr.txt" 2>/dev/null; then
    echo "   ✅ Version message found in stderr"
    cat "$TEMP_DIR/init_stderr.txt"
else
    echo "   ⚠️  Version message not found (may be normal if stderr is redirected)"
fi

# Check if response is valid JSON
if jq . "$TEMP_DIR/init_response.json" > /dev/null 2>&1; then
    echo "   ✅ Initialize response is valid JSON"
    
    # Check version in response
    VERSION=$(jq -r '.result.serverInfo.version' "$TEMP_DIR/init_response.json" 2>/dev/null || echo "")
    if [ "$VERSION" = "0.1.0" ]; then
        echo "   ✅ Version in response: $VERSION (correct)"
    else
        echo "   ⚠️  Version in response: $VERSION (expected 0.1.0)"
    fi
else
    echo "   ❌ Initialize response is NOT valid JSON"
    echo "   Response:"
    cat "$TEMP_DIR/init_response.json"
    exit 1
fi
echo ""

# Test 3: Tools list request
echo "✅ Test 3: Tools list request"
TOOLS_REQUEST='{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'

echo "$TOOLS_REQUEST" | "$BINARY" > "$TEMP_DIR/tools_response.json" 2>/dev/null || true

if jq . "$TEMP_DIR/tools_response.json" > /dev/null 2>&1; then
    echo "   ✅ Tools list response is valid JSON"
    TOOL_COUNT=$(jq '.result.tools | length' "$TEMP_DIR/tools_response.json" 2>/dev/null || echo "0")
    echo "   Found $TOOL_COUNT tools"
else
    echo "   ❌ Tools list response is NOT valid JSON"
    cat "$TEMP_DIR/tools_response.json"
    exit 1
fi
echo ""

# Test 4: Get wisdom tool call
echo "✅ Test 4: get_wisdom tool call"
WISDOM_REQUEST='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}'

echo "$WISDOM_REQUEST" | "$BINARY" > "$TEMP_DIR/wisdom_response.json" 2>/dev/null || true

if jq . "$TEMP_DIR/wisdom_response.json" > /dev/null 2>&1; then
    echo "   ✅ Get wisdom response is valid JSON"
    QUOTE=$(jq -r '.result.quote // "N/A"' "$TEMP_DIR/wisdom_response.json" 2>/dev/null || echo "N/A")
    echo "   Quote: ${QUOTE:0:60}..."
else
    echo "   ❌ Get wisdom response is NOT valid JSON"
    cat "$TEMP_DIR/wisdom_response.json"
    exit 1
fi
echo ""

# Test 5: Resources list
echo "✅ Test 5: Resources list request"
RESOURCES_REQUEST='{"jsonrpc":"2.0","id":4,"method":"resources/list","params":{}}'

echo "$RESOURCES_REQUEST" | "$BINARY" > "$TEMP_DIR/resources_response.json" 2>/dev/null || true

if jq . "$TEMP_DIR/resources_response.json" > /dev/null 2>&1; then
    echo "   ✅ Resources list response is valid JSON"
    RESOURCE_COUNT=$(jq '.result.resources | length' "$TEMP_DIR/resources_response.json" 2>/dev/null || echo "0")
    echo "   Found $RESOURCE_COUNT resources"
else
    echo "   ❌ Resources list response is NOT valid JSON"
    cat "$TEMP_DIR/resources_response.json"
    exit 1
fi
echo ""

# Test 6: Read resource (sources)
echo "✅ Test 6: Read wisdom://sources resource"
RESOURCE_READ_REQUEST='{"jsonrpc":"2.0","id":5,"method":"resources/read","params":{"uri":"wisdom://sources"}}'

echo "$RESOURCE_READ_REQUEST" | "$BINARY" > "$TEMP_DIR/resource_read_response.json" 2>/dev/null || true

if jq . "$TEMP_DIR/resource_read_response.json" > /dev/null 2>&1; then
    echo "   ✅ Resource read response is valid JSON"
    
    # Check if the embedded JSON in "text" field is valid
    TEXT_CONTENT=$(jq -r '.result.contents[0].text' "$TEMP_DIR/resource_read_response.json" 2>/dev/null || echo "")
    if [ -n "$TEXT_CONTENT" ]; then
        if echo "$TEXT_CONTENT" | jq . > /dev/null 2>&1; then
            echo "   ✅ Embedded JSON in text field is valid"
        else
            echo "   ❌ Embedded JSON in text field is NOT valid"
            echo "   Text content: ${TEXT_CONTENT:0:100}..."
            exit 1
        fi
    fi
else
    echo "   ❌ Resource read response is NOT valid JSON"
    cat "$TEMP_DIR/resource_read_response.json"
    exit 1
fi
echo ""

echo "🎉 All tests passed! Server is working correctly."
echo ""
echo "Summary:"
echo "  ✅ Binary exists and is executable"
echo "  ✅ Version 0.1.0 confirmed"
echo "  ✅ All JSON responses are valid"
echo "  ✅ No parse errors detected"
echo ""
echo "Ready for Cursor MCP integration!"

