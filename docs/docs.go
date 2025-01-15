// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/cancel": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "注销用户",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        },
        "/getCaptcha": {
            "get": {
                "description": "获取验证码生成所需的ID和图片URL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.CaptchaResponseData"
                        }
                    }
                }
            }
        },
        "/passwordLogin": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "账号密码登录",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.PasswordLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.LoginResponseData"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "目前只支持邮箱登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        },
        "/user/getCollege": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取学院信息",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.GetCollegeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.CollegeResponseData"
                        }
                    }
                }
            }
        },
        "/user/getCollegeList": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取学院列表",
                "responses": {
                    "200": {
                        "description": "返回学院信息列表",
                        "schema": {
                            "$ref": "#/definitions/v1.GetCollegeListDataResponse"
                        }
                    }
                }
            }
        },
        "/user/getUserInfo": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取用户信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.GetUserInfoResponseData"
                        }
                    }
                }
            }
        },
        "/user/logout": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "退出用户",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        },
        "/user/updateProfile": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.UpdateProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        },
        "/user/userAuth": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户认证",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.UserAuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.CaptchaResponseData": {
            "type": "object",
            "properties": {
                "CaptchaBase64": {
                    "type": "string"
                },
                "captchaId": {
                    "type": "string"
                }
            }
        },
        "v1.CollegeResponseData": {
            "type": "object",
            "properties": {
                "collegeId": {
                    "type": "integer"
                },
                "collegeName": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                }
            }
        },
        "v1.CreateArticleRequest": {
            "type": "object",
            "required": [
                "authorId",
                "content",
                "title",
                "visibleRange"
            ],
            "properties": {
                "authorId": {
                    "description": "作者ID",
                    "type": "string"
                },
                "categoryId": {
                    "description": "文章分类ID",
                    "type": "integer"
                },
                "commentDisabled": {
                    "description": "是否禁用评论",
                    "type": "boolean"
                },
                "content": {
                    "description": "文章内容",
                    "type": "string"
                },
                "contentShort": {
                    "description": "文章摘要",
                    "type": "string"
                },
                "importance": {
                    "description": "文章重要性",
                    "type": "integer"
                },
                "sourceUri": {
                    "description": "文章外链",
                    "type": "string"
                },
                "title": {
                    "description": "文章标题",
                    "type": "string"
                },
                "uploadedFiles": {
                    "description": "上传的文件列表",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.FileUpload"
                    }
                },
                "visibleRange": {
                    "description": "可见范围",
                    "type": "string"
                }
            }
        },
        "v1.CreateArticleResponseData": {
            "type": "object",
            "properties": {
                "articleId": {
                    "description": "文章ID",
                    "type": "integer"
                }
            }
        },
        "v1.FileUpload": {
            "type": "object",
            "required": [
                "fileName",
                "fileUrl"
            ],
            "properties": {
                "fileName": {
                    "description": "文件名",
                    "type": "string"
                },
                "fileUrl": {
                    "description": "文件URL",
                    "type": "string"
                }
            }
        },
        "v1.GetCollegeListDataResponse": {
            "type": "object",
            "properties": {
                "collegeList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.CollegeResponseData"
                    }
                }
            }
        },
        "v1.GetCollegeRequest": {
            "type": "object",
            "properties": {
                "collegeId": {
                    "type": "integer"
                }
            }
        },
        "v1.GetUserInfoResponseData": {
            "type": "object",
            "properties": {
                "collegeId": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string",
                    "example": "alan"
                },
                "phone": {
                    "type": "string",
                    "example": "10012239028"
                },
                "roleType": {
                    "type": "integer",
                    "example": 0
                },
                "studentId": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "v1.LoginResponseData": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                }
            }
        },
        "v1.PasswordLoginRequest": {
            "type": "object",
            "required": [
                "captchaAnswer",
                "captchaId",
                "password",
                "phone"
            ],
            "properties": {
                "captchaAnswer": {
                    "description": "验证码字段",
                    "type": "string"
                },
                "captchaId": {
                    "description": "验证码ID字段",
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "phone": {
                    "type": "string",
                    "example": "10012239028"
                }
            }
        },
        "v1.RegisterRequest": {
            "type": "object",
            "required": [
                "captchaAnswer",
                "captchaId",
                "password",
                "phone"
            ],
            "properties": {
                "captchaAnswer": {
                    "description": "验证码字段",
                    "type": "string"
                },
                "captchaId": {
                    "description": "验证码ID字段",
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "phone": {
                    "type": "string",
                    "example": "10012239028"
                }
            }
        },
        "v1.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.UpdateProfileRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "nickname": {
                    "type": "string",
                    "example": "alan"
                }
            }
        },
        "v1.UserAuthRequest": {
            "type": "object",
            "properties": {
                "collegeId": {
                    "type": "integer"
                },
                "remarks": {
                    "type": "string"
                },
                "studentId": {
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

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Nunu Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
