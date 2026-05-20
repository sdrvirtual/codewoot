package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func removeInstanceFiles(instancesDir, instance string) error {
	if instancesDir == "" {
		return nil
	}

	parsed, err := uuid.Parse(instance)
	if err != nil {
		return fmt.Errorf("instance must be a uuid: %w", err)
	}
	instance = parsed.String()

	base, err := filepath.Abs(instancesDir)
	if err != nil {
		return err
	}
	if base == string(os.PathSeparator) {
		return fmt.Errorf("instances dir cannot be filesystem root")
	}
	info, err := os.Stat(base)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("instances dir is not a directory")
	}

	target, err := filepath.Abs(filepath.Join(base, instance))
	if err != nil {
		return err
	}

	if target != filepath.Join(base, instance) || !strings.HasPrefix(target, base+string(os.PathSeparator)) {
		return fmt.Errorf("invalid instance path")
	}

	return os.RemoveAll(target)
}
