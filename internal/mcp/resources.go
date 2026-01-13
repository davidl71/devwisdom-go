package mcp

// GetResourceList returns the list of all available wisdom resources.
// This function eliminates duplication between server.go and sdk_adapter.go.
func GetResourceList() []Resource {
	return []Resource{
		{
			URI:         "wisdom://tools",
			Name:        "Available Tools",
			Description: "List all available MCP tools with descriptions and parameters",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://sources",
			Name:        "Wisdom Sources",
			Description: "List all available wisdom sources",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://advisors",
			Name:        "Wisdom Advisors",
			Description: "List all available advisors",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://advisor/{id}",
			Name:        "Advisor Details",
			Description: "Get details for a specific advisor",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://consultations/{days}",
			Name:        "Consultation Log",
			Description: "Get consultation log entries for the specified number of days",
			MimeType:    "application/json",
		},
	}
}
