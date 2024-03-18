package anytrust

import (
	"context"

	"github.com/eniac-x-labs/anytrustDA/das"
	"github.com/eniac-x-labs/anytrustDA/util/signature"
)

type AnytrustDA struct {
	das.DataAvailabilityServiceWriter
	das.DataAvailabilityServiceReader
	*das.LifecycleManager
}

func NewAnytrustDA(ctx context.Context, daConfig *das.DataAvailabilityConfig, dataSigner signature.DataSignerFunc) (*AnytrustDA, error) {
	daWriter, daReader, lifeManager, err := das.CreateAggregatorComponents(ctx, daConfig, dataSigner)
	if err != nil {
		return nil, err
	}
	return &AnytrustDA{
		DataAvailabilityServiceWriter: daWriter,
		DataAvailabilityServiceReader: daReader,
		LifecycleManager:              lifeManager,
	}, nil
}
