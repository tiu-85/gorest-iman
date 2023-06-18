package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
)

func GetTask(ctx *gin.Context, client pbv1.PostFetchServiceClient) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	res, err := client.GetTask(context.Background(), &pbv1.GetTaskRequest{
		Id: uint32(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func CreateTask(ctx *gin.Context, client pbv1.PostFetchServiceClient) {
	body := pbv1.CreateTaskRequest{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.CreateTask(context.Background(), &body)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
