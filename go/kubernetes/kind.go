package kubernetes

type NamespaceKind string
type ClusterKind string

const (
	Deployments            NamespaceKind = "deployments"
	ReplicaSets            NamespaceKind = "replicasets"
	ReplicationControllers NamespaceKind = "replicationcontrollers"
	StatefulSets           NamespaceKind = "statefulsets"
	DaemonSets             NamespaceKind = "daemonsets"
	CronJobs               NamespaceKind = "cronjobs"
	Services               NamespaceKind = "services"
	Jobs                   NamespaceKind = "jobs"
	Pods                   NamespaceKind = "pods"
	ConfigMaps             NamespaceKind = "configmaps"
	Roles                  NamespaceKind = "roles"
	RoleBindings           NamespaceKind = "rolebindings"
	NetworkPolicys         NamespaceKind = "networkpolicies"
	Ingresss               NamespaceKind = "ingresses"
	ResourceQuotas         NamespaceKind = "resourcequotas"
	LimitRanges            NamespaceKind = "limitranges"
	Secrets                NamespaceKind = "secrets"
)

const (
	ComponentStatus     ClusterKind = "componentstatuses"
	Nodes               ClusterKind = "nodes"
	Namespaces          ClusterKind = "namespaces"
	PersistentVolumes   ClusterKind = "persistentvolumes"
	ClusterRoles        ClusterKind = "clusterroles"
	ClusterRoleBindings ClusterKind = "clusterrolebindings"
	PodSecurityPolicies ClusterKind = "podsecuritypolicies"
)

func (k NamespaceKind) String() string {
	return string(k)
}

func (k ClusterKind) String() string {
	return string(k)
}

func GetNamespaceKinds() []NamespaceKind {
	return []NamespaceKind{
		Deployments,
		ReplicaSets,
		ReplicationControllers,
		StatefulSets,
		DaemonSets,
		CronJobs,
		Services,
		Jobs,
		Pods,
		ConfigMaps,
		Roles,
		RoleBindings,
		NetworkPolicys,
		Ingresss,
		ResourceQuotas,
		LimitRanges,
		Secrets,
	}
}

func GetClusterKinds() []ClusterKind {
	return []ClusterKind{
		ComponentStatus,
		Nodes,
		Namespaces,
		PersistentVolumes,
		ClusterRoles,
		ClusterRoleBindings,
		PodSecurityPolicies,
	}
}

func IsNamespaceKind(kind string) bool {
	for _, namespaceKind := range GetNamespaceKinds() {
		if namespaceKind.String() == kind {
			return true
		}
	}

	return false
}

func IsClusterKind(kind string) bool {
	for _, clusterKind := range GetClusterKinds() {
		if clusterKind.String() == kind {
			return true
		}
	}

	return false
}
