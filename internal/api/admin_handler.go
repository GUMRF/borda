package api

import (
	"borda/internal/domain"
	"borda/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAdminRoutes(router fiber.Router) {
	admin := router.Group("/admin", h.authRequired, h.adminPermissionRequired)

	tasks := admin.Group("/tasks")
	tasks.Get("", h.adminGetAllTasks)
	tasks.Post("", h.createNewTask)
	tasks.Patch("/:id", h.updateTask)
}

// @Summary      Get all tasks
// @Description  Get all tasks with admin access.
// @Tags         Admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}   domain.Task
// @Failure      400   {object}  ErrorsResponse
// @Failure      404   {object}  ErrorsResponse
// @Failure      500   {object}  ErrorsResponse
// @Router       /admin/tasks [get]
func (h *Handler) adminGetAllTasks(c *fiber.Ctx) error {
	uc := usecase.NewAdminUsecaseGetAlTasks(h.Repository.Tasks)

	tasks, err := uc.Execute()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Internal server error", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}

// @Summary      Update task
// @Description  Update task.
// @Tags         Admin
// @Accept       json
// @Security     ApiKeyAuth
// @Produce      json
// @Param        task_id  path      int          true  "Task ID"
// @Param        task     body      domain.Task  true  "Task"
// @Success      200      string    OK
// @Failure      400      {object}  ErrorsResponse
// @Failure      404      {object}  ErrorsResponse
// @Failure      500      {object}  ErrorsResponse
// @Router       /admin/tasks/{task_id} [patch]
func (h *Handler) updateTask(c *fiber.Ctx) error {
	// Get task id from request url
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusConflict, "Parse url string", err.Error())
	}

	// Get request payload
	var update domain.TaskUpdate
	if err := c.BodyParser(&update); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid.",
			err.Error())
	}

	//TODO: Fix validation: ints
	// Validate request payload
	if err := update.Validate(); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest,
			"Validation is not passed.", err.Error())
	}

	//Init usecase
	uc := usecase.NewAdminUsecaseUpdateTask(h.Repository.Tasks)

	// Update task
	updatedTask, err := uc.Execute(taskId, update)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest,
			"Internal server error.", err.Error())
	}

	// TODO: fetch and return task after update

	return c.Status(fiber.StatusOK).JSON(updatedTask)
}

// @Summary      Create new task
// @Description  Create new task.
// @Tags         Admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        task  body      domain.Task  true  "Task"
// @Success      200   {object}  domain.Task
// @Failure      400  {object}  ErrorsResponse
// @Failure      404  {object}  ErrorsResponse
// @Failure      500  {object}  ErrorsResponse
// @Router       /admin/tasks [post]
func (h *Handler) createNewTask(ctx *fiber.Ctx) error {
	var task domain.Task

	if err := ctx.BodyParser(&task); err != nil {
		return NewErrorResponse(ctx, fiber.StatusBadRequest,
			"Input is invalid.", err.Error())
	}

	//TODO: Fix validator: id, author id, points, isActive, isDisable
	//if err := task.Validate(); err != nil {
	//	return NewErrorResponse(ctx, fiber.StatusBadRequest,
	//		"Validation is not passed.", err.Error())
	//}

	uc := usecase.NewAdminUsecaseCreateNewTask(h.Repository.Tasks)

	createdTask, err := uc.Execute(task)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Internal server error.", err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(createdTask)
}
