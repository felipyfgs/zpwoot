#!/bin/bash

# Replace validation calls with basic validation comments
sed -i 's/if err := h\.GetValidator()\.ValidateStruct(&req); err != nil {/\/\/ Basic validation removed/g' internal/adapters/http/handler/message_handler.go

# Remove the error handling lines that follow
sed -i '/h\.GetWriter()\.WriteBadRequest(w, "Validation failed", err\.Error())/d' internal/adapters/http/handler/message_handler.go

# Remove the return statements that are now orphaned
sed -i '/\/\/ Basic validation removed/{n;/return/d;}' internal/adapters/http/handler/message_handler.go
sed -i '/\/\/ Basic validation removed/{n;/}/d;}' internal/adapters/http/handler/message_handler.go
