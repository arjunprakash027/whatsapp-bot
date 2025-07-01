package agents

import (
//	"context"
	"log"
//	"time"

)

type AgentHouseResponse struct {
	AiAddress         string // Address of the AI that processed the message
	AiPrimaryContact  string // Primary contact of the AI that processed the message
	AiSecondaryContact string // Secondary contact of the AI that processed the message
	AiMessage         string // The message that was processed by the AI
}

func AIProcessHouseMessage(message string) (* AgentHouseResponse, error) {

	log.Printf("Processing house message: %s", message)

	var response AgentHouseResponse
	response.AiAddress = "123 AI Street"
	response.AiPrimaryContact = "AI Primary Contact"
	response.AiSecondaryContact = "AI Secondary Contact"
	response.AiMessage = "Processed message: " + message

	return &response, nil

}

