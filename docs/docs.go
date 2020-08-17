// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/category": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns list of categories",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "category"
                ],
                "summary": "Returns list of categories",
                "operationId": "get-categories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Category"
                        }
                    }
                }
            }
        },
        "/donation": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns list of user's donations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "donation"
                ],
                "summary": "Returns list of user's donations",
                "operationId": "get-user-donations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Donation"
                        }
                    }
                }
            }
        },
        "/donation/project": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns list of project donations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "donation"
                ],
                "summary": "Returns list of project donations",
                "operationId": "get-project-donations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.ProjectDonation"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "get token for user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Returns access token",
                "operationId": "get-token",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenResponse"
                        }
                    }
                }
            }
        },
        "/project/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns project by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "Show a single project",
                "operationId": "get-project-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.ProjectDetailView"
                        }
                    }
                }
            }
        },
        "/project_type": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns list of project types",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project type"
                ],
                "summary": "return list of project types",
                "operationId": "get-project-types",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ProjectType"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns user by ID from token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Show a current user",
                "operationId": "get-user-by-token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.extendedUser"
                        }
                    }
                }
            }
        },
        "/user/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns user by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Show a specific user",
                "operationId": "get-user-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.extendedUser"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.ProjectDetailView": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "object",
                    "$ref": "#/definitions/models.Category"
                },
                "description": {
                    "type": "string"
                },
                "event_date": {
                    "type": "string"
                },
                "goal_amount": {
                    "type": "integer"
                },
                "goal_people": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "image_link": {
                    "type": "string"
                },
                "instructions": {
                    "type": "string"
                },
                "owner": {
                    "type": "object",
                    "$ref": "#/definitions/models.User"
                },
                "percent": {
                    "type": "integer"
                },
                "project_type": {
                    "type": "object",
                    "$ref": "#/definitions/models.ProjectType"
                },
                "release_date": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "subtitle": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "handlers.ProjectDonation": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "locked": {
                    "type": "boolean"
                },
                "paid": {
                    "type": "boolean"
                },
                "user": {
                    "type": "object",
                    "$ref": "#/definitions/models.User"
                }
            }
        },
        "handlers.TokenRequest": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "handlers.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "handlers.extendedUser": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_name": {
                    "type": "string"
                },
                "participation": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Participation"
                    }
                },
                "project_count": {
                    "type": "integer"
                },
                "success_rate": {
                    "type": "number"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Category": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.Donation": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "locked": {
                    "type": "boolean"
                },
                "paid": {
                    "type": "boolean"
                },
                "payment": {
                    "type": "integer"
                },
                "project": {
                    "type": "integer"
                }
            }
        },
        "models.Participation": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "models.ProjectType": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "end_by_goal_gain": {
                    "type": "boolean"
                },
                "goal_by_amount": {
                    "type": "boolean"
                },
                "goal_by_people": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "options": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_name": {
                    "type": "string"
                },
                "project_count": {
                    "type": "integer"
                },
                "success_rate": {
                    "type": "number"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Launchpad API",
	Description: "This is a launchpad backend.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}