# firestoredata-to-struct

A simple Go module that converts Firestore event data into native Go structs.  
Built to simplify Firestore-triggered Cloud Functions in Go.

## Features

- Convert Firestore event data (`google.events.cloud.firestore.v1.DocumentEventData`) to Go structs
- Handles nested structs, maps, arrays/slices, and primitive types
- Supports null values, Firestore timestamps, and recursive conversion

## Installation

```go
go get github.com/Lucy-Teknologi/firestoredata-to-struct
```

## Installation

```go
import (
	"context"
	"fmt"

	"github.com/Lucy-Teknologi/firestoredata-to-struct"
	"google.golang.org/genproto/googleapis/events/cloud/firestore/v1"
)

type MyModel struct {
	ID    string   `firestore:"id"`
	Name  string   `firestore:"name"`
}

func HandleFirestoreEvent(ctx context.Context, event event.Event) error {
	before, after, err := firestoredata_to_struct.ConvertEventToStruct[MyModel](ctx, event)
	if err != nil {
		return fmt.Errorf("failed to convert Firestore event data: %w", err)
	}

	// Use before and after as strongly-typed Go structs
	fmt.Println("Before:", before)
	fmt.Println("After:", after)

	return nil
}

```
