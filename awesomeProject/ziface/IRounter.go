package ziface

/*
 * 使用zinx框架的人给链接自定义的处理业务的方法
 * IRequest里面包含了链接的信息以及链接的请求数据
 * todo: 需要去了解什么叫路由
 */
type IRounter interface {

	//处理conn业务之前的处理方法
	PreHandle(request IRequest)

	//处理conn业务的处理方法
	Handle(request IRequest)

	//处理conn之后的钩子方法
	PostHandle(request IRequest)
}