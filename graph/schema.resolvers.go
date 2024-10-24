package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"
	"fmt"
	"github.com/imrishuroy/legal-referral/graph/model"
)

// CreateFaq is the resolver for the createFAQ field.
func (r *mutationResolver) CreateFaq(ctx context.Context, input model.NewFaq) (*model.Faq, error) {
	panic(fmt.Errorf("not implemented: CreateFaq - createFAQ"))
}

// Faqs is the resolver for the faqs field.
func (r *queryResolver) Faqs(ctx context.Context) ([]*model.Faq, error) {

	faqs, err := r.Store.ListFAQs(ctx)

	if err != nil {
		return nil, err
	}

	var result []*model.Faq
	for _, faq := range faqs {
		result = append(result, &model.Faq{
			ID:       int(faq.FaqID),
			Question: faq.Question,
			Answer:   faq.Answer,
		})
	}

	return result, nil

}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct {
	*Resolver
	//Store db.Store

}
