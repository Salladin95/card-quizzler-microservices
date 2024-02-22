package repositories

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
)

func parseCreateTermsPayload(payload []entities.CreateTermDto, moduleID uuid.UUID) ([]models.Term, error) {
	var terms []models.Term
	for _, v := range payload {
		termModel, err := v.ToModel(moduleID)
		if err != nil {
			return nil, err
		}
		terms = append(terms, termModel)
	}
	return terms, nil
}
