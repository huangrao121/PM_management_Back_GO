package router

import (
	"pm_go_version/config"

	"github.com/gin-gonic/gin"
)

func Init(init *config.Initialization) *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		me := api.Group("/me")
		me.GET("/", AuthMiddleware(), init.Uc.GetMe)
		user := api.Group("/user")
		user.POST("/login", init.Uc.LoginUser)
		user.POST("/signup", init.Uc.SignupUser)
		user.GET("/current", AuthMiddleware(), init.Uc.GetUserInfo)

		workspace := api.Group("/workspace")
		workspace.GET("/all", init.Wc.GetListofWorkspaces)
		workspace.POST("/", AuthMiddleware(), init.Wc.CreateWorkspace)
		workspace.GET("/", AuthMiddleware(), init.Wc.GetWorkspacesById)
		workspace.DELETE("/:workspaceId", AuthMiddleware(), init.Wc.DeleteWorkspaceById)
		workspace.PATCH("/:workspaceId", AuthMiddleware(), init.Wc.UpdateWorkspaceById)
		workspace.GET("/:workspaceId", AuthMiddleware(), init.Wc.GetSingleWorkspaceById)
		workspace.PATCH("/:workspaceId/reset-invite-code", AuthMiddleware(), init.Wc.ResetInvite)
		workspace.POST("/:workspaceId/join", AuthMiddleware(), init.Wc.JoinWorkspace)
		workspace.GET("/:workspaceId/info", AuthMiddleware(), init.Wc.GetWorkspaceInfo)

		members := api.Group("/members")
		members.GET("/workspace/:workspaceId", AuthMiddleware(), init.Uwc.GetListofMembersByWorkspaceId)
		members.DELETE("/:memberId/workspace/:workspaceId", AuthMiddleware(), init.Uwc.DeleteMemberByWorkspaceId)

		projects := api.Group("/projects")
		projects.GET("/:workspaceId", AuthMiddleware(), init.Pc.GetListofProjects)
		projects.POST("/", AuthMiddleware(), init.Pc.CreateProject)
		projects.GET("/project/:projectId", AuthMiddleware(), init.Pc.GetProjectById)
		projects.PATCH("/:workspaceId/:projectId", AuthMiddleware(), init.Pc.UpdateProjectById)

		tasks := api.Group("/tasks")
		tasks.GET("/workspace/:workspaceId", AuthMiddleware(), init.Tc.GetListofTasks)
		tasks.POST("/", AuthMiddleware(), init.Tc.CreateTask)
		tasks.GET("/:taskId", AuthMiddleware(), init.Tc.GetTaskById)
		tasks.PATCH("/:taskId", AuthMiddleware(), init.Tc.UpdateTaskById)
		tasks.DELETE("/:taskId", AuthMiddleware(), init.Tc.DeleteTaskById)
		tasks.PATCH("/batchUpdate", AuthMiddleware(), init.Tc.BatchUpdateTask)
	}
	return router
}
