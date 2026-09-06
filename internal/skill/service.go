package skill

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ai-dev-manager-v2/internal/model"
)

const defaultMaxBytes = 1 << 20

type Content struct {
	SkillID string `json:"skill_id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

func Discover(root string, supportRoots []string, defaultInclude bool) ([]model.CatalogEntry, error) {
	resolvedRoot, err := canonicalDirectory(root)
	if err != nil {
		return nil, fmt.Errorf("skill discovery root: %w", err)
	}
	resolvedSupportRoots := make([]string, 0, len(supportRoots))
	seenSupportRoots := map[string]struct{}{}
	for _, supportRoot := range supportRoots {
		if strings.TrimSpace(supportRoot) == "" {
			continue
		}
		resolved, err := canonicalDirectory(supportRoot)
		if err != nil {
			return nil, fmt.Errorf("skill support root: %w", err)
		}
		key := strings.ToLower(filepath.Clean(resolved))
		if _, exists := seenSupportRoots[key]; exists {
			continue
		}
		seenSupportRoots[key] = struct{}{}
		resolvedSupportRoots = append(resolvedSupportRoots, resolved)
	}
	sort.Slice(resolvedSupportRoots, func(i, j int) bool {
		return strings.ToLower(resolvedSupportRoots[i]) < strings.ToLower(resolvedSupportRoots[j])
	})

	byID := map[string]model.CatalogEntry{}
	err = filepath.WalkDir(resolvedRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() || entry.Name() != "SKILL.md" {
			return nil
		}
		dir := filepath.Dir(path)
		name := strings.TrimSpace(filepath.Base(dir))
		if name == "" || name == "." || name == string(filepath.Separator) {
			return fmt.Errorf("invalid skill directory for %s", path)
		}
		id := StableID(name)
		artifact, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
		artifact = filepath.Clean(artifact)
		if existing, exists := byID[id]; exists && !samePath(existing.ArtifactPath, artifact) {
			return fmt.Errorf("duplicate skill name %q under discovery root %s", name, resolvedRoot)
		}
		byID[id] = model.CatalogEntry{
			ID:                  id,
			Name:                name,
			DefaultIncludeInEnv: defaultInclude,
			ArtifactPath:        artifact,
			SourceRoot:          resolvedRoot,
			SupportRoots:        append([]string(nil), resolvedSupportRoots...),
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover skills under %s: %w", resolvedRoot, err)
	}
	if len(byID) == 0 {
		return nil, fmt.Errorf("skill discovery root %s contains no SKILL.md artifacts", resolvedRoot)
	}
	out := make([]model.CatalogEntry, 0, len(byID))
	for _, entry := range byID {
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func StableID(name string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(name))))
	return fmt.Sprintf("skill_%x", sum[:16])
}

func Configured(entry model.CatalogEntry) bool {
	return strings.TrimSpace(entry.ArtifactPath) != "" && strings.TrimSpace(entry.SourceRoot) != ""
}

func Read(entry model.CatalogEntry, requestedPath string, maxBytes int) (Content, error) {
	if !Configured(entry) {
		return Content{}, fmt.Errorf("skill %q is not backed by a configured artifact", entry.ID)
	}
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}
	artifact, err := canonicalFile(entry.ArtifactPath)
	if err != nil {
		return Content{}, fmt.Errorf("skill %q artifact: %w", entry.ID, err)
	}
	artifactDir := filepath.Dir(artifact)

	target := artifact
	requestedPath = strings.TrimSpace(requestedPath)
	if requestedPath != "" && requestedPath != "SKILL.md" {
		if filepath.IsAbs(requestedPath) {
			target = filepath.Clean(requestedPath)
		} else {
			target = filepath.Clean(filepath.Join(artifactDir, requestedPath))
		}
	}
	resolvedTarget, err := canonicalFile(target)
	if err != nil {
		return Content{}, fmt.Errorf("skill %q read %q: %w", entry.ID, requestedPath, err)
	}

	allowedRoots := []string{artifactDir}
	for _, root := range entry.SupportRoots {
		resolved, err := canonicalDirectory(root)
		if err != nil {
			return Content{}, fmt.Errorf("skill %q support root: %w", entry.ID, err)
		}
		allowedRoots = append(allowedRoots, resolved)
	}
	allowed := false
	for _, root := range allowedRoots {
		if within(root, resolvedTarget) {
			allowed = true
			break
		}
	}
	if !allowed {
		return Content{}, fmt.Errorf("skill %q path is outside configured artifact/support roots", entry.ID)
	}

	info, err := os.Stat(resolvedTarget)
	if err != nil {
		return Content{}, err
	}
	if !info.Mode().IsRegular() {
		return Content{}, fmt.Errorf("skill path is not a regular file")
	}
	if info.Size() > int64(maxBytes) {
		return Content{}, fmt.Errorf("skill file exceeds max_bytes=%d", maxBytes)
	}
	data, err := os.ReadFile(resolvedTarget)
	if err != nil {
		return Content{}, err
	}
	if containsNUL(data) {
		return Content{}, fmt.Errorf("binary skill file is not readable as text")
	}
	return Content{SkillID: entry.ID, Name: entry.Name, Path: resolvedTarget, Content: string(data)}, nil
}

func canonicalDirectory(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return "", fmt.Errorf("path is not a real directory: %s", abs)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func canonicalFile(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func within(root, target string) bool {
	rel, err := filepath.Rel(filepath.Clean(root), filepath.Clean(target))
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func samePath(a, b string) bool { return strings.EqualFold(filepath.Clean(a), filepath.Clean(b)) }

func containsNUL(data []byte) bool {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		if strings.IndexByte(scanner.Text(), 0) >= 0 {
			return true
		}
	}
	return strings.IndexByte(string(data), 0) >= 0
}
