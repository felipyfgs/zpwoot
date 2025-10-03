package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/zpwoot/zpwoot/blob/main/LICENSE",
        "contact": {
            "name": "ZPWoot API Support",
            "url": "https://github.com/zpwoot/zpwoot",
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
        "/sessions/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new WhatsApp session with optional proxy configuration. If qrCode is true, returns QR code immediately for connection.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Create new session",
                "parameters": [
                    {
                        "description": "Session creation request with optional qrCode flag",
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
                        "description": "Session created successfully. If qrCode was true, includes QR code data.",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/CreateSessionResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
                        "description": "Internal Server Error",
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
                "description": "Get a list of all WhatsApp sessions with optional filtering",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "List sessions",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "Filter by connection status",
                        "name": "isConnected",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by device JID",
                        "name": "deviceJid",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of sessions to return (default: 20)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of sessions to skip (default: 0)",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sessions retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/ListSessionsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/stats": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get statistics about all sessions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get session statistics",
                "responses": {
                    "200": {
                        "description": "Session statistics retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SessionStatsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/chatwoot/find": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Find the current Chatwoot configuration for the session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chatwoot"
                ],
                "summary": "Find Chatwoot configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/chatwoot/set": {
            "post": {
                "description": "Create a new Chatwoot configuration for the session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chatwoot"
                ],
                "summary": "Create Chatwoot configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
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
                "description": "Connect a WhatsApp session to start receiving messages. Automatically returns QR code (both string and base64 image) if device needs to be paired. If session is already connected, returns confirmation message.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Connect session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session connection initiated successfully with QR code if needed, or confirmation if already connected",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/ConnectSessionResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "List contacts with pagination and filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "List contacts",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit (default: 50, max: 100)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset (default: 0)",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search term",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.ListContactsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/all": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all contacts without pagination",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Get all contacts",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/internal_adapters_server_handler.ContactInfo"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/avatar": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get profile picture of a contact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Get profile picture",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Contact JID",
                        "name": "jid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.GetProfilePictureResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/business": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get business profile of a contact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Get business profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Contact JID",
                        "name": "jid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.BusinessProfileResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/check": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Check if phone numbers are registered on WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Check WhatsApp numbers",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Phone numbers to check",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_server_handler.CheckWhatsAppRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.CheckWhatsAppResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/detailed-info": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get detailed information about WhatsApp users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Get detailed user info (batch)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User JIDs to get detailed info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_server_handler.GetUserInfoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.GetUserInfoResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/info": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get information about WhatsApp users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Get user info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User JIDs to get info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_server_handler.GetUserInfoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.GetUserInfoResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/is-on-whatsapp": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Check if multiple phone numbers are registered on WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Check if numbers are on WhatsApp (batch)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Phone numbers to check",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_server_handler.CheckWhatsAppRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.CheckWhatsAppResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/profile-picture-info": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get profile picture information of a contact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Get profile picture info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Contact JID",
                        "name": "jid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.GetProfilePictureResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/contacts/sync": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Sync contacts from the device",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Sync contacts",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapters_server_handler.SyncContactsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
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
                "description": "Delete a WhatsApp session and all associated data",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Delete session",
                "parameters": [
                    {
                        "type": "string",
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
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
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
                "description": "Disconnect from WhatsApp session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Disconnect session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session disconnected successfully",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/groups": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "List all WhatsApp groups for a session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Groups"
                ],
                "summary": "List WhatsApp groups",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.ListGroupsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new WhatsApp group with specified participants",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Groups"
                ],
                "summary": "Create new WhatsApp group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Group creation request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.CreateGroupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.CreateGroupResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/groups/info": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get detailed information about a WhatsApp group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Groups"
                ],
                "summary": "Get group information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Group JID",
                        "name": "groupJid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.GetGroupInfoResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/groups/name": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Change the name of a WhatsApp group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Groups"
                ],
                "summary": "Set group name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Group name request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.SetGroupNameRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.SetGroupNameResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/groups/participants": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Add, remove, promote or demote group participants",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Groups"
                ],
                "summary": "Update group participants",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Participants update request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.UpdateParticipantsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.UpdateParticipantsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
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
                "description": "Get detailed information about a specific WhatsApp session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get session information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session information retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SessionInfoResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
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
                "description": "Logout from WhatsApp session and disconnect",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Logout session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session logout successful",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/media/clear-cache": {
            "post": {
                "description": "Clear all cached media files for the session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "Clear media cache",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/media/download": {
            "post": {
                "description": "Download media file from WhatsApp message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "Download media from WhatsApp",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/media/info": {
            "get": {
                "description": "Get information about media files",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "Get media information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/media/list": {
            "get": {
                "description": "List all cached media files for the session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "List cached media files",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/media/stats": {
            "get": {
                "description": "Get statistics about media usage for the session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "Get media statistics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/edit": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Edit a previously sent message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Edit message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Edit message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/EditMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/mark-read": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Mark messages as read in WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Mark messages as read",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Mark as read request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/MarkAsReadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/pending-sync": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get messages that are pending synchronization with Chatwoot",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Get pending sync messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit (default: 50, max: 100)",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.MessageDTO"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/poll/{messageId}/results": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get results of a poll message via WhatsApp",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Get poll results",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Message ID",
                        "name": "messageId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/GetPollResultsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/revoke": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Revoke a previously sent message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Revoke message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Revoke message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RevokeMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/audio": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send an audio message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send audio message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Audio message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendAudioMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/button": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a button message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send button message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Button message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendButtonMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/contact": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a contact message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send contact message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Contact message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendContactMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/contact-list": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a contact list message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send contact list message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Contact list message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendContactListMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendContactListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/document": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a document message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send document message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Document message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendDocumentMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/image": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send an image message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send image message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Image message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendImageMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/list": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a list message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send list message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "List message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendListMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/location": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a location message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send location message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Location message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendLocationMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/media": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a media message (image, video, audio, document) via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send media message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Media message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendMediaMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/poll": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a poll message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send poll message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Poll message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendPollMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/presence": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send presence status (typing, recording, etc.) via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send presence status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Presence message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendPresenceMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/profile/business": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a business profile message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send business profile message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Business profile message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendBusinessProfileMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/reaction": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a reaction to a message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send reaction message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Reaction message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendReactionMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/sticker": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a sticker message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send sticker message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Sticker message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendStickerMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/text": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a text message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send text message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Text message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendTextMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/send/video": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send a video message via WhatsApp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Send video message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Video message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendVideoMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/SendMessageResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/messages/{messageId}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete a message from the system",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Delete message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Message ID",
                        "name": "messageId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/pair": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Pair WhatsApp session with phone number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Pair phone number",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Phone pairing request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/PairPhoneRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Phone pairing initiated successfully",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/proxy": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get proxy configuration for a WhatsApp session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get proxy",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Proxy configuration retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/ProxyResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/proxy/set": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Configure proxy settings for a WhatsApp session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Set proxy",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Proxy configuration",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SetProxyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Proxy configured successfully",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
                        "description": "Internal Server Error",
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
                "description": "Get QR code for WhatsApp session pairing. Returns both raw QR code string and base64 image.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get QR code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "QR code generated successfully with base64 image",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/QRCodeResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/qr/generate": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Generate a new QR code for WhatsApp session pairing",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Generate QR code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "QR code generated successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/QRCodeResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Session not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/webhook/find": {
            "get": {
                "description": "Get the current webhook configuration for the session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Get webhook configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/webhook/set": {
            "post": {
                "description": "Configure webhook settings for the session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Set webhook configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{sessionId}/webhook/test": {
            "post": {
                "description": "Test the webhook configuration by sending a test event",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhooks"
                ],
                "summary": "Test webhook configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/SuccessResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ButtonInfo": {
            "type": "object",
            "required": [
                "id",
                "text"
            ],
            "properties": {
                "id": {
                    "type": "string",
                    "example": "btn-1"
                },
                "text": {
                    "type": "string",
                    "example": "Option 1"
                },
                "type": {
                    "type": "string",
                    "example": "reply"
                }
            }
        },
        "ConnectSessionResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Session connection initiated successfully"
                },
                "qrCode": {
                    "type": "string",
                    "example": "2@abc123..."
                },
                "qrCodeImage": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "ContactInfo": {
            "type": "object",
            "required": [
                "name",
                "phone"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511888888888"
                }
            }
        },
        "ContactResult": {
            "type": "object",
            "properties": {
                "contact_name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "error": {
                    "type": "string",
                    "example": ""
                },
                "message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511888888888"
                },
                "status": {
                    "type": "string",
                    "example": "sent"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "ContextInfo": {
            "type": "object",
            "required": [
                "stanzaId"
            ],
            "properties": {
                "participant": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "stanzaId": {
                    "type": "string",
                    "example": "ABCD1234abcd"
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
                    "maxLength": 50,
                    "minLength": 3,
                    "example": "my-session"
                },
                "proxyConfig": {
                    "$ref": "#/definitions/ProxyConfig"
                },
                "qrCode": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "CreateSessionResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2024-01-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "1b2e424c-a2a0-41a4-b992-15b7ec06b9bc"
                },
                "isConnected": {
                    "type": "boolean",
                    "example": false
                },
                "name": {
                    "type": "string",
                    "example": "my-session"
                },
                "proxyConfig": {
                    "$ref": "#/definitions/ProxyConfig"
                },
                "qrCode": {
                    "type": "string",
                    "example": "2@abc123..."
                },
                "qrCodeImage": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
                }
            }
        },
        "DeviceInfoResponse": {
            "type": "object",
            "properties": {
                "appVersion": {
                    "type": "string",
                    "example": "2.21.4.18"
                },
                "deviceModel": {
                    "type": "string",
                    "example": "Samsung Galaxy S21"
                },
                "osVersion": {
                    "type": "string",
                    "example": "11"
                },
                "platform": {
                    "type": "string",
                    "example": "android"
                }
            }
        },
        "EditMessageRequest": {
            "type": "object",
            "required": [
                "message_id",
                "new_text",
                "to"
            ],
            "properties": {
                "message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "new_body": {
                    "type": "string",
                    "example": "Updated message"
                },
                "new_text": {
                    "type": "string",
                    "example": "Updated message"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "VALIDATION_ERROR"
                },
                "details": {},
                "error": {
                    "type": "string",
                    "example": "Invalid request"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "GetPollResultsResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "data": {},
                "message": {
                    "type": "string",
                    "example": "Operation completed successfully"
                },
                "message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "poll_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "poll_name": {
                    "type": "string",
                    "example": "Favorite Color Poll"
                },
                "question": {
                    "type": "string",
                    "example": "What's your favorite color?"
                },
                "request_id": {
                    "type": "string",
                    "example": "req-123"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "timestamp": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "total_votes": {
                    "type": "integer",
                    "example": 15
                },
                "vote_results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/PollVoteInfo"
                    }
                },
                "votes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/PollVoteInfo"
                    }
                }
            }
        },
        "ListRowInfo": {
            "type": "object",
            "required": [
                "id",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Description of option 1"
                },
                "id": {
                    "type": "string",
                    "example": "row-1"
                },
                "title": {
                    "type": "string",
                    "example": "Option 1"
                }
            }
        },
        "ListSectionInfo": {
            "type": "object",
            "required": [
                "rows",
                "title"
            ],
            "properties": {
                "rows": {
                    "type": "array",
                    "maxItems": 10,
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/ListRowInfo"
                    }
                },
                "title": {
                    "type": "string",
                    "example": "Section 1"
                }
            }
        },
        "ListSessionsResponse": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer",
                    "example": 20
                },
                "offset": {
                    "type": "integer",
                    "example": 0
                },
                "sessions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/SessionInfoResponse"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "MarkAsReadRequest": {
            "type": "object",
            "required": [
                "chat_jid",
                "message_ids"
            ],
            "properties": {
                "chat_jid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "message_ids": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "[\"3EB0C767D71D\"]"
                    ]
                }
            }
        },
        "PairPhoneRequest": {
            "type": "object",
            "required": [
                "phoneNumber"
            ],
            "properties": {
                "phoneNumber": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "PollOptionInfo": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "Option 1"
                }
            }
        },
        "PollVoteInfo": {
            "type": "object",
            "properties": {
                "option_name": {
                    "type": "string",
                    "example": "Option 1"
                },
                "vote_count": {
                    "type": "integer",
                    "example": 5
                },
                "voters": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "[\"5511888888888@s.whatsapp.net\"]"
                    ]
                }
            }
        },
        "ProxyConfig": {
            "type": "object",
            "required": [
                "host",
                "port",
                "type"
            ],
            "properties": {
                "host": {
                    "type": "string",
                    "example": "proxy.example.com"
                },
                "password": {
                    "type": "string",
                    "example": "proxypass123"
                },
                "port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1,
                    "example": 8080
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "http",
                        "socks5"
                    ],
                    "example": "http"
                },
                "username": {
                    "type": "string",
                    "example": "proxyuser"
                }
            }
        },
        "ProxyResponse": {
            "type": "object",
            "properties": {
                "proxyConfig": {
                    "$ref": "#/definitions/ProxyConfig"
                }
            }
        },
        "QRCodeResponse": {
            "type": "object",
            "properties": {
                "expiresAt": {
                    "type": "string",
                    "example": "2024-01-01T00:01:00Z"
                },
                "qrCode": {
                    "type": "string",
                    "example": "2@abc123def456..."
                },
                "qrCodeImage": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
                },
                "timeoutSeconds": {
                    "type": "integer",
                    "example": 60
                }
            }
        },
        "RevokeMessageRequest": {
            "type": "object",
            "required": [
                "message_id",
                "to"
            ],
            "properties": {
                "message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendAudioMessageRequest": {
            "type": "object",
            "required": [
                "file",
                "to"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Audio message"
                },
                "file": {
                    "type": "string",
                    "example": "base64_audio_data"
                },
                "filename": {
                    "type": "string",
                    "example": "audio.mp3"
                },
                "mime_type": {
                    "type": "string",
                    "example": "audio/mpeg"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendBusinessProfileMessageRequest": {
            "type": "object",
            "required": [
                "to"
            ],
            "properties": {
                "business_jid": {
                    "type": "string",
                    "example": "5511888888888@s.whatsapp.net"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendButtonMessageRequest": {
            "type": "object",
            "required": [
                "buttons",
                "text",
                "to"
            ],
            "properties": {
                "buttons": {
                    "type": "array",
                    "maxItems": 3,
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/ButtonInfo"
                    }
                },
                "footer": {
                    "type": "string",
                    "example": "Powered by ZPWoot"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "text": {
                    "type": "string",
                    "example": "Choose an option:"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendContactListMessageRequest": {
            "type": "object",
            "required": [
                "contacts",
                "to"
            ],
            "properties": {
                "contacts": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/ContactInfo"
                    }
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendContactListResponse": {
            "type": "object",
            "properties": {
                "contact_count": {
                    "type": "integer",
                    "example": 3
                },
                "contact_results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ContactResult"
                    }
                },
                "data": {},
                "message": {
                    "type": "string",
                    "example": "Operation completed successfully"
                },
                "remote_jid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "request_id": {
                    "type": "string",
                    "example": "req-123"
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ContactResult"
                    }
                },
                "sent_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "session_id": {
                    "type": "string",
                    "example": "session-123"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "timestamp": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                }
            }
        },
        "SendContactMessageRequest": {
            "type": "object",
            "required": [
                "name",
                "phone",
                "to"
            ],
            "properties": {
                "contact_name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "contact_phone": {
                    "type": "string",
                    "example": "+5511888888888"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511888888888"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendDocumentMessageRequest": {
            "type": "object",
            "required": [
                "file",
                "filename",
                "to"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Document"
                },
                "file": {
                    "type": "string",
                    "example": "base64_document_data"
                },
                "filename": {
                    "type": "string",
                    "example": "document.pdf"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendImageMessageRequest": {
            "type": "object",
            "required": [
                "file",
                "to"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Check this image!"
                },
                "file": {
                    "type": "string",
                    "example": "base64_image_data"
                },
                "filename": {
                    "type": "string",
                    "example": "image.jpg"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendListMessageRequest": {
            "type": "object",
            "required": [
                "body",
                "button_text",
                "sections",
                "to"
            ],
            "properties": {
                "body": {
                    "type": "string",
                    "example": "Choose from the list:"
                },
                "button_text": {
                    "type": "string",
                    "example": "View Options"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "sections": {
                    "type": "array",
                    "maxItems": 10,
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/ListSectionInfo"
                    }
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendLocationMessageRequest": {
            "type": "object",
            "required": [
                "latitude",
                "longitude",
                "to"
            ],
            "properties": {
                "address": {
                    "type": "string",
                    "example": "São Paulo, SP, Brazil"
                },
                "latitude": {
                    "type": "number",
                    "example": -23.5505
                },
                "longitude": {
                    "type": "number",
                    "example": -46.6333
                },
                "name": {
                    "type": "string",
                    "example": "São Paulo"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendMediaMessageRequest": {
            "type": "object",
            "required": [
                "media_url",
                "to",
                "type"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Check this out!"
                },
                "filename": {
                    "type": "string",
                    "example": "image.jpg"
                },
                "media_url": {
                    "type": "string",
                    "example": "https://example.com/image.jpg"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "image",
                        "audio",
                        "video",
                        "document"
                    ],
                    "example": "image"
                }
            }
        },
        "SendMessageResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "delivered_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:05Z"
                },
                "message": {
                    "type": "string",
                    "example": "Operation completed successfully"
                },
                "message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "read_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:10Z"
                },
                "request_id": {
                    "type": "string",
                    "example": "req-123"
                },
                "status": {
                    "type": "string",
                    "example": "sent"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "timestamp": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendPollMessageRequest": {
            "type": "object",
            "required": [
                "name",
                "options",
                "question",
                "to"
            ],
            "properties": {
                "allow_multiple_vote": {
                    "type": "boolean",
                    "example": false
                },
                "name": {
                    "type": "string",
                    "example": "Favorite Color Poll"
                },
                "options": {
                    "type": "array",
                    "maxItems": 12,
                    "minItems": 2,
                    "items": {
                        "$ref": "#/definitions/PollOptionInfo"
                    }
                },
                "question": {
                    "type": "string",
                    "example": "What's your favorite color?"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "selectable_count": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 1
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendPresenceMessageRequest": {
            "type": "object",
            "required": [
                "presence",
                "to"
            ],
            "properties": {
                "presence": {
                    "type": "string",
                    "enum": [
                        "typing",
                        "recording",
                        "online",
                        "offline",
                        "paused"
                    ],
                    "example": "typing"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendReactionMessageRequest": {
            "type": "object",
            "required": [
                "message_id",
                "reaction",
                "to"
            ],
            "properties": {
                "message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "reaction": {
                    "type": "string",
                    "example": "👍"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendStickerMessageRequest": {
            "type": "object",
            "required": [
                "file",
                "to"
            ],
            "properties": {
                "file": {
                    "type": "string",
                    "example": "base64_sticker_data"
                },
                "mime_type": {
                    "type": "string",
                    "example": "image/webp"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendTextMessageRequest": {
            "type": "object",
            "required": [
                "body",
                "remoteJid"
            ],
            "properties": {
                "body": {
                    "type": "string",
                    "example": "Hello, World!"
                },
                "contextInfo": {
                    "$ref": "#/definitions/ContextInfo"
                },
                "remoteJid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SendVideoMessageRequest": {
            "type": "object",
            "required": [
                "file",
                "to"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Check this video!"
                },
                "file": {
                    "type": "string",
                    "example": "base64_video_data"
                },
                "filename": {
                    "type": "string",
                    "example": "video.mp4"
                },
                "reply_to": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "to": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "SessionInfoResponse": {
            "type": "object",
            "properties": {
                "deviceInfo": {
                    "$ref": "#/definitions/DeviceInfoResponse"
                },
                "session": {
                    "$ref": "#/definitions/SessionResponse"
                }
            }
        },
        "SessionResponse": {
            "type": "object",
            "properties": {
                "connectedAt": {
                    "type": "string",
                    "example": "2024-01-01T00:00:30Z"
                },
                "connectionError": {
                    "type": "string",
                    "example": "Connection timeout"
                },
                "createdAt": {
                    "type": "string",
                    "example": "2024-01-01T00:00:00Z"
                },
                "deviceJid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "id": {
                    "type": "string",
                    "example": "session-123"
                },
                "isConnected": {
                    "type": "boolean",
                    "example": false
                },
                "name": {
                    "type": "string",
                    "example": "my-whatsapp-session"
                },
                "proxyConfig": {
                    "$ref": "#/definitions/ProxyConfig"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2024-01-01T00:00:00Z"
                }
            }
        },
        "SessionStatsResponse": {
            "type": "object",
            "properties": {
                "connected": {
                    "type": "integer",
                    "example": 3
                },
                "offline": {
                    "type": "integer",
                    "example": 7
                },
                "total": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "SetProxyRequest": {
            "type": "object",
            "required": [
                "proxyConfig"
            ],
            "properties": {
                "proxyConfig": {
                    "$ref": "#/definitions/ProxyConfig"
                }
            }
        },
        "SuccessResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string",
                    "example": "Operation completed successfully"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "internal_adapters_server_handler.BusinessProfileResponse": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "category": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "isBusiness": {
                    "type": "boolean"
                },
                "jid": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "internal_adapters_server_handler.CheckWhatsAppRequest": {
            "type": "object",
            "required": [
                "phoneNumbers"
            ],
            "properties": {
                "phoneNumbers": {
                    "type": "array",
                    "maxItems": 50,
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "internal_adapters_server_handler.CheckWhatsAppResponse": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_adapters_server_handler.CheckWhatsAppResult"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "internal_adapters_server_handler.CheckWhatsAppResult": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "isOnWhatsApp": {
                    "type": "boolean"
                },
                "jid": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                }
            }
        },
        "internal_adapters_server_handler.ContactInfo": {
            "type": "object",
            "properties": {
                "isBusiness": {
                    "type": "boolean"
                },
                "isMyContact": {
                    "type": "boolean"
                },
                "isOnWhatsApp": {
                    "type": "boolean"
                },
                "jid": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "pushName": {
                    "type": "string"
                },
                "shortName": {
                    "type": "string"
                }
            }
        },
        "internal_adapters_server_handler.GetProfilePictureResponse": {
            "type": "object",
            "properties": {
                "hasPicture": {
                    "type": "boolean"
                },
                "jid": {
                    "type": "string"
                },
                "pictureId": {
                    "type": "string"
                },
                "pictureUrl": {
                    "type": "string"
                }
            }
        },
        "internal_adapters_server_handler.GetUserInfoRequest": {
            "type": "object",
            "required": [
                "jids"
            ],
            "properties": {
                "jids": {
                    "type": "array",
                    "maxItems": 20,
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "internal_adapters_server_handler.GetUserInfoResponse": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_adapters_server_handler.UserInfoResult"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "internal_adapters_server_handler.ListContactsResponse": {
            "type": "object",
            "properties": {
                "contacts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_adapters_server_handler.ContactInfo"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "internal_adapters_server_handler.SyncContactsResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "syncedContacts": {
                    "type": "integer"
                },
                "totalContacts": {
                    "type": "integer"
                }
            }
        },
        "internal_adapters_server_handler.UserInfoResult": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "isBusiness": {
                    "type": "boolean"
                },
                "isOnWhatsApp": {
                    "type": "boolean"
                },
                "jid": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "pictureId": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.CreateGroupRequest": {
            "type": "object",
            "required": [
                "name",
                "participants"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "maxLength": 512
                },
                "name": {
                    "type": "string",
                    "maxLength": 25,
                    "minLength": 1
                },
                "participants": {
                    "type": "array",
                    "maxItems": 256,
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.CreateGroupResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "group_jid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.GetGroupInfoResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "group_jid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.ParticipantInfo"
                    }
                },
                "settings": {
                    "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.GroupSettings"
                },
                "success": {
                    "type": "boolean"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.GroupInfo": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "group_jid": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "participants": {
                    "type": "integer"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.GroupSettings": {
            "type": "object",
            "properties": {
                "announce": {
                    "type": "boolean"
                },
                "join_approval_mode": {
                    "type": "string"
                },
                "locked": {
                    "type": "boolean"
                },
                "member_add_mode": {
                    "type": "string"
                },
                "restrict": {
                    "type": "boolean"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.ListGroupsResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "groups": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/zpwoot_internal_adapters_server_contracts.GroupInfo"
                    }
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.MessageDTO": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string",
                    "example": "Hello World"
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "cw_conversation_id": {
                    "type": "integer",
                    "example": 456
                },
                "cw_message_id": {
                    "type": "integer",
                    "example": 123
                },
                "id": {
                    "type": "string",
                    "example": "1b2e424c-a2a0-41a4-b992-15b7ec06b9bc"
                },
                "media_type": {
                    "type": "string",
                    "example": "image"
                },
                "media_url": {
                    "type": "string",
                    "example": "https://example.com/image.jpg"
                },
                "session_id": {
                    "type": "string",
                    "example": "session-123"
                },
                "sync_status": {
                    "type": "string",
                    "example": "synced"
                },
                "synced_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:05Z"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "zp_chat": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "zp_from_me": {
                    "type": "boolean",
                    "example": false
                },
                "zp_message_id": {
                    "type": "string",
                    "example": "3EB0C767D71D"
                },
                "zp_sender": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "zp_timestamp": {
                    "type": "string",
                    "example": "2024-01-01T12:00:00Z"
                },
                "zp_type": {
                    "type": "string",
                    "example": "text"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.ParticipantInfo": {
            "type": "object",
            "properties": {
                "jid": {
                    "type": "string"
                },
                "joined_at": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.SetGroupNameRequest": {
            "type": "object",
            "required": [
                "group_jid",
                "name"
            ],
            "properties": {
                "group_jid": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 25,
                    "minLength": 1
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.SetGroupNameResponse": {
            "type": "object",
            "properties": {
                "group_jid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.UpdateParticipantsRequest": {
            "type": "object",
            "required": [
                "action",
                "group_jid",
                "participants"
            ],
            "properties": {
                "action": {
                    "type": "string",
                    "enum": [
                        "add",
                        "remove",
                        "promote",
                        "demote"
                    ]
                },
                "group_jid": {
                    "type": "string"
                },
                "participants": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "zpwoot_internal_adapters_server_contracts.UpdateParticipantsResponse": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "group_jid": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "API Key authentication. Use: YOUR_API_KEY",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "2.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "ZPWoot WhatsApp API",
	Description:      "A comprehensive WhatsApp Business API built with Go. Provides endpoints for session management, messaging, contacts, groups, media handling, and integrations with Chatwoot.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
