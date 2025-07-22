package nacosx

import (
	"github.com/go-xuan/quanx/serverx"
	"github.com/go-xuan/utilx/errorx"
)

// RegisterServerInstance 注册服务实例
func RegisterServerInstance(server *serverx.Config) error {
	group := server.NacosGroup()
	center := &ServerCenter{group: group}
	if err := center.Register(&ServerInstance{
		Name:  server.Name,
		Group: server.NacosGroup(),
		IP:    server.IP,
		Port:  server.Port,
	}); err != nil {
		return errorx.Wrap(err, "register server instance failed")
	}
	serverx.Init(center)
	return nil
}

// ServerCenter 服务中心
type ServerCenter struct {
	group string
}

func (c *ServerCenter) Register(instance serverx.Instance) error {
	if err := this().RegisterInstance(&ServerInstance{
		Group: c.group,
		Name:  instance.GetName(),
		IP:    instance.GetIP(),
		Port:  instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "nacos instance register failed")
	}
	return nil
}

func (c *ServerCenter) Deregister(instance serverx.Instance) error {
	if err := this().DeregisterInstance(&ServerInstance{
		Group: c.group,
		Name:  instance.GetName(),
		IP:    instance.GetIP(),
		Port:  instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "nacos instance deregister failed")
	}
	return nil
}

func (c *ServerCenter) SelectOne(server string) (serverx.Instance, error) {
	if instance, err := this().SelectOneHealthyInstance(server, c.group); err != nil {
		return nil, errorx.Wrap(err, "select one healthy instance failed")
	} else {
		return &ServerInstance{
			Id:    instance.InstanceId,
			Name:  instance.ServiceName,
			Group: c.group,
			IP:    instance.Ip,
			Port:  int(instance.Port),
		}, nil
	}
}

func (c *ServerCenter) SelectList(server string) ([]serverx.Instance, error) {
	if instances, err := this().SelectInstances(server, c.group); err != nil {
		return nil, errorx.Wrap(err, "select instances failed")
	} else {
		var result []serverx.Instance
		for _, instance := range instances {
			result = append(result, &ServerInstance{
				Id:    instance.InstanceId,
				Name:  instance.ServiceName,
				Group: c.group,
				IP:    instance.Ip,
				Port:  int(instance.Port),
			})
		}
		return result, nil
	}
}
