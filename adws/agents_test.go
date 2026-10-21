package main

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Los roles de adw están en la definición (references/roles.md), sin arnés, y el ejecutor de
// Claude Code los traduce a agentes (references/executors/claude-code.md y .claude/agents/).
// Estos tests comprueban que las tres piezas dicen lo mismo.

const (
	rolesPath    = "../.claude/skills/adw/references/roles.md"
	executorPath = "../.claude/skills/adw/references/executors/claude-code.md"
	agentsDir    = "../.claude/agents"
	skillsDir    = "../.claude/skills"
)

// tableRows devuelve las filas de datos de las tablas de una sección de un markdown, celda a celda.
func tableRows(t *testing.T, path, section string) [][]string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var rows [][]string
	in := false
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "## ") {
			in = line == "## "+section
			continue
		}
		if !in || !strings.HasPrefix(line, "| ") || strings.HasPrefix(line, "| ---") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "| "), " | ")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		rows = append(rows, cells)
	}
	if len(rows) < 2 {
		t.Fatalf("%s: no hay tabla en «%s»", path, section)
	}
	return rows[1:] // sin la cabecera
}

// cellList convierte «a, b» o «`a`, `b`» en una lista ordenada; «—» es la lista vacía.
func cellList(cell string) []string {
	cell = strings.TrimSpace(strings.ReplaceAll(cell, "`", ""))
	if cell == "—" || cell == "" {
		return nil
	}
	var items []string
	for _, item := range strings.Split(cell, ",") {
		items = append(items, strings.TrimSpace(item))
	}
	sort.Strings(items)
	return items
}

type agentDef struct {
	name, model             string
	tools, disallowed, uses []string
}

func readAgent(t *testing.T, path string) agentDef {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(data), "---\n", 3)
	if len(parts) != 3 || parts[0] != "" {
		t.Fatalf("%s: sin frontmatter", path)
	}
	fields := map[string]string{}
	for _, line := range strings.Split(parts[1], "\n") {
		if key, value, ok := strings.Cut(line, ": "); ok {
			fields[key] = value
		}
	}
	if strings.TrimSpace(parts[2]) == "" || fields["description"] == "" {
		t.Fatalf("%s: le falta la descripción o el prompt", path)
	}
	return agentDef{fields["name"], fields["model"], cellList(fields["tools"]), cellList(fields["disallowedTools"]), cellList(fields["skills"])}
}

func TestAgentsMatchTheRoles(t *testing.T) {
	levels := map[string]string{}
	for _, row := range tableRows(t, executorPath, "Niveles") {
		levels[row[0]] = strings.Trim(row[1], "`")
	}
	roleLevel := map[string]string{}
	for _, row := range tableRows(t, rolesPath, "Catálogo") {
		roleLevel[strings.Trim(row[0], "`")] = row[2]
	}
	executor := map[string][]string{}
	for _, row := range tableRows(t, executorPath, "Agentes") {
		executor[strings.Trim(row[0], "`")] = row[1:]
	}

	files, err := filepath.Glob(filepath.Join(agentsDir, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	agents := map[string]agentDef{}
	for _, file := range files {
		a := readAgent(t, file)
		if a.name != strings.TrimSuffix(filepath.Base(file), ".md") {
			t.Errorf("%s: name %q no es el del fichero", file, a.name)
		}
		agents[a.name] = a
	}

	for role, level := range roleLevel {
		if role == "coordinator" {
			if _, ok := agents[role]; ok {
				t.Errorf("coordinator es la sesión que lanza el trigger: no tiene agente")
			}
			cfg, err := readConfig("../config/adw.yml")
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Model != levels[level] {
				t.Errorf("config/adw.yml: model %q, y el nivel del coordinador (%s) es %q", cfg.Model, level, levels[level])
			}
			continue
		}
		a, ok := agents[role]
		if !ok {
			t.Errorf("el rol %s no tiene agente en .claude/agents/", role)
			continue
		}
		model, ok := levels[level]
		if !ok {
			t.Errorf("el rol %s tiene un nivel desconocido: %q", role, level)
		}
		if a.model != model {
			t.Errorf("%s: model %q, y su nivel (%s) es %q", role, a.model, level, model)
		}
		row, ok := executor[role]
		if !ok {
			t.Errorf("el rol %s no está en la tabla de agentes del ejecutor", role)
			continue
		}
		if !reflect.DeepEqual(a.tools, cellList(row[0])) || !reflect.DeepEqual(a.disallowed, cellList(row[1])) || !reflect.DeepEqual(a.uses, cellList(row[2])) {
			t.Errorf("%s: tools %v, disallowedTools %v y skills %v no son los del ejecutor: %v", role, a.tools, a.disallowed, a.uses, row)
		}
		for _, skill := range a.uses {
			if _, err := os.Stat(filepath.Join(skillsDir, skill, "SKILL.md")); err != nil {
				t.Errorf("%s: la skill %s no existe", role, skill)
			}
		}
		for _, tool := range a.tools {
			if tool == "Agent" {
				t.Errorf("%s: un rol no lanza otros, y no tiene la herramienta Agent", role)
			}
		}
	}
	for name := range agents {
		if _, ok := roleLevel[name]; !ok {
			t.Errorf("el agente %s no es ningún rol del catálogo", name)
		}
	}
}

func TestRolesThatDoNotChangeCodeCannotEdit(t *testing.T) {
	data, err := os.ReadFile(rolesPath)
	if err != nil {
		t.Fatal(err)
	}
	var keep []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "- **Quién no cambia el código.**") {
			for _, m := range regexp.MustCompile("`([a-z-]+)`").FindAllStringSubmatch(line, -1) {
				keep = append(keep, m[1])
			}
		}
	}
	if len(keep) == 0 {
		t.Fatal("roles.md no dice qué roles no cambian el código")
	}
	readOnly := map[string]bool{}
	for _, role := range keep {
		readOnly[role] = true
	}
	files, err := filepath.Glob(filepath.Join(agentsDir, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		a := readAgent(t, file)
		has := func(tools []string) bool {
			for _, tool := range tools {
				if tool == "Edit" {
					return true
				}
			}
			return false
		}
		if readOnly[a.name] && (has(a.tools) || !has(a.disallowed)) {
			t.Errorf("%s no cambia el código: no lleva Edit en tools y lo lleva en disallowedTools", a.name)
		}
		if !readOnly[a.name] && !has(a.tools) {
			t.Errorf("%s cambia el código: lleva Edit", a.name)
		}
	}
}
