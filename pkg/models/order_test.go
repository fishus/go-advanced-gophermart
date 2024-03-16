package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatusValidate(t *testing.T) {
	testCases := []struct {
		name    string
		status  OrderStatus
		wantErr bool
	}{
		{
			"Status Processing",
			OrderStatusProcessing,
			false,
		},
		{
			"Status Invalid",
			OrderStatusInvalid,
			false,
		},
		{
			"Status Processed",
			OrderStatusProcessed,
			false,
		},
		{
			"Status New",
			OrderStatusNew,
			false,
		},
		{
			"Wrong status",
			"test1234",
			true,
		},
		{
			"Undefined status",
			OrderStatusUndefined,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.status.Validate()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
