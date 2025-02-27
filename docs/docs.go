// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/users": {
            "get": {
                "description": "List user with pages",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "List users",
                "operationId": "list-users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "pageNo",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Order asc/desc",
                        "name": "order",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Order by create_time",
                        "name": "orderBy",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Status 1/2",
                        "name": "status",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ListUserResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create user, user_id is generated random",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create user",
                "operationId": "create-user",
                "parameters": [
                    {
                        "description": "Set up user",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/httpv1.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpv1.CreateUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/:username": {
            "get": {
                "description": "Get user by username",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user",
                "operationId": "get-user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpv1.GetUserResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update user by username",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user",
                "operationId": "update-user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpv1.UpdateUserResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete user by username",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Delete user",
                "operationId": "delete-user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpv1.DeleteUserResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpv1.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.User": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "创建时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "description": {
                    "description": "备注",
                    "type": "string",
                    "example": "twfbmbsr"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string",
                    "example": "xxx@example.com"
                },
                "status": {
                    "description": "用户状态，1代表启用，2代表禁用",
                    "type": "integer",
                    "example": 1
                },
                "updatedAt": {
                    "description": "更新时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "username": {
                    "description": "用户的名称",
                    "type": "string",
                    "example": "twfbmbsr"
                }
            }
        },
        "httpv1.CreateUserRequest": {
            "type": "object",
            "required": [
                "confirmPassword",
                "email",
                "password",
                "username"
            ],
            "properties": {
                "confirmPassword": {
                    "type": "string"
                },
                "description": {
                    "type": "string",
                    "maxLength": 140,
                    "minLength": 0
                },
                "email": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 0
                },
                "password": {
                    "description": "密码的验证比较复杂，有单独的方法进行验证",
                    "type": "string"
                },
                "username": {
                    "description": "username 应该是长度1-64位的字母、数字或\"_\"",
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 1
                }
            }
        },
        "httpv1.CreateUserResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "创建时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "description": {
                    "description": "备注",
                    "type": "string",
                    "example": "twfbmbsr"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string",
                    "example": "xxx@example.com"
                },
                "status": {
                    "description": "用户状态，1代表启用，2代表禁用",
                    "type": "integer",
                    "example": 1
                },
                "updatedAt": {
                    "description": "更新时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "username": {
                    "description": "用户的名称",
                    "type": "string",
                    "example": "twfbmbsr"
                }
            }
        },
        "httpv1.DeleteUserResponse": {
            "type": "object"
        },
        "httpv1.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "Conflict"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                },
                "requestId": {
                    "type": "string",
                    "example": "b5953bf0-9f15-4c42-afb4-1c125b40d7ce"
                }
            }
        },
        "httpv1.GetUserResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "创建时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "description": {
                    "description": "备注",
                    "type": "string",
                    "example": "twfbmbsr"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string",
                    "example": "xxx@example.com"
                },
                "status": {
                    "description": "用户状态，1代表启用，2代表禁用",
                    "type": "integer",
                    "example": 1
                },
                "updatedAt": {
                    "description": "更新时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "username": {
                    "description": "用户的名称",
                    "type": "string",
                    "example": "twfbmbsr"
                }
            }
        },
        "httpv1.ListUserResponse": {
            "type": "object",
            "properties": {
                "pageNo": {
                    "type": "integer"
                },
                "pageSize": {
                    "type": "integer"
                },
                "result": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.User"
                    }
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "httpv1.UpdateUserResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "创建时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "description": {
                    "description": "备注",
                    "type": "string",
                    "example": "twfbmbsr"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string",
                    "example": "xxx@example.com"
                },
                "status": {
                    "description": "用户状态，1代表启用，2代表禁用",
                    "type": "integer",
                    "example": 1
                },
                "updatedAt": {
                    "description": "更新时间",
                    "type": "string",
                    "example": "2020-01-01T00:00:00Z"
                },
                "username": {
                    "description": "用户的名称",
                    "type": "string",
                    "example": "twfbmbsr"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/.",
	Schemes:          []string{},
	Title:            "GO WEBAPP TEMPLATE API",
	Description:      "GO WEBAPP TEMPLATE API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
