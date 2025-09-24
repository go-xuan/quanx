package nacosx

import (
	"github.com/go-xuan/quanx/serverx"
	"github.com/go-xuan/utilx/errorx"
)

// RegisterServerInstance 注册服务实例
func RegisterServerInstance(group string, server *serverx.Config) error {
	serverx.Init(&ServerCenter{})
	if err := serverx.GetCenter().Register(&ServerInstance{
		Name: server.Name,
		IP:   server.GetIP(),
		Port: server.Port,
	}); err != nil {
		return errorx.Wrap(err, "server center register instance failed")
	}
	return nil
}

// ServerCenter 服务中心
type ServerCenter struct{}

func (c *ServerCenter) Register(instance serverx.Instance) error {
	if err := this().RegisterInstance(&ServerInstance{
		Name: instance.GetName(),
		IP:   instance.GetIP(),
		Port: instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "nacos server center register instance failed")
	}
	return nil
}

func (c *ServerCenter) Deregister(instance serverx.Instance) error {
	if err := this().DeregisterInstance(&ServerInstance{
		Name: instance.GetName(),
		IP:   instance.GetIP(),
		Port: instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "nacos server center deregister instance failed")
	}
	return nil
}

func (c *ServerCenter) SelectOne(server string) (serverx.Instance, error) {
	if instance, err := this().SelectOneHealthyInstance(server); err != nil {
		return nil, errorx.Wrap(err, "nacos server center select one healthy instance failed")
	} else {
		return &ServerInstance{
			Id:   instance.InstanceId,
			Name: instance.ServiceName,
			IP:   instance.Ip,
			Port: int(instance.Port),
		}, nil
	}
}

func (c *ServerCenter) SelectList(server string) ([]serverx.Instance, error) {
	if instances, err := this().SelectInstances(server); err != nil {
		return nil, errorx.Wrap(err, "nacos server center select server instances failed")
	} else {
		var result []serverx.Instance
		for _, instance := range instances {
			result = append(result, &ServerInstance{
				Id:   instance.InstanceId,
				Name: instance.ServiceName,
				IP:   instance.Ip,
				Port: int(instance.Port),
			})
		}
		return result, nil
	}
}
