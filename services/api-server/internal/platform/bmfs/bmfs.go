package bmfs

import (
	"fmt"
	"time"

	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

type GuardFunc func(ctx ActionContext) bool

type OnTransitionFunc func(ctx ActionContext, from, to, action string)

type Transition struct {
	From         string
	Action       string
	To           string
	Guard        GuardFunc
	OnTransition OnTransitionFunc
}

type Machine struct {
	name        string
	states      map[string]bool
	transitions []Transition
	index       map[string][]int
}

type ActionContext struct {
	Principal auth.Principal
	IsOwner   bool
	UserRole  string
	Now       time.Time
	Extra     map[string]any
}

type ExecuteResult struct {
	FromStatus string
	ToStatus   string
	Action     string
	CanActions map[string]bool
}

func NewMachine(name string) *Machine {
	return &Machine{
		name:   name,
		states: make(map[string]bool),
		index:  make(map[string][]int),
	}
}

func (m *Machine) AddState(state string) *Machine {
	m.states[state] = true
	return m
}

func (m *Machine) AddTransition(from, action, to string, guard GuardFunc, onTransition OnTransitionFunc) *Machine {
	idx := len(m.transitions)
	m.transitions = append(m.transitions, Transition{
		From:         from,
		Action:       action,
		To:           to,
		Guard:        guard,
		OnTransition: onTransition,
	})
	m.index[from] = append(m.index[from], idx)
	m.states[from] = true
	m.states[to] = true
	return m
}

func (m *Machine) Name() string {
	return m.name
}

func (m *Machine) Execute(currentStatus, action string, actx ActionContext) (*ExecuteResult, error) {
	indices, ok := m.index[currentStatus]
	if !ok {
		return nil, fmt.Errorf("bmfs: unknown status %q in machine %q", currentStatus, m.name)
	}

	for _, idx := range indices {
		t := m.transitions[idx]
		if t.Action != action {
			continue
		}
		if t.Guard != nil && !t.Guard(actx) {
			continue
		}
		if t.OnTransition != nil {
			t.OnTransition(actx, currentStatus, t.To, action)
		}
		result := &ExecuteResult{
			FromStatus: currentStatus,
			ToStatus:   t.To,
			Action:     action,
			CanActions: m.availableActions(t.To, actx),
		}
		return result, nil
	}

	return nil, fmt.Errorf("bmfs: action %q not allowed from status %q in machine %q", action, currentStatus, m.name)
}

func (m *Machine) AvailableActions(currentStatus string, actx ActionContext) map[string]bool {
	return m.availableActions(currentStatus, actx)
}

func (m *Machine) availableActions(currentStatus string, actx ActionContext) map[string]bool {
	result := make(map[string]bool)
	indices, ok := m.index[currentStatus]
	if !ok {
		return result
	}
	for _, idx := range indices {
		t := m.transitions[idx]
		allowed := true
		if t.Guard != nil {
			allowed = t.Guard(actx)
		}
		result["can_"+t.Action] = allowed
	}
	return result
}
