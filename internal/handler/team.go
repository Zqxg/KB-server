package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/service/team"
	"projectName/pkg/utils"
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
// @Router /v1/team/createTeam [post]
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
// @Router /v1/team/updateTeam [post]
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
// @Router /v1/team/deleteTeam [post]
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
// @Router /v1/team/getTeamList [get]
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

// GetTeamInfo godoc
// @Summary 获取团队详细
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query int true "Team ID"
// @Success 200 {object} v1.GetTeamInfoResp
// @Router /v1/team/getTeamInfo [get]
func (h *TeamHandler) GetTeamInfo(ctx *gin.Context) {
	// 从查询参数中获取参数
	if !utils.IsNumeric(ctx.Query("id")) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	teamID, _ := utils.ToInt(ctx.Query("id")) // 获取 teamID 参数
	if teamID < 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	teamInfoResp, err := h.teamService.GetTeamInfo(ctx, teamID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, teamInfoResp)
}

// GetUserTeamList godoc
// @Summary 获取个人团队列表
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param pageIndex query int true "Page Index"
// @Param pageSize query int true "Page Size"
// @Success 200 {object} v1.GetUserTeamListResp
// @Router /v1/team/getUserTeamList [GET]
func (h *TeamHandler) GetUserTeamList(ctx *gin.Context) {
	// 从查询参数中获取参数
	if !utils.IsNumeric(ctx.Query("pageIndex")) || !utils.IsNumeric(ctx.Query("pageSize")) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	pageIndex, _ := utils.ToInt(ctx.Query("pageIndex")) // 获取参数
	pageSize, _ := utils.ToInt(ctx.Query("pageSize"))   // 获取参数
	if pageIndex < 0 || pageSize < 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	teamListResp, err := h.teamService.GetUserTeamList(ctx, userId, pageIndex, pageSize)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, teamListResp)
}

// GetTeamMemberList godoc
// @Summary 获取团队成员列表
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query int true "Team ID"
// @Success 200 {object} v1.GetTeamMemberListResp
// @Router /v1/team/getTeamMemberList [get]
func (h *TeamHandler) GetTeamMemberList(ctx *gin.Context) {
	// 从查询参数中获取参数
	if !utils.IsNumeric(ctx.Query("id")) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	teamID, _ := utils.ToInt(ctx.Query("id")) // 获取 teamID 参数
	if teamID < 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	teamMemberList, err := h.teamService.GetTeamMemberList(ctx, teamID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.GetTeamMemberListResp{
		MemberList: teamMemberList,
	})

}

// AddTeamMember godoc
// @Summary 添加团队成员
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.AddTeamMemberReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/team/addTeamMember [post]
func (h *TeamHandler) AddTeamMember(ctx *gin.Context) {
	var req v1.AddTeamMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.teamService.AddTeamMember(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteTeamMember godoc
// @Summary 删除团队成员
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.DeleteTeamMemberReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/team/deleteTeamMember [post]
func (h *TeamHandler) DeleteTeamMember(ctx *gin.Context) {
	var req v1.DeleteTeamMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.teamService.DeleteTeamMember(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UpdateTeamMemberRole godoc
// @Summary 更新团队成员角色
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateTeamMemberRoleReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/team/updateTeamMemberRole [post]
func (h *TeamHandler) UpdateTeamMemberRole(ctx *gin.Context) {
	var req v1.UpdateTeamMemberRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.teamService.UpdateTeamMemberRole(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// QuitTeam godoc
// @Summary 退出团队
// @Schemes
// @Description
// @Tags 团队模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.QuitTeamReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/team/quitTeam [post]
func (h *TeamHandler) QuitTeam(ctx *gin.Context) {
	var req v1.QuitTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.teamService.QuitTeam(ctx, userId, req.TeamID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
