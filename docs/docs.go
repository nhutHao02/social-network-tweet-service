// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/all": {
            "get": {
                "description": "Get All Tweets",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tweet"
                ],
                "summary": "GetAllTweets",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003cyour_token\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common.PagingSuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.GetTweetsRes"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "default": {
                        "description": "failure",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/tweet": {
            "get": {
                "description": "Get tweet by user id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tweet"
                ],
                "summary": "GetTweetByUserID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003cyour_token\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successful",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common.PagingSuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.GetTweetByUserRes"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "default": {
                        "description": "failure",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "description": "Post new Tweet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tweet"
                ],
                "summary": "PostTweet",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003cyour_token\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Post Tweet Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostTweetReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "boolean"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "default": {
                        "description": "failure",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common.ErrorMessage": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "common.PagingSuccessResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "success": {
                    "type": "boolean"
                },
                "totalPage": {
                    "type": "integer"
                }
            }
        },
        "common.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "$ref": "#/definitions/common.ErrorMessage"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "model.GetTweetByUserRes": {
            "type": "object",
            "properties": {
                "action": {
                    "$ref": "#/definitions/model.UserAction"
                },
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "statistics": {
                    "$ref": "#/definitions/model.Statistics"
                },
                "userID": {
                    "type": "integer"
                },
                "userInfo": {
                    "$ref": "#/definitions/model.UserInfo"
                }
            }
        },
        "model.GetTweetsRes": {
            "type": "object",
            "properties": {
                "action": {
                    "$ref": "#/definitions/model.UserAction"
                },
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "statistics": {
                    "$ref": "#/definitions/model.Statistics"
                },
                "userID": {
                    "type": "integer"
                },
                "userInfo": {
                    "$ref": "#/definitions/model.UserInfo"
                }
            }
        },
        "model.PostTweetReq": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "urlImg": {
                    "type": "string"
                },
                "urlVideo": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "model.Statistics": {
            "type": "object",
            "properties": {
                "totalBookmark": {
                    "type": "integer"
                },
                "totalComment": {
                    "type": "integer"
                },
                "totalLove": {
                    "type": "integer"
                },
                "totalRepost": {
                    "type": "integer"
                }
            }
        },
        "model.UserAction": {
            "type": "object",
            "properties": {
                "bookmark": {
                    "type": "boolean"
                },
                "love": {
                    "type": "boolean"
                },
                "repost": {
                    "type": "boolean"
                }
            }
        },
        "model.UserInfo": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "urlAvt": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Social Network Service",
	Description:      "This is tweet service of the social network implament using Go",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
