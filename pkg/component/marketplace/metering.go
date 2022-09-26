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
	// event
	qty, err := s.convert(req.Quantity)
	if err != nil {
		return cloud.UsageEventRes{}, err
	}
	event := &marketplacemetering.MeterUsageInput{
		ProductCode:    &s.config.ProductCode,
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
