package api

import (
	"huspass/repo"
)

type Service struct {
	Repo     repo.MsisdnRepository
	UserRepo repo.UserRepo
}

//	type Router struct {
//		R *gin.Engine
//	}
//
//	type Resource struct {
//		*Router
//		*Service
//	}
//
//	func NewResource(router *Router, service *Service) *Resource {
//		return &Resource{Router: router, Service: service}
//	}
//
// // ResourceProvider provides Resource with Router and Service to do api and db operations
//
//	func ResourceProvider(router *Router, service *Service) (*Resource, error) {
//		return NewResource(router, service), nil
//	}
//
//	func (r *Router) CreateResource(msisdn *model.Msisdn) {
//		r.R.POST("api/msisdn", createMsisdn)
//	}
//
//	func createMsisdn(ctx *gin.Context) {
//		var msisdn model.Msisdn
//		err := ctx.ShouldBindJSON(&msisdn)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{
//				"message": "Provide valid parameters",
//				"err":     err.Error(),
//			})
//			return
//		}
//
//		ctx.JSON(http.StatusCreated, gin.H{
//			"message": "Successfully created resource",
//			"data":,
//		})
//	}
//
//	func (r *Router) UpdateResource(msisdn *model.Msisdn) {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func (r *Router) DeleteResource(id uuid.UUID) {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func (r *Router) GetResources() {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func (r *Router) GetResource(id uuid.UUID) {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func NewRouter(r *gin.Engine) *Router {
//		return &Router{R: r}
//	}
//
// // RouterProvider provide Router to do restful operation
//
//	func RouterProvider() (*Router, error) {
//		return NewRouter(gin.Default()), nil
//	}
//
// ServiceProvider provide Service to do db operations
func ServiceProvider(r repo.MsisdnRepository, u repo.UserRepo) (*Service, error) {
	return &Service{r, u}, nil
}

//type Endpoint interface {
//	CreateResource(msisdn *model.Msisdn)
//	UpdateResource(msisdn *model.Msisdn)
//	DeleteResource(id uuid.UUID)
//	GetResources()
//	GetResource(id uuid.UUID)
//}
