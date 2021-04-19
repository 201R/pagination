## Gorm pagination with gin

# How to use it ?

1. integrate the module to your gorm request

```go
// GetAllUsers return all Users
func GetAllUsers(p *pagination.Pagination) (*pagination.Pagination, error) {

	var Users []Users

	tx := GetDB().Preload(clause.Associations).Limit(p.Limit).Offset(p.Offset()).Order(p.Sort)

	var totalRows int64
	if err := GetDB().Model(&Users{}).Count(&totalRows).Error; err != nil {
		return nil, err
	}

	p.SetTotalRows(totalRows)

	if err := tx.Find(&Users).Error; err != nil {
		return nil, err
	}

	p.Rows = Users
	p.Paginate()

	return p, nil
}
```

2. initialize the variables

```go
    pagination.FirstPage = pagination.PageLinkFirst()
	pagination.LastPage = pagination.PageLinkLast()
	pagination.NextPage = pagination.PageLinkNext()
	pagination.PreviousPage = pagination.PageLinkPrev()

```

3. you are ok you can pass the context to your function

```go

// GetAllUsers godoc
// @Summary Get all Users
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response{data=[]models.Users}
// @Failure 500 {object} app.Response
// @Router /api/v1/Users [get]
// @tags Users
func GetAllUsers(c *gin.Context) {
	appG := app.Gin{C: c}

	pagination, err := models.GetAllUsers(pagination.New(c))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DATA, nil)
		return
	}

	// ..... ///

    // Return the response with pagination
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{
		"Users": pagination,
	})
}

```
