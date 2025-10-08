package firestoredata_to_struct_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	firestoredata_to_struct "github.com/Lucy-Teknologi/firestoredata-to-struct"
	"github.com/Lucy-Teknologi/firestoredata-to-struct/models"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// dataToProtobufValue converts a single Go interface{} value into a *firestoredata.Value.
func dataToProtobufValue(v interface{}) (*firestoredata.Value, error) {
	if v == nil {
		// Explicitly return a null value
		return &firestoredata.Value{
			ValueType: &firestoredata.Value_NullValue{},
		}, nil
	}

	switch val := v.(type) {
	case string:
		return &firestoredata.Value{ValueType: &firestoredata.Value_StringValue{StringValue: val}}, nil
	case bool:
		return &firestoredata.Value{ValueType: &firestoredata.Value_BooleanValue{BooleanValue: val}}, nil
	case int: // Firestore integers are int64
		return &firestoredata.Value{ValueType: &firestoredata.Value_IntegerValue{IntegerValue: int64(val)}}, nil
	case int64:
		return &firestoredata.Value{ValueType: &firestoredata.Value_IntegerValue{IntegerValue: val}}, nil
	case float64:
		return &firestoredata.Value{ValueType: &firestoredata.Value_DoubleValue{DoubleValue: val}}, nil
	case time.Time:
		// Firestore stores time as a TimestampValue
		return &firestoredata.Value{
			ValueType: &firestoredata.Value_TimestampValue{
				TimestampValue: timestamppb.New(val),
			},
		}, nil
	case map[string]interface{}:
		// This is a MapValue (used for nested objects)
		fields := make(map[string]*firestoredata.Value)
		for k, v := range val {
			pv, err := dataToProtobufValue(v)
			if err != nil {
				return nil, err
			}
			fields[k] = pv
		}
		return &firestoredata.Value{
			ValueType: &firestoredata.Value_MapValue{
				MapValue: &firestoredata.MapValue{Fields: fields},
			},
		}, nil
	case []any:
		// This is an ArrayValue
		arrayValues := make([]*firestoredata.Value, len(val))
		for i, item := range val {
			pv, err := dataToProtobufValue(item)
			if err != nil {
				return nil, err
			}
			arrayValues[i] = pv
		}
		return &firestoredata.Value{
			ValueType: &firestoredata.Value_ArrayValue{
				ArrayValue: &firestoredata.ArrayValue{Values: arrayValues},
			},
		}, nil
	// You would need to add support for []byte, []interface{}, latlng, etc. here.

	default:
		return nil, fmt.Errorf("unsupported type for Protobuf conversion: %T", v)
	}
}

// simulateRawDataToProtobufFields converts a Go map into the Protobuf Fields map.
func simulateRawDataToProtobufFields(data map[string]interface{}) (map[string]*firestoredata.Value, error) {
	fields := make(map[string]*firestoredata.Value)
	for key, val := range data {
		pv, err := dataToProtobufValue(val)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", key, err)
		}
		fields[key] = pv
	}
	return fields, nil
}

func TestFirestoreDataToStructOnOrderingModel(t *testing.T) {
	ctx := t.Context()

	fs, err := firestore.NewClient(ctx, "lucy-cashier-dev", option.WithCredentialsFile("service-account.json"))
	if err != nil {
		t.Fatalf("failed to create Firestore client: %v", err)
	}
	defer fs.Close()

	doc, err := fs.Collection("cancelled_orders").Doc("cancel-123").Get(ctx)
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}

	fmt.Println("Document Data:", doc.Data())

	// --- 2. Convert Raw Data to Protobuf Fields ---
	pbFieldsBefore, err := simulateRawDataToProtobufFields(doc.Data())
	if err != nil {
		t.Fatalf("Failed to convert raw data (before) to Protobuf fields: %v", err)
	}
	newDoc := doc.Data()
	newDoc["approval"] = map[string]interface{}{
		"by": "admin-user",
	}
	pbFieldsAfter, err := simulateRawDataToProtobufFields(newDoc)
	if err != nil {
		t.Fatalf("Failed to convert raw data (after) to Protobuf fields: %v", err)
	}

	// --- 3. Construct the Full Protobuf Event Data ---
	pbDocBefore := &firestoredata.Document{
		Name:   "projects/p/databases/d/documents/cancelled_orders/cancel-123",
		Fields: pbFieldsBefore,
	}
	pbDocAfter := &firestoredata.Document{
		Name:   "projects/p/databases/d/documents/cancelled_orders/cancel-123",
		Fields: pbFieldsAfter,
	}

	eventData := &firestoredata.DocumentEventData{
		OldValue: pbDocBefore, // Populated for an UPDATE or DELETE
		Value:    pbDocAfter,  // Populated for a CREATE or UPDATE
	}

	// --- 4. Marshal the Protobuf Event Data into JSON ---
	// The firestruct package expects the CloudEvent data to be JSON.
	binPayload, err := proto.Marshal(eventData)
	if err != nil {
		t.Fatalf("Failed to marshal event data to JSON: %v", err)
	}

	// --- 5. Create the CloudEvent with the JSON Payload ---
	e := event.New()
	e.SetSource("test-source")
	// e.SetSource("projects/lucy-cashier-dev/databases/(default)/documents/orders/cancelled_orders")
	e.SetType("google.cloud.firestore.document.v1.written")

	// THIS IS THE CRITICAL FIX: Set the marshaled *JSON bytes*
	if err := e.SetData("application/protobuf", binPayload); err != nil {
		t.Fatalf("Failed to set event data: %v", err) // Should now succeed
	}

	before, after, err := firestoredata_to_struct.
		ConvertEventToStruct[models.CancelledOrder](
		ctx,
		e,
	)
	if err != nil {
		t.Fatalf("failed to convert Firestore document to struct: %v", err)
	}

	// Test the before and after states
	if before != nil {
		b, _ := json.MarshalIndent(before, "", "\t")
		t.Logf("Before: %s", b)
	}
	if after != nil {
		b, _ := json.MarshalIndent(after, "", "\t")
		t.Logf("After: %s", b)
	}
}
