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
	AgreedToProcess    bool // true or false to see if a particualr message is agreed to be processed by AI 
}

func AIProcessHouseMessage(message string, BenchmarkMessage string) (*AgentHouseResponse, error) {

	var EditDistanceNormalized float64
	log.Printf("Processing house message: %s", utils.NormalizeText(message))
	log.Println("Benchmark message normalized = ", utils.NormalizeText(BenchmarkMessage))

	EditDistanceNormalized = utils.NormalizedLevenshteinDistance(
		utils.NormalizeText(message),
		utils.NormalizeText(BenchmarkMessage),
	)

	log.Println("Edit Distance = ", EditDistanceNormalized)
	var response AgentHouseResponse
	response.AiAddress = "123 AI Street"
	response.AiPrimaryContact = "AI Primary Contact"
	response.AiSecondaryContact = "AI Secondary Contact"
	response.AiMessage = "Processed message: " + message
	response.AgreedToProcess = true

	return &response, nil

}
