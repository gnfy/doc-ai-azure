package main

import (
	"doc-ai-azure/pkg/api"
	"doc-ai-azure/pkg/common"
	"fmt"
)

func main() {
	r := api.SetupRouter()
	if err := r.Run(":" + common.GlobalObject.Server.Port); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
