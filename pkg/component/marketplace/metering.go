package marketplace

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/ydataai/go-core/pkg/services/cloud"
)

type awsMeteringService struct {
	config      AWSMarketplaceConfiguration
	marketplace *marketplacemetering.MarketplaceMetering
}

func NewMarketplaceMetering(config AWSMarketplaceConfiguration, marketplace *marketplacemetering.MarketplaceMetering) cloud.MeteringService {
	return awsMeteringService{
		config:      config,
		marketplace: marketplace,
	}
}

func (s awsMeteringService) CreateUsageEvent(ctx context.Context, req cloud.UsageEventReq) (cloud.UsageEventRes, error) {
	// event
	event := &marketplacemetering.MeterUsageInput{
		ProductCode:    &s.config.ProductCode,
		UsageDimension: &req.DimensionID,
		UsageQuantity:  s.round(req.Quantity),
		Timestamp:      &req.StartAt,
	}
	// send
	output, err := s.marketplace.MeterUsageWithContext(ctx, event)
	if err != nil {
		return cloud.UsageEventRes{}, err
	}
	// result
	return cloud.UsageEventRes{
		UsageEventID: *output.MeteringRecordId,
		DimensionID:  req.DimensionID,
		Status:       output.String(),
	}, nil
}

func (s awsMeteringService) CreateUsageEventBatch(ctx context.Context, req cloud.UsageEventBatchReq) (cloud.UsageEventBatchRes, error) {
	output := []cloud.UsageEventRes{}
	for _, event := range req.Request {
		// events
		req := cloud.UsageEventReq{
			DimensionID: event.DimensionID,
			Quantity:    event.Quantity,
			StartAt:     event.StartAt,
		}
		// send
		res, err := s.CreateUsageEvent(ctx, req)
		if err != nil {
			return cloud.UsageEventBatchRes{}, err
		}
		output = append(output, res)
	}
	// result
	return cloud.UsageEventBatchRes{Result: output}, nil
}

// round
// returns the nearest integer, rounding half away from zero.
func (s awsMeteringService) round(quantity float32) *int64 {
	rounded := int64(math.Round(float64(quantity)))
	return &rounded
}
