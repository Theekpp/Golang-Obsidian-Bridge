package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/user/mcp-obsidian/internal/obsidian"
)

// Register adds all Obsidian tools to the MCP server.
func Register(s *server.MCPServer, client *obsidian.Client) {
	registerListFiles(s, client)
	registerReadNote(s, client)
	registerCreateNote(s, client)
	registerUpdateNote(s, client)
	registerAppendNote(s, client)
	registerDeleteNote(s, client)
	registerSearchNotes(s, client)
	registerGetActiveNote(s, client)
	registerOpenNote(s, client)
	registerGetPeriodicNote(s, client)
	registerServerInfo(s, client)
}

// ---- list_files -----------------------------------------------------------

func registerListFiles(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"list_files",
		mcp.WithDescription("List all files and folders in the Obsidian vault. Optionally filter by a subdirectory path."),
		mcp.WithString("path",
			mcp.Description("Optional subdirectory path within the vault (e.g. 'Projects/' or 'Daily Notes'). Leave empty to list the entire vault."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")

		files, err := client.ListFiles(path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list files: %v", err)), nil
		}

		if len(files) == 0 {
			msg := "No files found"
			if path != "" {
				msg += " in " + path
			}
			return mcp.NewToolResultText(msg), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d item(s):\n\n", len(files)))
		for _, f := range files {
			sb.WriteString("- " + f + "\n")
		}

		return mcp.NewToolResultText(sb.String()), nil
	})
}

// ---- read_note ------------------------------------------------------------

func registerReadNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"read_note",
		mcp.WithDescription("Read the full content of a note from the Obsidian vault."),
		mcp.WithString("path",
			mcp.Description("Path to the note within the vault (e.g. 'Projects/MyProject.md')."),
			mcp.Required(),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("'path' is required"), nil
		}

		note, err := client.ReadNote(path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read note: %v", err)), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# %s\n\n", note.Path))
		if len(note.Tags) > 0 {
			sb.WriteString("**Tags:** " + strings.Join(note.Tags, ", ") + "\n\n")
		}
		sb.WriteString("---\n\n")
		sb.WriteString(note.Content)

		return mcp.NewToolResultText(sb.String()), nil
	})
}

// ---- create_note ----------------------------------------------------------

func registerCreateNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"create_note",
		mcp.WithDescription("Create a new note in the Obsidian vault. Fails if the note already exists (use update_note to overwrite)."),
		mcp.WithString("path",
			mcp.Description("Path for the new note (e.g. 'Projects/MyIdea.md'). Must end with .md."),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Markdown content for the new note."),
			mcp.Required(),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")
		content := req.GetString("content", "")

		if path == "" {
			return mcp.NewToolResultError("'path' is required"), nil
		}

		// Check if note already exists.
		existing, _ := client.ReadNote(path)
		if existing != nil {
			return mcp.NewToolResultError(fmt.Sprintf(
				"Note already exists at '%s'. Use update_note to overwrite it.", path,
			)), nil
		}

		if err := client.CreateOrUpdateNote(path, content); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create note: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Note created successfully at '%s'.", path)), nil
	})
}

// ---- update_note ----------------------------------------------------------

func registerUpdateNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"update_note",
		mcp.WithDescription("Overwrite an existing note with new content. Creates the note if it does not exist."),
		mcp.WithString("path",
			mcp.Description("Path to the note (e.g. 'Projects/MyProject.md')."),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("New markdown content for the note."),
			mcp.Required(),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")
		content := req.GetString("content", "")

		if path == "" {
			return mcp.NewToolResultError("'path' is required"), nil
		}

		if err := client.CreateOrUpdateNote(path, content); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update note: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Note '%s' updated successfully.", path)), nil
	})
}

// ---- append_to_note -------------------------------------------------------

func registerAppendNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"append_to_note",
		mcp.WithDescription("Append text to the end of an existing note without overwriting it."),
		mcp.WithString("path",
			mcp.Description("Path to the note (e.g. 'Daily/2024-01-15.md')."),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Text to append. A newline will be added before the content automatically."),
			mcp.Required(),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")
		content := req.GetString("content", "")

		if path == "" {
			return mcp.NewToolResultError("'path' is required"), nil
		}

		text := "\n" + content
		if err := client.AppendToNote(path, text); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to append to note: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Content appended to '%s'.", path)), nil
	})
}

