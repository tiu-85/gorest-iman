package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
)

func GetPostList(ctx *gin.Context, client pbv1.PostCrudServiceClient) {
	body := pbv1.GetPostListRequest{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.GetPostList(context.Background(), &body)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func GetPost(ctx *gin.Context, client pbv1.PostCrudServiceClient) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	res, err := client.GetPost(context.Background(), &pbv1.GetPostRequest{
		Id: uint32(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func EditPost(ctx *gin.Context, client pbv1.PostCrudServiceClient) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	body := pbv1.EditPostRequestParams{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.EditPost(context.Background(), &pbv1.EditPostRequest{
		Id:     uint32(id),
		Params: &body,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func DeletePost(ctx *gin.Context, client pbv1.PostCrudServiceClient) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	res, err := client.DeletePost(context.Background(), &pbv1.DeletePostRequest{
		Id: uint32(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
