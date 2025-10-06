package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "zpwoot Support",
            "url": "https://github.com/your-org/zpwoot",
            "email": "support@zpwoot.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Get basic information about the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Service Information",
                "responses": {
                    "200": {
                        "description": "Service information",
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_http_handlers.InfoResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check if the service and database are healthy",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_http_handlers.HealthResponse"
                        }
                    },
                    "503": {
                        "description": "Service is unhealthy",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Creates a new WhatsApp session with the specified configuration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Create WhatsApp Session",
                "parameters": [
                    {
                        "description": "Session configuration",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateSessionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Session created successfully",
                        "schema": {
                            "$ref": "#/definitions/SessionResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body or validation error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Session already exists",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieves a list of all WhatsApp sessions with their current status (without QR codes)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "List WhatsApp Sessions",
                "responses": {
                    "200": {
                        "description": "List of sessions (without QR codes)",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/connect": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Connects a WhatsApp session. If already connected, returns current status with appropriate message.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Connect WhatsApp Session",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session connected successfully or already connected",
                        "schema": {
                            "$ref": "#/definitions/SessionStatusResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid session ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Connection error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Permanently deletes a WhatsApp session and all associated data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Delete WhatsApp Session",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session deleted successfully",
                        "schema": {
                            "$ref": "#/definitions/SessionResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid session ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Deletion error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/disconnect": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Disconnects an active WhatsApp session temporarily. If already disconnected, returns current status with appropriate message.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Disconnect WhatsApp Session",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session disconnected successfully or already disconnected",
                        "schema": {
                            "$ref": "#/definitions/SessionStatusResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid session ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Disconnection error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/info": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieves detailed information about a specific WhatsApp session (without QR code)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get WhatsApp Session",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session information (without QR code)",
                        "schema": {
                            "$ref": "#/definitions/SessionListInfo"
                        }
                    },
                    "400": {
                        "description": "Invalid session ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/logout": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Logs out a WhatsApp session permanently. Unlinks device from WhatsApp. Requires QR scan to reconnect.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Logout WhatsApp Session",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session logged out successfully",
                        "schema": {
                            "$ref": "#/definitions/SessionResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid session ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Logout error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/qr": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieves the QR code for WhatsApp session authentication. Scan with WhatsApp mobile app.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get QR Code",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "QR code data",
                        "schema": {
                            "$ref": "#/definitions/QRCodeResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid session ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Session already connected",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "QR code generation error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "APIResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "$ref": "#/definitions/ErrorInfo"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "timestamp": {
                    "type": "string",
                    "example": "2025-01-15T10:30:00Z"
                }
            }
        },
        "CreateSessionRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1,
                    "example": "my-session"
                },
                "qrCode": {
                    "type": "boolean",
                    "example": true
                },
                "settings": {
                    "$ref": "#/definitions/SessionSettings"
                }
            }
        },
        "ErrorInfo": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "validation_error"
                },
                "details": {
                    "type": "object",
                    "additionalProperties": true
                },
                "message": {
                    "type": "string",
                    "example": "Validation failed"
                }
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "validation_error"
                },
                "message": {
                    "type": "string",
                    "example": "name is required"
                }
            }
        },
        "ProxySettings": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean",
                    "example": true
                },
                "host": {
                    "type": "string",
                    "example": "proxy.example.com"
                },
                "pass": {
                    "type": "string",
                    "example": "proxyPass123"
                },
                "port": {
                    "type": "string",
                    "example": "8080"
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "http",
                        "https",
                        "socks5"
                    ],
                    "example": "http"
                },
                "user": {
                    "type": "string",
                    "example": "proxyUser123"
                }
            }
        },
        "QRCodeResponse": {
            "type": "object",
            "properties": {
                "expiresAt": {
                    "type": "string",
                    "example": "2025-01-15T10:35:00Z"
                },
                "qrCode": {
                    "type": "string",
                    "example": "2@abc123..."
                },
                "qrCodeBase64": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
                },
                "status": {
                    "type": "string",
                    "example": "generated"
                }
            }
        },
        "SessionListInfo": {
            "type": "object",
            "properties": {
                "connected": {
                    "type": "boolean",
                    "example": true
                },
                "connectedAt": {
                    "type": "string",
                    "example": "2025-01-15T10:32:00Z"
                },
                "createdAt": {
                    "type": "string",
                    "example": "2025-01-15T10:30:00Z"
                },
                "deviceJid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "lastSeen": {
                    "type": "string",
                    "example": "2025-01-15T10:35:00Z"
                },
                "name": {
                    "type": "string",
                    "example": "my-session"
                },
                "sessionId": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "settings": {
                    "$ref": "#/definitions/SessionSettings"
                },
                "status": {
                    "type": "string",
                    "example": "connected"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2025-01-15T10:35:00Z"
                }
            }
        },
        "SessionResponse": {
            "type": "object",
            "properties": {
                "connected": {
                    "type": "boolean",
                    "example": true
                },
                "connectedAt": {
                    "type": "string",
                    "example": "2025-01-15T10:32:00Z"
                },
                "createdAt": {
                    "type": "string",
                    "example": "2025-01-15T10:30:00Z"
                },
                "deviceJid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "lastSeen": {
                    "type": "string",
                    "example": "2025-01-15T10:35:00Z"
                },
                "name": {
                    "type": "string",
                    "example": "my-session"
                },
                "qrCode": {
                    "type": "string",
                    "example": "2@abc123..."
                },
                "qrCodeBase64": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgo..."
                },
                "qrCodeExpiresAt": {
                    "type": "string",
                    "example": "2025-01-15T10:35:00Z"
                },
                "sessionId": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "settings": {
                    "$ref": "#/definitions/SessionSettings"
                },
                "status": {
                    "type": "string",
                    "example": "connected"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2025-01-15T10:35:00Z"
                }
            }
        },
        "SessionSettings": {
            "type": "object",
            "properties": {
                "proxy": {
                    "$ref": "#/definitions/ProxySettings"
                },
                "webhook": {
                    "$ref": "#/definitions/WebhookSettings"
                }
            }
        },
        "SessionStatusResponse": {
            "type": "object",
            "properties": {
                "connected": {
                    "type": "boolean",
                    "example": true
                },
                "id": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "message": {
                    "type": "string",
                    "example": "Session is already connected"
                },
                "status": {
                    "type": "string",
                    "example": "connected"
                }
            }
        },
        "WebhookSettings": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean",
                    "example": true
                },
                "events": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Message",
                        "Receipt",
                        "Connected"
                    ]
                },
                "secret": {
                    "type": "string",
                    "example": "supersecrettoken123"
                },
                "url": {
                    "type": "string",
                    "example": "https://api.example.com/webhook"
                }
            }
        },
        "internal_adapters_http_handlers.HealthResponse": {
            "type": "object",
            "properties": {
                "service": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "internal_adapters_http_handlers.InfoResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "service": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "API key for authentication. Use the value from .env file (API_KEY variable). Send as 'Authorization: your-api-key' (without Bearer prefix)",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "WhatsApp session management operations",
            "name": "Sessions"
        },
        {
            "description": "Message sending and retrieval operations",
            "name": "Messages"
        },
        {
            "description": "Contact management operations",
            "name": "Contacts"
        },
        {
            "description": "Health check and system status",
            "name": "Health"
        }
    ]
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "zpwoot WhatsApp API",
	Description:      "A comprehensive WhatsApp Business API built with Go, following Clean Architecture principles.\nProvides endpoints for session management, messaging, contacts, groups, media handling, and integrations.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
