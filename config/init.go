package config

import (
	"pm_go_version/app/controller"
	//"pm_go_version/app/repository"
	//"pm_go_version/app/service"
)

type Initialization struct {
	Uc controller.UserController
	//Us service.UserService
	//Ur repository.UserRepository
	Wc controller.WorkspaceController
	//Ws service.WorkspaceService
	//Wr repository.WorkspaceRepository
	Uwc controller.UserWorkspaceController
	Pc  controller.ProjectController
	Tc  controller.TaskController
}

func NewInitialization(uc controller.UserController,
	//us service.UserService,
	//ur repository.UserRepository,
	wc controller.WorkspaceController,
	//ws service.WorkspaceService,
	//wr repository.WorkspaceRepository,
	uwc controller.UserWorkspaceController,
	pc controller.ProjectController,
	tc controller.TaskController,
) *Initialization {
	return &Initialization{
		Uc: uc,
		//Us: us,
		//Ur: ur,
		Wc: wc,
		//Ws: ws,
		//Wr: wr,
		Uwc: uwc,
		Pc:  pc,
		Tc:  tc,
	}
}
