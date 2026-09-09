package upgrade

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type plan struct {
	PlanVersion int       `json:"plan_version"`
	CreatedAt   time.Time `json:"created_at"`
	RequestPID  int       `json:"request_pid"`

	Options Options `json:"options"`
}

func planPath(workDir string) string {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	return filepath.Join(workDir, "plan.json")
}

func statusPath(workDir string) string {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	return filepath.Join(workDir, "status.json")
}

type status struct {
	PlanVersion int       `json:"plan_version"`
	UpdatedAt   time.Time `json:"updated_at"`
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
}

func writePlan(workDir string, opt Options) (string, error) {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", fmt.Errorf("create work dir failed: %w", err)
	}

	p := plan{
		PlanVersion: 1,
		CreatedAt:   time.Now(),
		RequestPID:  os.Getpid(),
		Options:     opt,
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal plan failed: %w", err)
	}

	pPath := planPath(workDir)
	// overwrite (single plan at a time)
	if err := os.WriteFile(pPath, b, 0600); err != nil {
		return "", fmt.Errorf("write plan failed: %w", err)
	}
	return pPath, nil
}

func readPlan(planPath string) (plan, error) {
	var p plan
	b, err := os.ReadFile(planPath)
	if err != nil {
		return p, fmt.Errorf("read plan failed: %w", err)
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return p, fmt.Errorf("unmarshal plan failed: %w", err)
	}
	if p.PlanVersion != 1 {
		return p, fmt.Errorf("unsupported plan version: %d", p.PlanVersion)
	}
	return p, nil
}

func writeStatus(workDir string, success bool, msg string) error {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return fmt.Errorf("create work dir failed: %w", err)
	}
	s := status{
		PlanVersion: 1,
		UpdatedAt:   time.Now(),
		Success:     success,
		Message:     msg,
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal status failed: %w", err)
	}
	return os.WriteFile(statusPath(workDir), b, 0644)
}

// Status 对外暴露升级结果（给 `upgrade status` 使用）
type Status struct {
	PlanVersion int       `json:"plan_version"`
	UpdatedAt   time.Time `json:"updated_at"`
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
}

func ReadStatus(workDir string) (*Status, string, error) {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	p := statusPath(workDir)
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, p, err
	}
	var s Status
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, p, err
	}
	return &s, p, nil
}
