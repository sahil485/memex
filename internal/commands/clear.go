package commands

import (
	"fmt"

	"github.com/sahil485/memex/pkg/client"
)

func ClearIndex(args []string) error {
	c := client.New()
	idx := c.GetIndex()

	fmt.Println("Clearing all documents from index...")

	// Delete all documents from the index
	task, err := idx.DeleteAllDocuments(nil)
	if err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}

	// Wait for the deletion to complete
	taskInfo, err := c.WaitForTask(task.TaskUID)
	if err != nil {
		return fmt.Errorf("failed to wait for deletion task: %w", err)
	}

	// Check if the task failed
	if taskInfo.Status == "failed" {
		return fmt.Errorf("deletion task failed: %s", taskInfo.Error.Message)
	}

	fmt.Println("✓ Index cleared successfully")
	return nil
}
