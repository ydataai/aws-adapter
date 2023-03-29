package marketplace

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/ydataai/go-core/pkg/metering"
)

type awsMeteringService struct {
	config      AWSMarketplaceConfiguration
	marketplace *marketplacemetering.MarketplaceMetering
}

func NewMarketplaceMetering(
	config AWSMarketplaceConfiguration, marketplace *marketplacemetering.MarketplaceMetering,
) metering.Client {
	return awsMeteringService{
		config:      config,
		marketplace: marketplace,
	}
}

func (s awsMeteringService) CreateUsageEvent(
	ctx context.Context, req metering.UsageEvent,
) (metering.UsageEventResponse, error) {
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
		return metering.UsageEventResponse{}, err
	}
	// result
	return metering.UsageEventResponse{
		UsageEventID: *output.MeteringRecordId,
		DimensionID:  req.DimensionID,
		Status:       output.String(),
	}, nil
}

func (s awsMeteringService) CreateUsageEventBatch(
	ctx context.Context, req metering.UsageEventBatch,
) (metering.UsageEventBatchResponse, error) {
	output := []metering.UsageEventResponse{}
	for _, event := range req.Events {
		// events
		req := metering.UsageEvent{
			DimensionID: event.DimensionID,
			Quantity:    event.Quantity,
			StartAt:     event.StartAt,
		}
		// send
		res, err := s.CreateUsageEvent(ctx, req)
		if err != nil {
			return metering.UsageEventBatchResponse{}, err
		}
		output = append(output, res)
	}
	// result
	return metering.UsageEventBatchResponse{Result: output}, nil
}

// round
// returns the nearest integer, rounding half away from zero.
func (s awsMeteringService) round(quantity float32) *int64 {
	rounded := int64(math.Round(float64(quantity)))
	return &rounded
}
