package kubernetes

import apiv1 "k8s.io/api/core/v1"

func FilterOutNodesWithIgnoredLabel(ignoredLabel map[string]string, nodes []*apiv1.Node) []*apiv1.Node {
	retNode := make([]*apiv1.Node, 0)
	for _, node := range nodes {
		filter := false
		for label := range ignoredLabel {
			if _, hasLabel := node.Labels[label]; hasLabel {
				filter = true
				break
			}
		}
		if !filter {
			retNode = append(retNode, node)
		}
	}
	return retNode
}