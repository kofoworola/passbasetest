package conversion

import (
	"context"
	"log"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kofoworola/passbasetest/integrations/fixer"
	pb "github.com/kofoworola/passbasetest/proto/v1/conversion"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedConversionServiceServer

	fixer *fixer.Handler
}

func New(fixer *fixer.Handler) *Service {
	return &Service{fixer: fixer}
}

func (s *Service) ConvertAmount(ctx context.Context, req *pb.ConvertAmountRequest) (*pb.ConvertAmountResponse, error) {
	if err := validation.ValidateStruct(
		req,
		validation.Field(&req.InputCurrency, validation.Required, validation.By(func(a interface{}) error {
			b := a.(*string)
			if *b != "EUR" {
				return validation.NewError("invalid_currency", "only EUR is supported as from")
			}
			return nil
		})),
		validation.Field(&req.OutputCurrency, validation.Required, is.Alpha),
		validation.Field(&req.Amount, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	amount, err := s.fixer.Convert(*req.InputCurrency, *req.OutputCurrency, *req.Amount)
	if err != nil {
		log.Printf("error converting currency: %v", err)
		return nil, status.Error(codes.Internal, "error converting")
	}
	return &pb.ConvertAmountResponse{
		Amount: &amount,
	}, nil
}
