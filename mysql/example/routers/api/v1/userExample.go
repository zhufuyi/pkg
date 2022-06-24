package v1

/*
// CreateUserExample 创建
func CreateUserExample(c *gin.Context) {
	form := &service.CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.CreateUserExample(form)
	if err != nil {
		logger.Error("CreateUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.UserExampleCreateErr)
		return
	}

	render.Success(c)
}

// DeleteUserExample 删除一条记录
func DeleteUserExample(c *gin.Context) {
	form := &service.DeleteUserExampleRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.DeleteUserExample(form.ID)
	if err != nil {
		logger.Error("DeleteUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.UserExampleDeleteErr)
		return
	}

	render.Success(c)
}

// DeleteUserExamples 删除多条记录
func DeleteUserExamples(c *gin.Context) {
	form := &service.DeleteUserExamplesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.DeleteUserExample(form.IDs...)
	if err != nil {
		logger.Error("DeleteUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.UserExampleDeleteErr)
		return
	}

	render.Success(c)
}

// UpdateUserExample 更新
func UpdateUserExample(c *gin.Context) {
	form := &service.UpdateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.UpdateUserExample(form)
	if err != nil {
		logger.Error("CreateUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.UserExampleUpdateErr)
		return
	}

	render.Success(c)
}

// GetUserExample 获取一条记录
func GetUserExample(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	form := &service.GetUserExampleRequest{ID: id}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	userExample, err := svc.GetUserExample(form)
	if err != nil {
		logger.Error("GetUserExample error", logger.Err(err), logger.Any("form", form))
		if err.Error() == mysql.ErrNotFound.Error() {
			render.Error(c, errcode.NotFound)
		} else {
			render.Error(c, errcode.UserExampleGetErr)
		}
		return
	}

	render.Success(c, gin.H{"userExample": userExample})
}

// GetUserExamples 获取多条记录
func GetUserExamples(c *gin.Context) {
	form := &service.GetUserExamplesRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	users, total, err := svc.GetUserExamples(form)
	if err != nil {
		logger.Error("GetUserExample error", logger.Err(err), logger.Any("form", form))
		if err.Error() == mysql.ErrNotFound.Error() {
			render.Error(c, errcode.NotFound)
		} else {
			render.Error(c, errcode.UserExampleGetErr)
		}
		return
	}

	render.Success(c, gin.H{
		"userExamples": users,
		"total":     total,
	})
}
*/
