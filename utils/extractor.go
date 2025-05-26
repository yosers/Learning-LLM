package utils

import (
	"strings"
)

func ExtractThinkingProcess(response string) string {
	startTag := "<think>"
	endTag := "</think>"

	startPos := strings.Index(response, startTag)
	if startPos == -1 {
		return response // No think tags found, return original response
	}

	endPos := strings.Index(response, endTag)
	if endPos == -1 {
		return response // No end tag found, return original response
	}

	// Remove the think section and any newlines before the actual response
	cleanedResponse := response[endPos+len(endTag):]
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	return cleanedResponse
}
