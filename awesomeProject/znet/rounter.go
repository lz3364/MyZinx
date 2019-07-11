package znet


/*
 * 作为路由的基类,当有其他的路由类时,嵌入基类,并实现接口的方法
 */
type BaseRounter struct {

}

/*
 * 基类的3个方法都为空,为了方便其他路由类从基类继承时,不需要实现PreHandle 及 PostHandle等方法
 * 不是所有的Rounter都需要这两个方法
 */
func (br *BaseRounter) PreHandle(request Request)  {

}

func (br *BaseRounter) Handle(request Request)  {

}

func (br *BaseRounter) PostHandle(request Request)  {

}
