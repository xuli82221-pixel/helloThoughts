package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"

	"thoughts.local/thoughts/internal/thoughts"
)

func output(v any) { _ = json.NewEncoder(os.Stdout).Encode(v) }
func run(args []string) error {
	command := "open"
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}
	if command == "version" || command == "--version" {
		output(map[string]string{"version": thoughts.Version})
		return nil
	}
	if command == "capabilities" {
		output(map[string]any{"product": "thoughts", "version": thoughts.Version, "protocolVersion": 1, "managedRuntime": 1, "readerState": 1, "sharedDefaultHome": 1})
		return nil
	}
	if command == "help" || command == "--help" {
		fmt.Println("Thoughts: init | open | serve | doctor | paths | library | projects list/create | notes list/read/create/update/rename/move/trash | trash list/restore\nRuntime: runtime install|upgrade --source REVIEWED_EXE --sha256 SHA256 --json; runtime status|rollback --json\nDefault home: Windows account profile/Thoughts; THOUGHTS_HOME overrides it. Runtime changes keep library files.\nFlags: --vault DIR --json; open --browser none|system\nWrite: --path relative.md --content-file UTF8_FILE (or stdin) --revision SHA256\nRename: --path relative.md --name NEW_NAME --revision SHA256\nMove: --path relative.md --to existing/folder --revision SHA256\nTrash: notes trash --path relative.md --revision SHA256; trash restore --id ID --revision SHA256\nProject: projects create --name NAME --language en|ko|ja|zh-Hans|zh-Hant --template=true")
		return nil
	}
	noteCommand := ""
	if command == "notes" || command == "projects" || command == "trash" || command == "runtime" {
		if len(args) == 0 {
			return fmt.Errorf("请指定 notes create/read/update/list")
		}
		noteCommand = args[0]
		args = args[1:]
	}
	f := flag.NewFlagSet(command, flag.ContinueOnError)
	f.SetOutput(io.Discard)
	vault := f.String("vault", "", "资料库目录")
	_ = f.Bool("json", false, "JSON 输出")
	browser := f.String("browser", "system", "none 或 system")
	p := f.String("path", "", "相对 Markdown 路径")
	contentFile := f.String("content-file", "", "UTF-8 文件，不填写则读取 stdin")
	revision := f.String("revision", "", "读取时的 SHA256")
	destination := f.String("to", "", "Destination folder for notes move")
	name := f.String("name", "", "New name without .md")
	id := f.String("id", "", "Trash entry ID")
	language := f.String("language", "en", "Template language")
	source := f.String("source", "", "Reviewed Windows executable for a managed runtime install/update")
	expected := f.String("sha256", "", "Reviewed executable SHA256")
	template := f.Bool("template", true, "Create the seven editable starter nodes")
	if e := f.Parse(args); e != nil {
		return fmt.Errorf("参数无效，请运行 thoughts help")
	}
	if f.NArg() != 0 {
		return fmt.Errorf("存在未识别的位置参数")
	}
	if command == "mcp" {
		// stdout is reserved for JSON-RPC after entering MCP mode.
		return thoughts.RunMCP(context.Background(), *vault)
	}
	if command == "runtime" {
		value, e := thoughts.ManageRuntime(context.Background(), noteCommand, *source, *expected)
		if e != nil {
			return e
		}
		output(value)
		return nil
	}
	if command == "init" {
		v, e := thoughts.Init(*vault)
		if e != nil {
			return e
		}
		output(map[string]any{"vault": v, "version": thoughts.Version, "initialized": true})
		return nil
	}
	if command == "paths" {
		value, e := thoughts.Paths()
		if e != nil {
			return e
		}
		output(value)
		return nil
	}
	loc, e := thoughts.Resolve(*vault)
	if e != nil {
		return e
	}
	ctx := context.Background()
	switch command {
	case "serve":
		return thoughts.Serve(loc)
	case "doctor":
		output(thoughts.Doctor(ctx, loc))
		return nil
	case "open":
		if *browser != "none" && *browser != "system" {
			return fmt.Errorf("browser 必须为 none 或 system")
		}
		v, e := thoughts.Open(ctx, loc, *browser)
		if e != nil {
			return e
		}
		output(v)
		return nil
	case "library":
		var out any
		if e = thoughts.Call(ctx, loc, "GET", "/api/library", nil, &out); e != nil {
			return e
		}
		output(out)
		return nil
	case "projects":
		var out any
		switch noteCommand {
		case "list":
			e = thoughts.Call(ctx, loc, "GET", "/api/library", nil, &out)
		case "create":
			if *name == "" {
				return fmt.Errorf("--name is required")
			}
			e = thoughts.Call(ctx, loc, "POST", "/api/projects", map[string]any{"name": *name, "language": *language, "template": *template}, &out)
		default:
			return fmt.Errorf("Use projects list or projects create")
		}
		if e != nil {
			return e
		}
		output(out)
		return nil
	case "trash":
		var out any
		switch noteCommand {
		case "list":
			e = thoughts.Call(ctx, loc, "GET", "/api/trash", nil, &out)
		case "restore":
			if *id == "" || *revision == "" {
				return fmt.Errorf("--id and --revision are required")
			}
			e = thoughts.Call(ctx, loc, "POST", "/api/trash/restore", map[string]string{"id": *id, "revision": *revision}, &out)
		default:
			return fmt.Errorf("Use trash list or trash restore")
		}
		if e != nil {
			return e
		}
		output(out)
		return nil
	case "notes":
		var out any
		switch noteCommand {
		case "read":
			if *p == "" {
				return fmt.Errorf("需要 --path")
			}
			e = thoughts.Call(ctx, loc, "GET", "/api/notes?path="+url.QueryEscape(*p), nil, &out)
		case "list":
			e = thoughts.Call(ctx, loc, "GET", "/api/notes", nil, &out)
		case "move":
			if *p == "" || *revision == "" || *destination == "" {
				return fmt.Errorf("--path, --to and --revision are required")
			}
			e = thoughts.Call(ctx, loc, "POST", "/api/notes/move", map[string]string{"path": *p, "directory": *destination, "revision": *revision}, &out)
		case "rename", "trash", "delete":
			if *p == "" || *revision == "" {
				return fmt.Errorf("--path and --revision are required; read the current note first")
			}
			if noteCommand == "rename" {
				if *name == "" {
					return fmt.Errorf("--name is required")
				}
				e = thoughts.Call(ctx, loc, "POST", "/api/notes/rename", map[string]string{"path": *p, "name": *name, "revision": *revision}, &out)
			} else {
				e = thoughts.Call(ctx, loc, "POST", "/api/notes/trash", map[string]string{"path": *p, "revision": *revision}, &out)
			}
		case "create", "update":
			if *p == "" {
				return fmt.Errorf("需要 --path")
			}
			if noteCommand == "update" && *revision == "" {
				return fmt.Errorf("更新必须提供 --revision，先读取文件获取版本")
			}
			var reader io.Reader = os.Stdin
			var file *os.File
			if *contentFile != "" {
				file, e = os.Open(*contentFile)
				if e != nil {
					return fmt.Errorf("无法读取正文文件")
				}
				defer file.Close()
				reader = file
			}
			b, readErr := io.ReadAll(io.LimitReader(reader, thoughts.MaxNoteSize+1))
			if readErr != nil {
				return readErr
			}
			if len(b) > thoughts.MaxNoteSize {
				return fmt.Errorf("正文不能超过 1 MiB")
			}
			method := "POST"
			if noteCommand == "update" {
				method = "PUT"
			}
			e = thoughts.Call(ctx, loc, method, "/api/notes", map[string]string{"path": *p, "content": string(b), "revision": *revision}, &out)
		default:
			return fmt.Errorf("Use notes list/read/create/update/rename/move/trash")
		}
		if e != nil {
			return e
		}
		output(out)
		return nil
	default:
		return fmt.Errorf("未知命令，请运行 thoughts help")
	}
}
func main() {
	forwarded, forwardErr := thoughts.ForwardManaged(os.Args[1:])
	if forwarded {
		if forwardErr != nil {
			var exit *exec.ExitError
			if errors.As(forwardErr, &exit) {
				os.Exit(exit.ExitCode())
			}
			fmt.Fprintln(os.Stderr, "Thoughts runtime launch failed:", forwardErr)
			os.Exit(1)
		}
		return
	}
	if forwardErr != nil {
		output(map[string]any{"error": map[string]string{"code": "runtime_invalid", "message": forwardErr.Error()}})
		os.Exit(1)
	}
	if e := run(os.Args[1:]); e != nil {
		if len(os.Args) > 1 && os.Args[1] == "mcp" {
			fmt.Fprintln(os.Stderr, "Thoughts MCP stopped:", e)
			os.Exit(1)
		}
		var p *thoughts.Error
		body := map[string]any{"code": "command_failed", "message": e.Error()}
		if errors.As(e, &p) {
			body["code"] = p.Code
			body["status"] = p.Status
		}
		output(map[string]any{"error": body})
		os.Exit(1)
	}
}
