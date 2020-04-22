package kubernetes

import (
	"k8s.io/apimachinery/pkg/labels"
	"time"

	apiv1 "k8s.io/api/core/v1"
	v1lister "k8s.io/client-go/listers/core/v1"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// ListerRegistry is a registry providing various listers
type ListerRegistry interface {
	AllNodeLister() NodeLister
	ReadyNodeLister() NodeLister
}

type listerRegistryImpl struct {
	allNodeLister   NodeLister
	readyNodeLister NodeLister
}


func NewListerRegistry(allNode NodeLister, readyNode NodeLister) ListerRegistry {

}

func NewListerRegistryWithDefaultListers(kubeClient client.Interface, stopChannel <-chan struct{}) ListerRegistry {
	allNodeLister := NewAllNodeLister(kubeClient, stopChannel)
	readyNodeLister := NewReadyNodeLister(kubeClient, stopChannel)
	return NewListerRegistry(allNodeLister, readyNodeLister)
}


// NodeLister lists nodes.
type NodeLister interface {
	List() ([]*apiv1.Node, error)
	Get(name string) (*apiv1.Node, error)
}

// nodeLister implementation.
type nodeListerImpl struct {
	nodeLister v1lister.NodeLister
	filter     func(*apiv1.Node) bool
}


// NewNodeLister builds a node lister.
func NewNodeLister(kubeClient client.Interface, filter func(*apiv1.Node) bool, stopChannel <-chan struct{}) NodeLister {
	listWatcher := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "nodes", apiv1.NamespaceAll, fields.Everything())
	store, reflector := cache.NewNamespaceKeyedIndexerAndReflector(listWatcher, &apiv1.Node{}, time.Hour)
	nodeLister := v1lister.NewNodeLister(store)
	go reflector.Run(stopChannel)
	return &nodeListerImpl{
		nodeLister: nodeLister,
		filter:     filter,
	}
}

// List returns list of nodes.
func (l *nodeListerImpl) List() ([]*apiv1.Node, error) {
	var nodes []*apiv1.Node
	var err error

	nodes, err = l.nodeLister.List(labels.Everything())
	if err != nil {
		return []*apiv1.Node{}, err
	}

	if l.filter != nil {
		nodes = filterNodes(nodes, l.filter)
	}

	return nodes, nil
}

// Get returns the node with the given name.
func (l *nodeListerImpl) Get(name string) (*apiv1.Node, error) {
	node, err := l.nodeLister.Get(name)
	if err != nil {
		return nil, err
	}
	return node, nil
}


func filterNodes(nodes []*apiv1.Node, predicate func(*apiv1.Node) bool) []*apiv1.Node {
	var filtered []*apiv1.Node
	for i := range nodes {
		if predicate(nodes[i]) {
			filtered = append(filtered, nodes[i])
		}
	}
	return filtered
}

func NewAllNodeLister(kubeClient client.Interface, stopChannel <-chan struct{}) NodeLister {
	return NewNodeLister(kubeClient, nil, stopChannel)
}
