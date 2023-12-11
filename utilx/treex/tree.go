package treex

// 树形结构
type NodeList[T any] []*Node[T]
type Node[T any] struct {
	Id    string      `json:"id"`
	Pid   string      `json:"pid"`
	Data  T           `json:"data"`
	Child NodeList[T] `json:"child"`
}

// 一维数组转树形结构
func (list NodeList[T]) Convert2Tree(rootId ...string) NodeList[T] {
	var pid = "0"
	if len(rootId) > 0 {
		pid = rootId[0]
	}
	nodeMap := make(map[string]NodeList[T])
	for _, item := range list {
		nodeMap[item.Pid] = append(nodeMap[item.Pid], item)
	}
	if result, ok := nodeMap[pid]; ok {
		result = result.findChildFrom(nodeMap)
		return result
	}
	return nil
}

// 获取子节点
func (list NodeList[T]) findChildFrom(nodeMap map[string]NodeList[T]) (result NodeList[T]) {
	if list != nil && len(list) > 0 {
		for _, item := range list {
			var child = nodeMap[item.Id]
			child = child.findChildFrom(nodeMap)
			result = append(result, &Node[T]{item.Id, item.Pid, item.Data, child})
		}
	}
	return
}
