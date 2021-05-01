package adaptor

import (
	"errors"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

var _ persist.Adapter = &Adapter{}

// Adapter is the file adapter for Casbin.
// It can load policy from file or save policy to file.
type Adapter struct {
	policy map[string]string // policy_id, policy
	mu    *sync.RWMutex
	model model.Model
}

func (a *Adapter) UpdatePolicy(sec string, ptype string, oldRule, newPolicy []string) error {
	panic(errors.New("not implemented"))
}

func (a *Adapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	panic(errors.New("not implemented"))
}

// NewAdapter is the constructor for Adapter.
func NewAdapter() *Adapter {
	return &Adapter{
		policy: make(map[string]string, 0),
		mu:     &sync.RWMutex{},
	}
}

// LoadPolicy loads all policy rules from the storage into model.
func (a *Adapter) LoadPolicy(m model.Model) error {

	for _, policy := range a.policy {
		line := strings.TrimSpace(policy)

		tokens := strings.Split(line, model.DefaultSep)
		var tmp []string
		for _, token := range tokens {
			tmp = append(tmp, strings.TrimSpace(token))
		}
		tokens = tmp

		if len(tokens) <= 2 {
			continue
		}

		key := tokens[0]
		sec := key[:1]

		m[sec][key].Policy = append(m[sec][key].Policy, tokens[1:])
		m[sec][key].PolicyMap[strings.Join(tokens[1:], model.DefaultSep)] = len(m[sec][key].Policy) - 1

	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	panic(errors.New("not implemented"))
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(policyId string, ptype string, rule []string) error {
	s := ptype + "," + strings.Join(rule, ",")

	a.mu.Lock()
	a.policy[policyId] = s
	a.mu.Unlock()
	return nil
}

// AddPolicies adds policy rules to the storage.
func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	panic(errors.New("not implemented"))
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(policyId string, ptype string, rule []string) error {
	a.mu.Lock()
	delete(a.policy, policyId)
	a.mu.Unlock()
	return nil
}

// RemovePolicies removes policy rules from the storage.
func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	panic(errors.New("not implemented"))
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	panic(errors.New("not implemented"))
}
