package firestoredata_to_struct

import (
	"context"
	"encoding/json"
	"fmt"

	// Replace with your actual import paths
	// "your/path/to/models"

	"github.com/bennovw/firestruct" // Dependency on the external package
	"github.com/cloudevents/sdk-go/v2/event"
)

// ConvertEventToStruct parses a CloudEvent, unmarshals the JSON payload into
// the firestruct model, and converts the "before" and "after" document states
// into Go structs of type T.
func ConvertEventToStruct[T any](ctx context.Context, e event.Event) (*T, *T, error) {
	// 1. Check content type (Note: firestruct expects JSON-encoded protobuf)
	if e.DataContentType() != "application/json" {
		// Real Firestore CE often uses application/json for the full event body
		// containing the firestoredata.DocumentEventData structure.
		return nil, nil, fmt.Errorf("unexpected content type: %s (expected application/json)", e.DataContentType())
	}

	var cloudEvent firestruct.FirestoreCloudEvent

	// 2. Unmarshal the CloudEvent data (which is JSON-encoded firestoredata.DocumentEventData)
	// Note: e.Data() returns the raw bytes. e.DataEncoded is an older field/concept.
	err := json.Unmarshal(e.Data(), &cloudEvent)
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

		// **FIX 2: Pass the unwrapped Go map (map[string]any) to DataTo.**
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

		// Pass the unwrapped Go map (map[string]any) to DataTo.
		if err := firestruct.DataTo(&tmp, unwrappedAfter); err != nil {
			return nil, nil, fmt.Errorf("failed to convert 'after' document data to struct: %w", err)
		}
		after = &tmp
	}

	return before, after, nil
}
