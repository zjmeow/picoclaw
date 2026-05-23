package agent

import (
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/session"
)

func sessionSkillStore(store session.SessionStore) session.SkillAwareSessionStore {
	if store == nil {
		return nil
	}
	skillStore, ok := store.(session.SkillAwareSessionStore)
	if !ok {
		return nil
	}
	return skillStore
}

func sessionSkillNames(store session.SessionStore, sessionKey string) []string {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return nil
	}
	skillStore := sessionSkillStore(store)
	if skillStore == nil {
		return nil
	}
	skills := skillStore.GetSessionSkills(sessionKey)
	if len(skills) == 0 {
		return nil
	}
	return append([]string(nil), skills...)
}

func (al *AgentLoop) getSessionSkills(agent *AgentInstance, sessionKey string) []string {
	if agent == nil {
		return nil
	}
	return sessionSkillNames(agent.Sessions, sessionKey)
}

func (al *AgentLoop) setSessionSkills(agent *AgentInstance, sessionKey string, skillNames []string) error {
	if agent == nil || agent.Sessions == nil {
		return fmt.Errorf("session skill selection is unavailable")
	}
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return fmt.Errorf("session skill selection is unavailable")
	}
	skillStore := sessionSkillStore(agent.Sessions)
	if skillStore == nil {
		return fmt.Errorf("session skill selection is unavailable")
	}
	skillStore.SetSessionSkills(sessionKey, skillNames)
	if err := agent.Sessions.Save(sessionKey); err != nil {
		return err
	}
	return nil
}

func sessionPureStore(store session.SessionStore) session.PureAwareSessionStore {
	if store == nil {
		return nil
	}
	pureStore, ok := store.(session.PureAwareSessionStore)
	if !ok {
		return nil
	}
	return pureStore
}

func sessionPureConfig(store session.SessionStore, sessionKey string) session.PureSessionConfig {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return session.PureSessionConfig{}
	}
	pureStore := sessionPureStore(store)
	if pureStore == nil {
		return session.PureSessionConfig{}
	}
	return pureStore.GetSessionPureConfig(sessionKey)
}

func (al *AgentLoop) getSessionPureConfig(agent *AgentInstance, sessionKey string) session.PureSessionConfig {
	if agent == nil {
		return session.PureSessionConfig{}
	}
	return sessionPureConfig(agent.Sessions, sessionKey)
}

func (al *AgentLoop) setSessionPureConfig(agent *AgentInstance, sessionKey string, cfg session.PureSessionConfig) error {
	if agent == nil || agent.Sessions == nil {
		return fmt.Errorf("pure session configuration is unavailable")
	}
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return fmt.Errorf("pure session configuration is unavailable")
	}
	pureStore := sessionPureStore(agent.Sessions)
	if pureStore == nil {
		return fmt.Errorf("pure session configuration is unavailable")
	}
	cfg.Skill = strings.TrimSpace(cfg.Skill)
	pureStore.SetSessionPureConfig(sessionKey, cfg)
	if err := agent.Sessions.Save(sessionKey); err != nil {
		return err
	}
	return nil
}
