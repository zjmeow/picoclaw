package commands

import (
	"context"
	"fmt"
	"strings"
)

func toolsCommand() Definition {
	return Definition{
		Name:        "tools",
		Description: "Manage session tool injection",
		SubCommands: []SubCommand{
			{
				Name:        "add",
				Description: "Add tools to this session",
				ArgsUsage:   "default",
				Handler:     toolsAddHandler(),
			},
		},
	}
}

func toolsAddHandler() Handler {
	return func(_ context.Context, req Request, rt *Runtime) error {
		if rt == nil || rt.GetSessionPure == nil || rt.SetSessionPure == nil {
			return req.Reply(unavailableMsg)
		}

		target := nthToken(req.Text, 2)
		if !strings.EqualFold(target, "default") {
			return req.Reply("Usage: /tools add default")
		}

		cfg := rt.GetSessionPure()
		if !cfg.Enabled {
			return req.Reply("Default tools can only be added after enabling a pure session with /pure.")
		}
		if cfg.DefaultTools {
			return req.Reply("Default tools are already enabled for this pure session.")
		}

		cfg.DefaultTools = true
		if err := rt.SetSessionPure(cfg); err != nil {
			return req.Reply(fmt.Sprintf("Failed to update pure session tools: %v", err))
		}
		return req.Reply("Default tools are now enabled for this pure session.")
	}
}
