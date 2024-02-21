package repositories

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
)

func parseCreateTermsPayload(payload []entities.CreateTermDto, module models.Module) ([]models.Term, error) {
	var terms []models.Term
	for _, v := range payload {
		termModel, err := v.ToModel()
		if err != nil {
			return nil, err
		}
		termModel.Modules = []models.Module{module}
		terms = append(terms, termModel)
	}
	return terms, nil
}

func extractAssociatedTerms(terms []models.Term) []models.Term {
	var filteredTerms []models.Term

	for _, v := range terms {
		if len(v.Modules) <= 1 {
			filteredTerms = append(filteredTerms, v)
		}
	}
	return filteredTerms
}
