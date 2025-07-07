package agents

import (
	//	"context"
	"log"
	"whatsapp-bot/utils"
	//	"time"
)

type AgentHouseResponse struct {
	AiAddress          string // Address of the AI that processed the message
	AiPrimaryContact   string // Primary contact of the AI that processed the message
	AiSecondaryContact string // Secondary contact of the AI that processed the message
	AiMessage          string // The message that was processed by the AI
	AgreedToProcess    bool   // true or false to see if a particualr message is agreed to be processed by AI
}

func AIProcessHouseMessage(message string, BenchmarkMessage string) (*AgentHouseResponse, error) {

	var EditDistanceNormalized float64
	log.Printf("Processing house message: %s", message)

	EditDistanceNormalized = utils.NormalizedLevenshteinDistance(
		utils.NormalizeText(message),
		utils.NormalizeText(BenchmarkMessage),
	)

	log.Println("Edit Distance = ", EditDistanceNormalized)
	var response AgentHouseResponse

	if EditDistanceNormalized > 0.74 {
		response.AgreedToProcess = false
	} else {
		response.AiAddress = "123 AI Street"
		response.AiPrimaryContact = "AI Primary Contact"
		response.AiSecondaryContact = "AI Secondary Contact"
		response.AiMessage = "we are group of 3 students looking for 2 rooms in Dublin (even number parts of dublin) and our budget is 550-600 per person. If the house is still up, we would love to have a chat \n Reference: \n" + message
		response.AgreedToProcess = true
	}

	return &response, nil

}
