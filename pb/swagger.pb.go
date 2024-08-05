package pb

const (
	Swagger = `{
  "swagger": "2.0",
  "info": {
    "title": "inference_engine.proto",
    "version": "1.0"
  },
  "tags": [
    {
      "name": "GRPCInferenceService"
    }
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/comate/v1/grpc/dataserver/grpcinferenceservice/modelfetchrequest": {
      "post": {
        "summary": "拉取一个请求，给inference server调用",
        "operationId": "GRPCInferenceService_ModelFetchRequest3",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestResult"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestParams"
            }
          }
        ],
        "tags": [
          "GRPCInferenceService"
        ]
      }
    },
    "/comate/v1/grpc/dataserver/grpcinferenceservice/modelsendresponse": {
      "post": {
        "summary": "拉取一个请求，给inference server调用",
        "operationId": "GRPCInferenceService_ModelFetchRequest4",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestResult"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestParams"
            }
          }
        ],
        "tags": [
          "GRPCInferenceService"
        ]
      }
    },
    "/comate/v1/grpc/dataserver/grpcinferenceservice/modelstreaminfer": {
      "post": {
        "summary": "模型推理请求入口\n输入一个请求，流式返回多个response",
        "operationId": "GRPCInferenceService_ModelStreamInfer2",
        "responses": {
          "200": {
            "description": "A successful response.(streaming responses)",
            "schema": {
              "type": "object",
              "properties": {
                "result": {
                  "$ref": "#/definitions/language_inferenceModelInferResponse"
                },
                "error": {
                  "$ref": "#/definitions/rpcStatus"
                }
              },
              "title": "Stream result of language_inferenceModelInferResponse"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/language_inferenceModelInferRequest"
            }
          }
        ],
        "tags": [
          "GRPCInferenceService"
        ]
      }
    },
    "/comate/v1/grpc/ic/grpcinferenceservice/modelfetchrequest": {
      "post": {
        "summary": "拉取一个请求，给inference server调用",
        "operationId": "GRPCInferenceService_ModelFetchRequest",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestResult"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestParams"
            }
          }
        ],
        "tags": [
          "GRPCInferenceService"
        ]
      }
    },
    "/comate/v1/grpc/ic/grpcinferenceservice/modelsendresponse": {
      "post": {
        "summary": "拉取一个请求，给inference server调用",
        "operationId": "GRPCInferenceService_ModelFetchRequest2",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestResult"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/language_inferenceModelFetchRequestParams"
            }
          }
        ],
        "tags": [
          "GRPCInferenceService"
        ]
      }
    },
    "/comate/v1/grpc/ic/grpcinferenceservice/modelstreaminfer": {
      "post": {
        "summary": "模型推理请求入口\n输入一个请求，流式返回多个response",
        "operationId": "GRPCInferenceService_ModelStreamInfer",
        "responses": {
          "200": {
            "description": "A successful response.(streaming responses)",
            "schema": {
              "type": "object",
              "properties": {
                "result": {
                  "$ref": "#/definitions/language_inferenceModelInferResponse"
                },
                "error": {
                  "$ref": "#/definitions/rpcStatus"
                }
              },
              "title": "Stream result of language_inferenceModelInferResponse"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/language_inferenceModelInferRequest"
            }
          }
        ],
        "tags": [
          "GRPCInferenceService"
        ]
      }
    }
  },
  "definitions": {
    "language_inferenceContentType": {
      "type": "string",
      "enum": [
        "WENXIN",
        "TRITON"
      ],
      "default": "WENXIN"
    },
    "language_inferenceModelFetchRequestParams": {
      "type": "object",
      "properties": {
        "modelId": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "title": "模型全局唯一id"
        },
        "maxRequestNum": {
          "type": "integer",
          "format": "int32",
          "title": "一次返回的最大请求数"
        }
      },
      "title": "拉取请求时，需要给出模型参数"
    },
    "language_inferenceModelFetchRequestResult": {
      "type": "object",
      "properties": {
        "requests": {
          "type": "array",
          "items": {
            "type": "object",
            "$ref": "#/definitions/language_inferenceModelInferRequest"
          },
          "title": "获取到的请求数组"
        }
      }
    },
    "language_inferenceModelInferRequest": {
      "type": "object",
      "properties": {
        "modelId": {
          "type": "string",
          "title": "模型唯一id，"
        },
        "requestId": {
          "type": "string",
          "title": "请求唯一id，"
        },
        "traceId": {
          "type": "string",
          "title": "可用于跟踪同一请求，多次推理的应答，可选"
        },
        "input": {
          "type": "string",
          "title": "语言模型输入，各模型不同，"
        },
        "tenantId": {
          "type": "string",
          "title": "租户信息，可选"
        },
        "inputType": {
          "$ref": "#/definitions/language_inferenceContentType",
          "title": "输入类型，取值 triton、wenxin，默认值是 wenxin"
        }
      }
    },
    "language_inferenceModelInferResponse": {
      "type": "object",
      "properties": {
        "requestId": {
          "type": "string",
          "title": "请求唯一id"
        },
        "sentenceId": {
          "type": "integer",
          "format": "int32",
          "title": "返回的句子id，表示第几句，用于去重和排序"
        },
        "output": {
          "type": "string",
          "title": "语言模型输出"
        },
        "modelId": {
          "type": "string",
          "title": "模型唯一id"
        },
        "traceId": {
          "type": "string",
          "title": "可用于跟踪同一请求，多次推理的应答"
        },
        "tenantId": {
          "type": "string",
          "title": "租户信息，可选"
        },
        "outputType": {
          "$ref": "#/definitions/language_inferenceContentType",
          "title": "输出类型，取值 triton、wenxin，默认值是 wenxin"
        }
      }
    },
    "language_inferenceModelSendResponseResult": {
      "type": "object",
      "title": "无需关心SendResponse的返回值"
    },
    "protobufAny": {
      "type": "object",
      "properties": {
        "@type": {
          "type": "string"
        }
      },
      "additionalProperties": {}
    },
    "rpcStatus": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        },
        "details": {
          "type": "array",
          "items": {
            "type": "object",
            "$ref": "#/definitions/protobufAny"
          }
        }
      }
    }
  }
}
`
)
