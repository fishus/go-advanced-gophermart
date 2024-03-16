package accrual

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestOrderAccrualStatusValidate(t *testing.T) {
	testCases := []struct {
		name    string
		status  OrderAccrualStatus
		wantErr bool
	}{
		{
			"Status Registered",
			OrderAccrualStatusRegistered,
			false,
		},
		{
			"Status Processing",
			OrderAccrualStatusProcessing,
			false,
		},
		{
			"Status Invalid",
			OrderAccrualStatusInvalid,
			false,
		},
		{
			"Status Processed",
			OrderAccrualStatusProcessed,
			false,
		},
		{
			"Wrong status",
			"test1234",
			true,
		},
		{
			"Undefined status",
			OrderAccrualStatusUndefined,
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
