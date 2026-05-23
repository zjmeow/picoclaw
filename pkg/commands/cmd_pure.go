package commands

import (
	"context"
	"fmt"
	"strings"
)

func pureCommand() Definition {
	return Definition{
		Name:        "pure",
		Description: "Use a minimal pure chat session",
		Usage:       "/pure [skill|off]",
		Handler:     pureHandler(),
	}
}

func pureHandler() Handler {
	return func(_ context.Context, req Request, rt *Runtime) error {
		if rt == nil || rt.SetSessionPure == nil {
			return req.Reply(unavailableMsg)
		}

		rawName := nthToken(req.Text, 1)
		if strings.EqualFold(rawName, "off") || strings.EqualFold(rawName, "clear") {
			if err := rt.SetSessionPure(PureSessionConfig{}); err != nil {
				return req.Reply(fmt.Sprintf("Failed to disable pure session: %v", err))
			}
			return req.Reply("Pure session disabled.")
		}

		cfg := PureSessionConfig{Enabled: true}
		if strings.TrimSpace(rawName) != "" {
			if rt.ResolveSkillName == nil {
				return req.Reply(unavailableMsg)
			}
			skillName, ok := rt.ResolveSkillName(rawName)
			if !ok {
				return req.Reply(fmt.Sprintf("Unknown skill: %s\nUse /list skills to see installed skills.", rawName))
			}
			cfg.Skill = skillName
		}

		if err := rt.SetSessionPure(cfg); err != nil {
			return req.Reply(fmt.Sprintf("Failed to update pure session: %v", err))
		}
		if cfg.Skill == "" {
			return req.Reply("Pure session enabled without a skill. Only dynamic context will be included in the system prompt.")
		}
		return req.Reply(fmt.Sprintf("Pure session enabled with skill %q. Default tools are disabled until you run /tools add default.", cfg.Skill))
	}
}
