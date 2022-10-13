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
	customer, err := s.marketplace.ResolveCustomerWithContext(ctx, &marketplacemetering.ResolveCustomerInput{
		RegistrationToken: &s.config.RegistrationToken,
	})
	if err != nil {
		return cloud.UsageEventRes{}, err
	}
	// event
	qty, err := s.round(req.Quantity)
	if err != nil {
		return cloud.UsageEventRes{}, err
	}
	event := &marketplacemetering.MeterUsageInput{
		ProductCode:    customer.ProductCode,
		UsageDimension: &req.DimensionID,
		UsageQuantity:  qty,
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
	customer, err := s.marketplace.ResolveCustomerWithContext(ctx, &marketplacemetering.ResolveCustomerInput{
		RegistrationToken: &s.config.RegistrationToken,
	})
	if err != nil {
		return cloud.UsageEventBatchRes{}, err
	}
	// events
	events := []*marketplacemetering.UsageRecord{}
	for _, event := range req.Request {
		qty, err := s.round(event.Quantity)
		if err != nil {
			return cloud.UsageEventBatchRes{}, err
		}
		events = append(events, &marketplacemetering.UsageRecord{
			CustomerIdentifier: customer.CustomerIdentifier,
			Dimension:          &event.DimensionID,
			Quantity:           qty,
			Timestamp:          &event.StartAt,
			UsageAllocations: []*marketplacemetering.UsageAllocation{
				{
					AllocatedUsageQuantity: qty,
					Tags:                   []*marketplacemetering.Tag{},
				},
			},
		})
	}
	// send
	output, err := s.marketplace.BatchMeterUsageWithContext(ctx, &marketplacemetering.BatchMeterUsageInput{
		ProductCode:  customer.ProductCode,
		UsageRecords: events,
	})
	if err != nil {
		return cloud.UsageEventBatchRes{}, err
	}
	// result
	results := []cloud.UsageEventRes{}
	for _, result := range output.Results {
		results = append(results, cloud.UsageEventRes{
			UsageEventID: *result.MeteringRecordId,
			DimensionID:  *result.UsageRecord.Dimension,
			Status:       result.String(),
		})
	}
	return cloud.UsageEventBatchRes{Result: results}, nil
}

// round
// returns the nearest integer, rounding half away from zero.
func (s awsMeteringService) round(quantity float32) (*int64, error) {
	rounded := int64(math.Round(float64(quantity)))
	return &rounded, nil
}
