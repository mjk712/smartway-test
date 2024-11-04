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
        "/document/{documentId}": {
            "put": {
                "description": "Обновляет данные документа по id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Документы"
                ],
                "summary": "Обновление информации о документе",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID документа",
                        "name": "documentId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные для обновления документа",
                        "name": "document",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.DocumentUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Документ успешно обновлён",
                        "schema": {
                            "$ref": "#/definitions/models.Document"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или обновления документа"
                    }
                }
            },
            "delete": {
                "description": "Удаляет документ по id.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Документы"
                ],
                "summary": "Удаление документа",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Document ID",
                        "name": "documentId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Документ успешно удалён"
                    },
                    "404": {
                        "description": "Ошибка в запросе или при удалении документа"
                    }
                }
            }
        },
        "/documents/{passengerId}": {
            "get": {
                "description": "Возвращает список документов. Требует id пассажира.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Документы"
                ],
                "summary": "Получение списка документов",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Passenger ID",
                        "name": "passengerId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список продуктов успешно получен",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Document"
                            }
                        }
                    },
                    "404": {
                        "description": "Ошибка в запросе или при получении списка документов"
                    }
                }
            }
        },
        "/passenger/{passengerId}": {
            "put": {
                "description": "Обновляет данные пассажира по заданному id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пассажиры"
                ],
                "summary": "Обновление информации о пассажире",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID пассажира",
                        "name": "passengerId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные для обновления информации о пассажире",
                        "name": "passenger",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.UpdatePassengerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пассажир успешно обновлён",
                        "schema": {
                            "$ref": "#/definitions/models.Passenger"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или обновления пассажира"
                    }
                }
            },
            "delete": {
                "description": "Удаляет пассажира по id.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Пассажиры"
                ],
                "summary": "Удаление пассажира",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Passenger ID",
                        "name": "passengerId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пассажир успешно удалён"
                    },
                    "404": {
                        "description": "Ошибка в запросе или при удалении пассажира"
                    }
                }
            }
        },
        "/passengers/{ticketNumber}": {
            "get": {
                "description": "Возвращает список пассажиров по номеру билета",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пассажиры"
                ],
                "summary": "Получение списка пассажиров по номеру билета",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Номер билета",
                        "name": "ticketNumber",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список пассажиров успешно получен",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Passenger"
                            }
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или получения списка пассажиров"
                    }
                }
            }
        },
        "/reports/passenger/{passengerId}": {
            "get": {
                "description": "Возвращает отчет о пассажире по заданному ` + "`" + `passengerId` + "`" + ` и диапазону дат",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пассажиры"
                ],
                "summary": "Получение отчета о пассажире",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пассажира",
                        "name": "passengerId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Дата начала в формате YYYY-MM-DD",
                        "name": "start_date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Дата окончания в формате YYYY-MM-DD",
                        "name": "end_date",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Отчёт успешно получен",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response.FlightReport"
                            }
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или получения отчёта"
                    }
                }
            }
        },
        "/ticket/{ticketId}": {
            "put": {
                "description": "Обновляет данные билета по заданному ` + "`" + `ticketId` + "`" + `",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Билеты"
                ],
                "summary": "Обновление информации о билете",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID билета",
                        "name": "ticketId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные для обновления информации о билете",
                        "name": "ticket",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.TicketUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Билет успешно обновлён",
                        "schema": {
                            "$ref": "#/definitions/models.Ticket"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или обновления билета"
                    }
                }
            },
            "delete": {
                "description": "Удаляет билет по id и связь билета с пассажиром.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Билеты"
                ],
                "summary": "Удаление билета",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Ticket ID",
                        "name": "ticketId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Билет успешно удалён"
                    },
                    "404": {
                        "description": "Ошибка в запросе или при удалении билета"
                    }
                }
            }
        },
        "/ticket/{ticketNumber}": {
            "get": {
                "description": "Возвращает полные данные о билете по заданному номеру",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Билеты"
                ],
                "summary": "Получение полной информации о билете",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Номер билета",
                        "name": "ticketNumber",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация о билете успешно получена",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response.FullTicketInfo"
                            }
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или получения полной информации о билете"
                    }
                }
            }
        },
        "/tickets": {
            "get": {
                "description": "Возвращает список всех доступных билетов",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Билеты"
                ],
                "summary": "Получение списка билетов",
                "responses": {
                    "200": {
                        "description": "Список билетов успешно получен",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Ticket"
                            }
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса или получения списка билетов"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Document": {
            "type": "object",
            "properties": {
                "documentId": {
                    "type": "integer"
                },
                "documentNumber": {
                    "type": "string"
                },
                "documentType": {
                    "type": "string"
                },
                "passengerId": {
                    "type": "integer"
                }
            }
        },
        "models.Passenger": {
            "type": "object",
            "properties": {
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "middleName": {
                    "type": "string"
                },
                "passengerId": {
                    "type": "integer"
                }
            }
        },
        "models.Ticket": {
            "type": "object",
            "properties": {
                "arrivalDate": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "departureDate": {
                    "type": "string"
                },
                "departurePoint": {
                    "type": "string"
                },
                "destinationPoint": {
                    "type": "string"
                },
                "orderNumber": {
                    "type": "string"
                },
                "serviceProvider": {
                    "type": "string"
                },
                "ticketId": {
                    "type": "integer"
                }
            }
        },
        "requests.DocumentUpdateRequest": {
            "type": "object",
            "properties": {
                "documentNumber": {
                    "type": "string"
                },
                "documentType": {
                    "type": "string"
                },
                "passengerId": {
                    "type": "integer"
                }
            }
        },
        "requests.TicketUpdateRequest": {
            "type": "object",
            "properties": {
                "arrivalDate": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "departureDate": {
                    "type": "string"
                },
                "departurePoint": {
                    "type": "string"
                },
                "destinationPoint": {
                    "type": "string"
                },
                "serviceProvider": {
                    "type": "string"
                }
            }
        },
        "requests.UpdatePassengerRequest": {
            "type": "object",
            "properties": {
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "middleName": {
                    "type": "string"
                }
            }
        },
        "response.FlightReport": {
            "type": "object",
            "properties": {
                "arrivalDate": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "departureDate": {
                    "type": "string"
                },
                "departurePoint": {
                    "type": "string"
                },
                "destinationPoint": {
                    "type": "string"
                },
                "flightStatus": {
                    "type": "string"
                },
                "orderNumber": {
                    "type": "string"
                },
                "serviceProvider": {
                    "type": "string"
                },
                "ticketId": {
                    "type": "integer"
                }
            }
        },
        "response.FullTicketInfo": {
            "type": "object",
            "properties": {
                "passengers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.PassengerWithDocs"
                    }
                },
                "ticket": {
                    "$ref": "#/definitions/models.Ticket"
                }
            }
        },
        "response.PassengerWithDocs": {
            "type": "object",
            "properties": {
                "documents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Document"
                    }
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "middleName": {
                    "type": "string"
                },
                "passengerId": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}