# firestoredata-to-struct

A simple Go module that converts Firestore event data into native Go structs.  
Built to simplify Firestore-triggered Cloud Functions in Go.

## Features

- Convert Firestore event data (`google.events.cloud.firestore.v1.DocumentEventData`) to Go structs
- Handles nested structs, maps, arrays/slices, and primitive types
- Supports null values, Firestore timestamps, and recursive conversion

## Installation

```bash
go get github.com/Lucy-Teknologi/firestoredata-to-struct
