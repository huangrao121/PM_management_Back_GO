package llm_schema

type CreateTaskSchema struct {
	Name        string `json:"name" required:"true"`
	DueDate     string `json:"due_date" required:"true"`
	AssigneeId  int    `json:"assignee_id" required:"true"`
	ProjectId   int    `json:"project_id" required:"true"`
	WorkspaceId int    `json:"workspace_id" required:"true"`
	Status      string `json:"status" required:"true"`
}

type CreateWOrkspaceSchema struct {
	Name string `json:"name" required:"true"`
}

type CreateProjectSchema struct {
	Name        string `json:"name" required:"true"`
	WorkspaceId int    `json:"workspace_id" required:"true"`
}
