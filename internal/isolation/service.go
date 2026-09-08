package isolation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/environment"
	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
	"ai-dev-manager-v2/internal/workspace"
)

type Service struct {
	store        *store.Store
	workspaces   *workspace.Service
	environments *environment.Service
	now          func() time.Time
}

type CreateResult struct {
	ManagedWorktree model.ManagedWorktree `json:"managed_worktree"`
	Environment     model.Environment     `json:"environment"`
}

type DestroyResult struct {
	ManagedWorktreeID string `json:"managed_worktree_id"`
	EnvironmentID     string `json:"environment_id"`
	Root              string `json:"root"`
	RetainedBranch    string `json:"retained_branch"`
	Forced            bool   `json:"forced"`
}

type SafetyStatus struct {
	Dirty       bool   `json:"dirty"`
	Unpublished bool   `json:"unpublished"`
	Head        string `json:"head"`
}

func New(s *store.Store, workspaces *workspace.Service, environments *environment.Service) *Service {
	return &Service{store: s, workspaces: workspaces, environments: environments, now: time.Now}
}

func (s *Service) List() ([]model.ManagedWorktree, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	items := append([]model.ManagedWorktree(nil), state.ManagedWorktrees...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items, nil
}

func (s *Service) GetByEnvironment(environmentID string) (model.ManagedWorktree, bool, error) {
	environmentID = strings.TrimSpace(environmentID)
	state, err := s.store.Load()
	if err != nil {
		return model.ManagedWorktree{}, false, err
	}
	for _, item := range state.ManagedWorktrees {
		if item.EnvironmentID == environmentID {
			return item, true, nil
		}
	}
	return model.ManagedWorktree{}, false, nil
}

func (s *Service) Create(ctx context.Context, workspaceID, name, baseRef string) (CreateResult, error) {
	ws, err := s.workspaces.Get(strings.TrimSpace(workspaceID))
	if err != nil {
		return CreateResult{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return CreateResult{}, fmt.Errorf("environment name is required")
	}
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		baseRef = "HEAD"
	}
	if strings.HasPrefix(baseRef, "-") || strings.ContainsRune(baseRef, '\x00') {
		return CreateResult{}, fmt.Errorf("invalid base_ref %q", baseRef)
	}
	if _, err := exec.LookPath("git"); err != nil {
		return CreateResult{}, fmt.Errorf("git worktree isolation is unavailable: %w", err)
	}

	top, err := s.gitOutput(ctx, ws.Path, "rev-parse", "--show-toplevel")
	if err != nil {
		return CreateResult{}, fmt.Errorf("workspace %s is not a usable Git worktree: %w", ws.ID, err)
	}
	topDir, err := canonicalExistingDir(strings.TrimSpace(top))
	if err != nil {
		return CreateResult{}, fmt.Errorf("resolve Git top-level: %w", err)
	}
	if !samePath(topDir, ws.Path) {
		return CreateResult{}, fmt.Errorf("managed worktree creation requires workspace root %s to be the Git top-level %s", ws.Path, topDir)
	}
	commonRaw, err := s.gitOutput(ctx, ws.Path, "rev-parse", "--git-common-dir")
	if err != nil {
		return CreateResult{}, err
	}
	commonDir, err := canonicalGitPath(ws.Path, strings.TrimSpace(commonRaw))
	if err != nil {
		return CreateResult{}, fmt.Errorf("resolve Git common directory: %w", err)
	}
	baseCommit, err := s.gitOutput(ctx, ws.Path, "rev-parse", "--verify", "--end-of-options", baseRef+"^{commit}")
	if err != nil {
		return CreateResult{}, fmt.Errorf("resolve base_ref %q: %w", baseRef, err)
	}
	baseCommit = strings.TrimSpace(baseCommit)

	managedID, err := identity.New("wt")
	if err != nil {
		return CreateResult{}, err
	}
	branch := "adm/" + managedID
	root := filepath.Join(s.ownedRoot(), ws.ID, managedID)
	if err := os.MkdirAll(filepath.Dir(root), 0o755); err != nil {
		return CreateResult{}, err
	}
	if _, err := os.Stat(root); err == nil {
		return CreateResult{}, fmt.Errorf("managed worktree destination already exists: %s", root)
	} else if !os.IsNotExist(err) {
		return CreateResult{}, err
	}

	if _, err := s.gitOutput(ctx, ws.Path, "worktree", "add", "-b", branch, root, baseCommit); err != nil {
		return CreateResult{}, fmt.Errorf("create managed worktree: %w", err)
	}
	rollback := true
	defer func() {
		if !rollback {
			return
		}
		_, _ = s.gitOutput(context.Background(), ws.Path, "worktree", "remove", "--force", root)
		_, _ = s.gitOutput(context.Background(), ws.Path, "branch", "-D", branch)
		_ = os.RemoveAll(root)
	}()

	managed := model.ManagedWorktree{
		ID:           managedID,
		Branch:       branch,
		BaseCommit:   baseCommit,
		GitCommonDir: commonDir,
		CreatedAt:    s.nowUTC(),
	}
	env, persisted, err := s.environments.CreateManaged(ws.ID, name, root, managed)
	if err != nil {
		return CreateResult{}, err
	}
	rollback = false
	return CreateResult{ManagedWorktree: persisted, Environment: env}, nil
}

func (s *Service) ValidateEnvironment(ctx context.Context, env model.Environment) error {
	ws, err := s.workspaces.Get(env.WorkspaceID)
	if err != nil {
		return err
	}
	managed, ok, err := s.GetByEnvironment(env.ID)
	if err != nil {
		return err
	}
	if !ok {
		if within(ws.Path, env.Root) {
			return nil
		}
		return fmt.Errorf("environment %s root %s is outside workspace %s without managed worktree metadata", env.ID, env.Root, ws.Path)
	}
	if managed.EnvironmentID != env.ID || managed.WorkspaceID != env.WorkspaceID || !samePath(managed.Root, env.Root) {
		return fmt.Errorf("managed worktree metadata does not match environment %s", env.ID)
	}
	owned := s.ownedRoot()
	if !within(owned, managed.Root) {
		return fmt.Errorf("managed worktree %s root escapes ADM-owned root", managed.ID)
	}
	root, err := canonicalExistingDir(managed.Root)
	if err != nil {
		return fmt.Errorf("managed worktree %s root is missing or invalid: %w", managed.ID, err)
	}
	if !samePath(root, managed.Root) {
		return fmt.Errorf("managed worktree %s root identity changed", managed.ID)
	}
	top, err := s.gitOutput(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("managed worktree %s Git top-level check failed: %w", managed.ID, err)
	}
	topDir, err := canonicalExistingDir(strings.TrimSpace(top))
	if err != nil || !samePath(topDir, root) {
		return fmt.Errorf("managed worktree %s Git top-level no longer matches root", managed.ID)
	}
	commonRaw, err := s.gitOutput(ctx, root, "rev-parse", "--git-common-dir")
	if err != nil {
		return fmt.Errorf("managed worktree %s Git common-dir check failed: %w", managed.ID, err)
	}
	commonDir, err := canonicalGitPath(root, strings.TrimSpace(commonRaw))
	if err != nil || !samePath(commonDir, managed.GitCommonDir) {
		return fmt.Errorf("managed worktree %s Git common directory changed", managed.ID)
	}
	branch, err := s.gitOutput(ctx, root, "branch", "--show-current")
	if err != nil {
		return fmt.Errorf("managed worktree %s branch check failed: %w", managed.ID, err)
	}
	if strings.TrimSpace(branch) != managed.Branch {
		return fmt.Errorf("managed worktree %s branch changed: got %q want %q", managed.ID, strings.TrimSpace(branch), managed.Branch)
	}
	return nil
}

func (s *Service) Safety(ctx context.Context, environmentID string) (SafetyStatus, error) {
	env, err := s.environments.Get(strings.TrimSpace(environmentID))
	if err != nil {
		return SafetyStatus{}, err
	}
	managed, ok, err := s.GetByEnvironment(env.ID)
	if err != nil {
		return SafetyStatus{}, err
	}
	if !ok {
		return SafetyStatus{}, fmt.Errorf("environment %s is not backed by a managed worktree", env.ID)
	}
	if err := s.ValidateEnvironment(ctx, env); err != nil {
		return SafetyStatus{}, err
	}
	status, err := s.gitOutput(ctx, managed.Root, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return SafetyStatus{}, err
	}
	head, err := s.gitOutput(ctx, managed.Root, "rev-parse", "HEAD")
	if err != nil {
		return SafetyStatus{}, err
	}
	head = strings.TrimSpace(head)
	unpublished := false
	if head != managed.BaseCommit {
		remoteRefs, err := s.gitOutput(ctx, managed.Root, "branch", "-r", "--contains", head, "--format=%(refname:short)")
		if err != nil {
			return SafetyStatus{}, err
		}
		unpublished = strings.TrimSpace(remoteRefs) == ""
	}
	return SafetyStatus{Dirty: strings.TrimSpace(status) != "", Unpublished: unpublished, Head: head}, nil
}

func (s *Service) Destroy(ctx context.Context, environmentID, writerOwner string, force bool) (DestroyResult, error) {
	env, err := s.environments.RequireWriter(strings.TrimSpace(environmentID), strings.TrimSpace(writerOwner))
	if err != nil {
		return DestroyResult{}, err
	}
	managed, ok, err := s.GetByEnvironment(env.ID)
	if err != nil {
		return DestroyResult{}, err
	}
	if !ok {
		return DestroyResult{}, fmt.Errorf("environment %s is not backed by a managed worktree", env.ID)
	}
	if err := s.ValidateEnvironment(ctx, env); err != nil {
		return DestroyResult{}, err
	}
	safety, err := s.Safety(ctx, env.ID)
	if err != nil {
		return DestroyResult{}, err
	}
	if !force && (safety.Dirty || safety.Unpublished) {
		reasons := make([]string, 0, 2)
		if safety.Dirty {
			reasons = append(reasons, "dirty worktree")
		}
		if safety.Unpublished {
			reasons = append(reasons, "locally advanced HEAD is not present in any remote-tracking ref")
		}
		return DestroyResult{}, fmt.Errorf("managed worktree %s destroy refused: %s; retry with force=true only after review", managed.ID, strings.Join(reasons, ", "))
	}
	ws, err := s.workspaces.Get(managed.WorkspaceID)
	if err != nil {
		return DestroyResult{}, err
	}
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, managed.Root)
	if _, err := s.gitOutput(ctx, ws.Path, args...); err != nil {
		return DestroyResult{}, fmt.Errorf("remove managed worktree %s: %w", managed.ID, err)
	}
	if _, _, err := s.environments.RemoveManaged(env.ID); err != nil {
		return DestroyResult{}, fmt.Errorf("managed worktree removed from Git but ADM metadata cleanup failed: %w", err)
	}
	_ = os.Remove(filepath.Dir(managed.Root))
	return DestroyResult{
		ManagedWorktreeID: managed.ID,
		EnvironmentID:     env.ID,
		Root:              managed.Root,
		RetainedBranch:    managed.Branch,
		Forced:            force,
	}, nil
}

func (s *Service) ownedRoot() string {
	root := filepath.Join(filepath.Dir(s.store.Path()), "worktrees")
	abs, err := filepath.Abs(root)
	if err == nil {
		root = abs
	}
	return filepath.Clean(root)
}

func (s *Service) gitOutput(ctx context.Context, dir string, args ...string) (string, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}
	fullArgs := append([]string{"-C", dir}, args...)
	cmd := exec.CommandContext(ctx, git, fullArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (s *Service) nowUTC() time.Time {
	if s.now == nil {
		return time.Now().UTC()
	}
	return s.now().UTC()
}

func canonicalExistingDir(path string) (string, error) {
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("not a directory: %s", abs)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func canonicalGitPath(base, value string) (string, error) {
	if !filepath.IsAbs(value) {
		value = filepath.Join(base, value)
	}
	return canonicalExistingDir(value)
}

func within(base, target string) bool {
	rel, err := filepath.Rel(filepath.Clean(base), filepath.Clean(target))
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func samePath(a, b string) bool { return strings.EqualFold(filepath.Clean(a), filepath.Clean(b)) }
