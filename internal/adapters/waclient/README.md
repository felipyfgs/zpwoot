# WhatsApp Client Adapter

This package implements the WhatsApp client adapter for the zpwoot application.

## Files

- **`adapter.go`** - Implements the `output.WhatsAppClient` interface (main adapter)
- **`manager.go`** - Core WhatsApp client management and session handling
- **`events.go`** - WhatsApp event handling (messages, connections, etc.)
- **`messages.go`** - Message sending functionality (text, media, location, contact)
- **`qr.go`** - QR code management for session authentication
- **`types.go`** - Type definitions and data structures

## Architecture

This adapter follows Clean Architecture principles:

- **Implements**: `internal/core/ports/output.WhatsAppClient`
- **Dependencies**: Uses whatsmeow library for WhatsApp integration
- **Responsibility**: Handles all WhatsApp-related operations

## Usage

The adapter is initialized in the dependency injection container and used by the application layer through the output port interface.
