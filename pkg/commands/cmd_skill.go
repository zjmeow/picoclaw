package commands

import (
	"context"
	"fmt"
	"strings"
)

func skillCommand() Definition {
	return Definition{
		Name:        "skill",
		Description: "Manage session-pinned skills",
		SubCommands: []SubCommand{
			{
				Name:        "add",
				Description: "Add a skill to this session",
				ArgsUsage:   "<skill>",
				Handler:     skillAddHandler(),
			},
			{
				Name:        "remove",
				Description: "Remove a session skill",
				ArgsUsage:   "<skill>",
				Handler:     skillRemoveHandler(),
			},
			{
				Name:        "clear",
				Description: "Clear all session skills",
				Handler:     skillClearHandler(),
			},
			{
				Name:        "show",
				Description: "Show active session skills",
				Handler:     skillShowHandler(),
			},
		},
	}
}

func skillAddHandler() Handler {
	return func(_ context.Context, req Request, rt *Runtime) error {
		if rt == nil || rt.ResolveSkillName == nil || rt.ListSessionSkills == nil || rt.SetSessionSkills == nil {
			return req.Reply(unavailableMsg)
		}

		rawName := nthToken(req.Text, 2)
		if strings.TrimSpace(rawName) == "" {
			return req.Reply("Usage: /skill add <skill>")
		}

		skillName, ok := rt.ResolveSkillName(rawName)
		if !ok {
			return req.Reply(fmt.Sprintf("Unknown skill: %s\nUse /list skills to see installed skills.", rawName))
		}

		current := append([]string(nil), rt.ListSessionSkills()...)
		for _, existing := range current {
			if strings.EqualFold(existing, skillName) {
				return req.Reply(fmt.Sprintf(
					"Skill %q is already active for this session.",
					skillName,
				))
			}
		}

		current = append(current, skillName)
		if err := rt.SetSessionSkills(current); err != nil {
			return req.Reply(fmt.Sprintf("Failed to update session skills: %v", err))
		}

		return req.Reply(fmt.Sprintf(
			"Skill %q is now active for this session and will be included in future prompts until removed.",
			skillName,
		))
	}
}

func skillRemoveHandler() Handler {
	return func(_ context.Context, req Request, rt *Runtime) error {
		if rt == nil || rt.ListSessionSkills == nil || rt.SetSessionSkills == nil {
			return req.Reply(unavailableMsg)
		}

		rawName := nthToken(req.Text, 2)
		if strings.TrimSpace(rawName) == "" {
			return req.Reply("Usage: /skill remove <skill>")
		}

		current := append([]string(nil), rt.ListSessionSkills()...)
		if len(current) == 0 {
			return req.Reply("No session skills are active.")
		}

		removeIndex := -1
		for i, existing := range current {
			if strings.EqualFold(existing, rawName) {
				removeIndex = i
				break
			}
		}
		if removeIndex < 0 && rt.ResolveSkillName != nil {
			if canonical, ok := rt.ResolveSkillName(rawName); ok {
				for i, existing := range current {
					if strings.EqualFold(existing, canonical) {
						removeIndex = i
						break
					}
				}
			}
		}
		if removeIndex < 0 {
			return req.Reply(fmt.Sprintf("Skill %q is not active for this session.", rawName))
		}

		removed := current[removeIndex]
		next := append(current[:removeIndex:removeIndex], current[removeIndex+1:]...)
		if err := rt.SetSessionSkills(next); err != nil {
			return req.Reply(fmt.Sprintf("Failed to update session skills: %v", err))
		}
		return req.Reply(fmt.Sprintf("Removed session skill %q.", removed))
	}
}

func skillClearHandler() Handler {
	return func(_ context.Context, req Request, rt *Runtime) error {
		if rt == nil || rt.ListSessionSkills == nil || rt.SetSessionSkills == nil {
			return req.Reply(unavailableMsg)
		}

		current := rt.ListSessionSkills()
		if len(current) == 0 {
			return req.Reply("No session skills are active.")
		}

		if err := rt.SetSessionSkills(nil); err != nil {
			return req.Reply(fmt.Sprintf("Failed to clear session skills: %v", err))
		}
		return req.Reply(fmt.Sprintf("Cleared %d session skill(s).", len(current)))
	}
}

func skillShowHandler() Handler {
	return func(_ context.Context, req Request, rt *Runtime) error {
		if rt == nil || rt.ListSessionSkills == nil {
			return req.Reply(unavailableMsg)
		}

		current := rt.ListSessionSkills()
		if len(current) == 0 {
			return req.Reply("No session skills are active.")
		}

		return req.Reply(fmt.Sprintf(
			"Session Skills:\n- %s\n\nUse /skill remove <skill> or /skill clear to stop applying them.",
			strings.Join(current, "\n- "),
		))
	}
}
