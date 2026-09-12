package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/FyberLabs/hypermesh-cli/internal/api"
)

func newChatCmd(r *run) *cobra.Command {
	var leaseID, model, message, system string
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "POST router /v1/chat/completions (never the control-plane 409 stub)",
		RunE: func(cmd *cobra.Command, args []string) error {
			text, err := messageOrStdin(message, args)
			if err != nil {
				return err
			}
			return runChat(r, cmd, leaseID, model, system, text)
		},
	}
	addChatFlags(cmd, &leaseID, &model, &message, &system)
	return cmd
}

func newPromptCmd(r *run) *cobra.Command {
	var leaseID, model, message, system string
	cmd := &cobra.Command{
		Use:   "prompt [text...]",
		Short: "One-shot Full Model prompt on the Fyber router",
		RunE: func(cmd *cobra.Command, args []string) error {
			text, err := messageOrStdin(message, args)
			if err != nil {
				return err
			}
			return runChat(r, cmd, leaseID, model, system, text)
		},
	}
	addChatFlags(cmd, &leaseID, &model, &message, &system)
	return cmd
}

func newCompletionsCmd(r *run) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completions",
		Short: "OpenAI-shaped chat completions against the Fyber router",
	}
	var leaseID, model, message, system string
	create := &cobra.Command{
		Use:   "create",
		Short: "POST https://chat.test.hyperme.sh/v1/chat/completions",
		RunE: func(cmd *cobra.Command, args []string) error {
			text, err := messageOrStdin(message, args)
			if err != nil {
				return err
			}
			return runChat(r, cmd, leaseID, model, system, text)
		},
	}
	addChatFlags(create, &leaseID, &model, &message, &system)
	cmd.AddCommand(create)
	return cmd
}

func addChatFlags(cmd *cobra.Command, leaseID, model, message, system *string) {
	cmd.Flags().StringVar(leaseID, "lease-id", "", "paid lease ticket (also HYPERMESH_LEASE_ID)")
	cmd.Flags().StringVar(model, "model", api.DefaultCatalogID, "OpenAI-shaped model name (Phase 1 catalog id)")
	cmd.Flags().StringVar(message, "message", "", "user message (omit to read remaining args or stdin)")
	cmd.Flags().StringVar(system, "system", "", "optional system message")
}

func messageOrStdin(flag string, args []string) (string, error) {
	if strings.TrimSpace(flag) != "" {
		return flag, nil
	}
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		b, err := io.ReadAll(io.LimitReader(os.Stdin, 1<<20))
		if err != nil {
			return "", err
		}
		text := strings.TrimSpace(string(b))
		if text != "" {
			return text, nil
		}
	}
	return "", fmt.Errorf("message is required (--message, args, or stdin)")
}

func runChat(r *run, cmd *cobra.Command, leaseID, model, system, text string) error {
	if leaseID == "" {
		leaseID = r.cfg.LeaseID
	}
	var messages []api.ChatMessage
	if strings.TrimSpace(system) != "" {
		messages = append(messages, api.ChatMessage{Role: "system", Content: system})
	}
	messages = append(messages, api.ChatMessage{Role: "user", Content: text})
	raw, err := r.client.ChatCompletions(leaseID, api.ChatRequest{
		Model:    model,
		Messages: messages,
		LeaseID:  leaseID,
	})
	if err != nil {
		return err
	}
	if r.json {
		return r.printRawJSON(raw)
	}
	if out := api.AssistantText(raw); out != "" {
		fmt.Fprintln(cmd.OutOrStdout(), out)
		return nil
	}
	return r.printRawJSON(raw)
}
