package repositories

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paginate struct {
	limit int
	page  int
}

func newPaginate(limit int, page int) *paginate {
	return &paginate{limit: limit, page: page}
}

func (p *paginate) paginatedResult(db *gorm.DB) *gorm.DB {
	offset := (p.page - 1) * p.limit
	return db.Offset(offset).Limit(p.limit)
}

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

// getTermsToDelete returns terms from the module's terms slice that are not included in termsToReplace.
func getTermsToDelete(module models.Module, termsToReplace []models.Term) []models.Term {
	// Create a map to store term IDs from termsToReplace for efficient lookup
	termsToReplaceIDs := make(map[uuid.UUID]struct{})
	for _, term := range termsToReplace {
		termsToReplaceIDs[term.ID] = struct{}{}
	}

	// Filter out terms from the module's terms slice that are not included in termsToReplace
	var termsToDelete []models.Term
	for _, term := range module.Terms {
		// Check if the term's ID is not in the termsToReplaceIDs map
		if _, ok := termsToReplaceIDs[term.ID]; !ok {
			// Term is not in termsToReplace, so add it to the termsToDelete slice
			termsToDelete = append(termsToDelete, term)
		}
	}
	return termsToDelete
}

func (r *repo) pushToQueue(ctx context.Context, routingKey string, data interface{}) {
	if err := r.broker.PushToQueue(ctx, routingKey, data); err != nil {
		r.log(ctx, err.Error(), "error", "pushToQueue")
	}
}
