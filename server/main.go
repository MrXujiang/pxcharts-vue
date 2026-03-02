package main

import (
	"mvtable/cmd"
)

// @title MVTable
// @version 1.0
// @description MVTable 后端服务 API
// @termsOfService http://swagger.io/terms/
// @contact.name MVTable
// @contact.email mvtable@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 输入Bearer token进行认证，格式: Bearer {token}。例如: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

func main() {
	cmd.Execute()
}
