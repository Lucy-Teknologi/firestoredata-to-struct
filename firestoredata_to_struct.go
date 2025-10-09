package firestoredata_to_struct

import (
	"context"
	"encoding/json"
	"fmt"

	// Replace with your actual import paths
	// "your/path/to/models"

	"github.com/bennovw/firestruct" // Dependency on the external package
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ConvertEventToStruct parses a CloudEvent, unmarshals the JSON payload into
// the firestruct model, and converts the "before" and "after" document states
// into Go structs of type T.
func ConvertEventToStruct[T any](ctx context.Context, e event.Event) (*T, *T, error) {
	dataBytes := e.Data()
	contentType := e.DataContentType()

	if len(dataBytes) == 0 {
		return nil, nil, fmt.Errorf("cloud event data is empty")
	}

	// 1. Handle "application/protobuf" by converting the binary data to JSON
	if contentType == "application/protobuf" {
		var protoData firestoredata.DocumentEventData

		// Unmarshal the raw Protobuf binary data
		if err := proto.Unmarshal(dataBytes, &protoData); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal binary protobuf data: %w", err)
		}

		// Marshal the Protobuf structure into JSON bytes (using protojson)
		// This creates the JSON payload that firestruct.FirestoreCloudEvent expects.
		jsonBytes, err := protojson.Marshal(&protoData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal protobuf data to JSON: %w", err)
		}

		dataBytes = jsonBytes // Use JSON data for the next step

	} else if contentType != "application/json" {
		// Only JSON or Protobuf is supported
		return nil, nil, fmt.Errorf("unexpected content type: %s (only application/protobuf or application/json is supported)", contentType)
	}

	var cloudEvent firestruct.FirestoreCloudEvent

	// 2. Unmarshal the CloudEvent data (which is JSON-encoded firestoredata.DocumentEventData)
	// Note: e.Data() returns the raw bytes. e.DataEncoded is an older field/concept.
	err := json.Unmarshal(dataBytes, &cloudEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal JSON payload into FirestoreCloudEvent: %w", err)
	}

	var before *T
	// 3. Handle 'before' document (OldValue)
	if cloudEvent.OldValue.Fields != nil {
		var tmp T

		// **FIX 1: Unwrap the Protobuf fields first.**
		unwrappedBefore, err := firestruct.UnwrapFirestoreFields(cloudEvent.OldValue.Fields)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unwrap 'before' document fields: %w", err)
		}

		// **FIX 2: Normalize numeric types based on struct T**
		// **FIX 3: Convert map → struct**
		unwrappedBefore = normalizeNumbers(unwrappedBefore).(map[string]any)
		if err := firestruct.DataTo(&tmp, unwrappedBefore); err != nil {
			return nil, nil, fmt.Errorf("failed to convert 'before' document data to struct: %w", err)
		}

		before = &tmp
	}

	var after *T
	// 4. Handle 'after' document (Value)
	if cloudEvent.Value.Fields != nil {
		var tmp T

		// Unwrap the Protobuf fields first.
		unwrappedAfter, err := firestruct.UnwrapFirestoreFields(cloudEvent.Value.Fields)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unwrap 'after' document fields: %w", err)
		}

		// Normalize numeric types based on struct T
		unwrappedAfter = normalizeNumbers(unwrappedAfter).(map[string]any)
		// Pass the unwrapped Go map (map[string]any) to DataTo.
		if err := firestruct.DataTo(&tmp, unwrappedAfter); err != nil {
			return nil, nil, fmt.Errorf("failed to convert 'after' document data to struct: %w", err)
		}
		after = &tmp
	}

	return before, after, nil
}

func normalizeNumbers(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			x[k] = normalizeNumbers(val)
		}
		return x
	case []any:
		for i, val := range x {
			x[i] = normalizeNumbers(val)
		}
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	default:
		return x
	}
}
