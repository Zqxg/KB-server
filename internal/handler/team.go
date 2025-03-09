package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/service/team"
)

type TeamHandler struct {
	*Handler
	teamService team.TeamService
}

func NewTeamHandler(
	handler *Handler,
	teamService team.TeamService,
) *TeamHandler {
	return &TeamHandler{
		Handler:     handler,
		teamService: teamService,
	}
}

// CreateTeam godoc
// @Summary 新建团队
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateTeamRequest true "params"
// @Success 200 {object} v1.CreateTeamResp
// @Router /team/createTeam [post]
func (h *TeamHandler) CreateTeam(ctx *gin.Context) {
	var req v1.CreateTeamRequest
	userId := GetUserIdFromCtx(ctx)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	teamId, err := h.teamService.CreateTeam(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.CreateTeamResp{
		TeamID: teamId,
	})
}

// UpdateTeam godoc
// @Summary 更新团队
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateTeamRequest true "params"
// @Success 200 {object} v1.Response
// @Router /team/updateTeam [post]
func (h *TeamHandler) UpdateTeam(ctx *gin.Context) {
	var req v1.UpdateTeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.teamService.UpdateTeam(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteTeam godoc
// @Summary 删除团队
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.DeleteTeamRequest true "params"
// @Success 200 {object} v1.Response
// @Router /team/deleteTeam [post]
func (h *TeamHandler) DeleteTeam(ctx *gin.Context) {
	var req v1.DeleteTeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.teamService.DeleteTeam(ctx, userId, req.TeamID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetTeamList godoc
// @Summary 获取团队列表
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.GetTeamListReq true "params"
// @Success 200 {object} v1.GetTeamListResp
// @Router /team/getTeamList [get]
func (h *TeamHandler) GetTeamList(ctx *gin.Context) {
	var req v1.GetTeamListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	teamListResp, err := h.teamService.GetTeamList(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, teamListResp)
}
