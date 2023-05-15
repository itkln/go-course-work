package memorypackage

import (
	"context"
	"errors"
	"movieapp/pkg/discovery"
	"sync"
	"time"
)

type serviceNameType string
type instanceIDType string

// Registry defines an in-memory service registry
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceNameType]map[instanceIDType]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceNameType]map[instanceIDType]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		r.serviceAddrs[serviceNameType(serviceName)] = map[instanceIDType]*serviceInstance{}
	}
	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceNameType(serviceName)], instanceIDType(instanceID))
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)]; !ok {
		return errors.New("service is not registered yet")
	}
	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceNameType(serviceName)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[serviceNameType(serviceName)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
