package kubernetes

import (
	"encoding/json"

	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
)

func StructedPod(data []byte) (*coreV1.Pod, error) {
	pod := &coreV1.Pod{}
	err := json.Unmarshal(data, pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func StructedService(data []byte) (*coreV1.Service, error) {
	service := &coreV1.Service{}
	err := json.Unmarshal(data, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func StructedNode(data []byte) (*coreV1.Node, error) {
	node := &coreV1.Node{}
	err := json.Unmarshal(data, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func StructedRole(data []byte) (*rbacV1.Role, error) {
	role := &rbacV1.Role{}
	err := json.Unmarshal(data, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func StructedRoleBinding(data []byte) (*rbacV1.RoleBinding, error) {
	roleBinding := &rbacV1.RoleBinding{}
	err := json.Unmarshal(data, roleBinding)
	if err != nil {
		return nil, err
	}
	return roleBinding, nil
}

func StructedClusterRole(data []byte) (*rbacV1.ClusterRole, error) {
	clusterRole := &rbacV1.ClusterRole{}
	err := json.Unmarshal(data, clusterRole)
	if err != nil {
		return nil, err
	}
	return clusterRole, nil
}

func StructedClusterRoleBinding(data []byte) (*rbacV1.ClusterRoleBinding, error) {
	clusterRoleBinding := &rbacV1.ClusterRoleBinding{}
	err := json.Unmarshal(data, clusterRoleBinding)
	if err != nil {
		return nil, err
	}
	return clusterRoleBinding, nil
}