// ---- delete_note ----------------------------------------------------------

func registerDeleteNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"delete_note",
		mcp.WithDescription("Permanently delete a note from the Obsidian vault."),
		mcp.WithString("path",
			mcp.Description("Path to the note to delete (e.g. 'Archive/OldNote.md')."),
			mcp.Required(),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("'path' is required"), nil
		}

		if err := client.DeleteNote(path); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete note: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Note '%s' deleted successfully.", path)), nil
	})
}

// ---- search_notes ---------------------------------------------------------

func registerSearchNotes(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"search_notes",
		mcp.WithDescription("Search the Obsidian vault for notes matching a query string. Returns matching file paths with relevance scores."),
		mcp.WithString("query",
			mcp.Description("Search query string. Supports Obsidian search syntax (e.g. 'tag:#project content')."),
			mcp.Required(),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return (default: 20)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := req.GetString("query", "")
		if query == "" {
			return mcp.NewToolResultError("'query' is required"), nil
		}

		limit := int(req.GetFloat("limit", 20))
		if limit <= 0 {
			limit = 20
		}

		results, err := client.SearchNotes(query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No notes found matching '%s'.", query)), nil
		}

		if len(results) > limit {
			results = results[:limit]
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d result(s) for '%s':\n\n", len(results), query))
		for i, r := range results {
			sb.WriteString(fmt.Sprintf("%d. %s (score: %.2f)\n", i+1, r.Filename, r.Score))
			for _, m := range r.Matches {
				if m.Context != "" {
					sb.WriteString(fmt.Sprintf("   › %s\n", strings.TrimSpace(m.Context)))
				}
			}
		}

		return mcp.NewToolResultText(sb.String()), nil
	})
}

// ---- get_active_note ------------------------------------------------------

func registerGetActiveNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"get_active_note",
		mcp.WithDescription("Get the content of the note currently open and active in Obsidian."),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		note, err := client.GetActiveNote()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get active note: %v", err)), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# Active Note: %s\n\n", note.Path))
		if len(note.Tags) > 0 {
			sb.WriteString("**Tags:** " + strings.Join(note.Tags, ", ") + "\n\n")
		}
		sb.WriteString("---\n\n")
		sb.WriteString(note.Content)

		return mcp.NewToolResultText(sb.String()), nil
	})
}

// ---- open_note ------------------------------------------------------------

func registerOpenNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"open_note",
		mcp.WithDescription("Open a specific note in the Obsidian application UI."),
		mcp.WithString("path",
			mcp.Description("Path to the note to open (e.g. 'Projects/MyProject.md')."),
			mcp.Required(),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := req.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("'path' is required"), nil
		}

		if err := client.OpenNote(path); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to open note: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Opened '%s' in Obsidian.", path)), nil
	})
}

// ---- get_periodic_note ----------------------------------------------------

func registerGetPeriodicNote(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"get_periodic_note",
		mcp.WithDescription("Get the content of a periodic note (today's daily note, this week's weekly note, etc.)."),
		mcp.WithString("period",
			mcp.Description("The type of periodic note: 'daily', 'weekly', 'monthly', or 'yearly'."),
			mcp.Required(),
			mcp.Enum("daily", "weekly", "monthly", "yearly"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		period := req.GetString("period", "")
		if period == "" {
			return mcp.NewToolResultError("'period' is required"), nil
		}

		note, err := client.GetPeriodicNote(obsidian.PeriodicNoteType(period))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get %s note: %v", period, err)), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# %s Periodic Note: %s\n\n", strings.Title(period), note.Path))
		if len(note.Tags) > 0 {
			sb.WriteString("**Tags:** " + strings.Join(note.Tags, ", ") + "\n\n")
		}
		sb.WriteString("---\n\n")
		sb.WriteString(note.Content)

		return mcp.NewToolResultText(sb.String()), nil
	})
}

// ---- server_info ----------------------------------------------------------

func registerServerInfo(s *server.MCPServer, client *obsidian.Client) {
	tool := mcp.NewTool(
		"server_info",
		mcp.WithDescription("Get information about the Obsidian Local REST API server (version, vault name, authentication status)."),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info, err := client.ServerInfo()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get server info: %v", err)), nil
		}

		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return mcp.NewToolResultError("Failed to format server info"), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	})
}
