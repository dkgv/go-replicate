package replicate

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestPredictionsService_Await(t *testing.T) {
	type destination struct {
		foo string
	}

	type args struct {
		ctx          context.Context
		predictionID string
	}
	tests := []struct {
		name     string
		args     args
		giveJSON string
		wantDest destination
		wantErr  error
	}{
		{
			name: "should return error if context is cancelled",
			args: args{
				ctx:          canceledContext(),
				predictionID: "prediction-id",
			},
			giveJSON: "{}",
			wantErr:  context.Canceled,
		},
		{
			name: "should return failed error if prediction failed",
			args: args{
				ctx:          context.Background(),
				predictionID: "prediction-id",
			},
			giveJSON: `{"status": "failed"}`,
			wantErr:  ErrPredictionFailed,
		},
		{
			name: "should return cancled error if prediction status is canceled",
			args: args{
				ctx:          context.Background(),
				predictionID: "prediction-id",
			},
			giveJSON: `{"status": "canceled"}`,
			wantErr:  ErrPredictionCanceled,
		},
		{
			name: "should unmarshal response into destination",
			args: args{
				ctx:          context.Background(),
				predictionID: "prediction-id",
			},
			giveJSON: `{"status": "succeeded", "output": {"foo": "bar"}}`,
			wantDest: destination{
				foo: "bar",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBaseURL, teardown := mockServer(
				endpoint{
					path: "/predictions/prediction-id",
					handler: func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(tt.giveJSON))
					},
				},
			)
			defer teardown()

			client := NewClient("token", WithBaseURL(mockBaseURL))

			err := client.Predictions.Await(tt.args.ctx, tt.args.predictionID, &tt.wantDest)
			if err != tt.wantErr {
				t.Errorf("PredictionsService.Await() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.wantDest, tt.wantDest) {
				t.Errorf("PredictionsService.Await() destination = %v, wantDest %v", tt.wantDest, tt.wantDest)
			}
		})
	}
}
