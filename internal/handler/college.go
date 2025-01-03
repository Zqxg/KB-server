package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/service"
)

type CollegeHandler struct {
	*Handler
	collegeService service.CollegeService
}

func NewCollegeHandler(
	handler *Handler,
	collegeService service.CollegeService,
) *CollegeHandler {
	return &CollegeHandler{
		Handler:        handler,
		collegeService: collegeService,
	}
}

// GetCollege 处理单个学院信息
func (h *CollegeHandler) GetCollege(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	var req v1.GetCollegeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrBadRequest, nil)
		return
	}
	college, err := h.collegeService.GetCollege(ctx, req.CollegeId)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	// 返回单个学院信息
	v1.HandleSuccess(ctx, v1.CollegeResponseData{
		CollegeId:   college.CollegeId,
		CollegeName: college.CollegeName,
		Description: college.Description,
	})
}

// GetCollegeList 获取所有学院信息
func (h *CollegeHandler) GetCollegeList(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	collegeList, err := h.collegeService.GetCollegeList(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	// 格式化多个学院的信息
	var responseCollegeList []*v1.CollegeResponseData
	for _, college := range collegeList {
		responseCollegeList = append(responseCollegeList, &v1.CollegeResponseData{
			CollegeId:   college.CollegeId,
			CollegeName: college.CollegeName,
			Description: college.Description,
		})
	}
	// 返回多个学院信息
	v1.HandleSuccess(ctx, map[string]interface{}{
		"collegeList": responseCollegeList,
	})
}
