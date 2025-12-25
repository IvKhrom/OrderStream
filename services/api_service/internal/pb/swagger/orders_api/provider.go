package orders_api

import _ "embed"

//go:embed orders.swagger.json
var swaggerJSON []byte

func SwaggerJSON() []byte {
	return swaggerJSON
}
