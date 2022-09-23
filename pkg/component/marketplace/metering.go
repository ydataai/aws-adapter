package marketplace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	qty, err := s.convert(req.Quantity)
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
	quantity := int64(1)
	customer, err := s.marketplace.ResolveCustomerWithContext(ctx, &marketplacemetering.ResolveCustomerInput{
		RegistrationToken: &s.config.RegistrationToken,
	})
	if err != nil {
		return cloud.UsageEventBatchRes{}, err
	}
	// events
	events := []*marketplacemetering.UsageRecord{}
	for _, event := range req.Request {
		qty, err := s.convert(event.Quantity)
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
					AllocatedUsageQuantity: &quantity,
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

// convert
// %f = 6 decimal places
// remove decimal point = quantity * 1000000
func (s awsMeteringService) convert(quantity float32) (*int64, error) {
	qntstr := strings.Replace(fmt.Sprintf("%f", quantity), ".", "", 1)
	qnt, err := strconv.ParseInt(qntstr, 10, 64)
	if err != nil {
		return nil, err
	}
	return &qnt, nil
}
