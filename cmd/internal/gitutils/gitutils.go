package gitutils

import (
	"bytes"
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GetBaseBranch returns the base branch of the current git repository
func GetBaseBranch() (string, error) {
	// Open the current Git repository located in the current directory.
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("error opening git repository: %w", err)
	}

	// Retrieve the list of branches in the repository.
	branches, err := r.Branches()
	if err != nil {
		return "", fmt.Errorf("error getting branches: %w", err)
	}

	var baseBranch string
	// Iterate over the branches to find "master" or "main".
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		if name == "master" || name == "main" {
			// Set baseBranch if "master" or "main" is found.
			baseBranch = name
			return nil
		}
		return nil
	})

	if err != nil {
		// Return an error if there was an issue iterating through branches.
		return "", fmt.Errorf("error iterating branches: %w", err)
	}

	// If neither "master" nor "main" was found, return an error.
	if baseBranch == "" {
		return "", fmt.Errorf("neither 'master' nor 'main' branches exist")
	}

	// Return the base branch name.
	return baseBranch, nil
}

// GetGitDiff gets the git diff for the specified base branch
func GetGitDiff(baseBranch string) (string, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("error opening git repository: %w", err)
	}

	// Get the reference to the base branch
	baseRef, err := r.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
	if err != nil {
		return "", fmt.Errorf("error getting base branch reference: %w", err)
	}

	// Get the commit object for the base branch
	baseCommit, err := r.CommitObject(baseRef.Hash())
	if err != nil {
		return "", fmt.Errorf("error getting base commit: %w", err)
	}

	// Get the tree object for the base commit
	baseTree, err := baseCommit.Tree()
	if err != nil {
		return "", fmt.Errorf("error getting base tree: %w", err)
	}

	// Get the HEAD reference
	headRef, err := r.Head()
	if err != nil {
		return "", fmt.Errorf("error getting head reference: %w", err)
	}

	// Get the commit object for HEAD
	headCommit, err := r.CommitObject(headRef.Hash())
	if err != nil {
		return "", fmt.Errorf("error getting head commit: %w", err)
	}

	// Get the tree object for the HEAD commit
	headTree, err := headCommit.Tree()
	if err != nil {
		return "", fmt.Errorf("error getting head tree: %w", err)
	}

	// Generate the diff between the two trees
	changes, err := baseTree.Diff(headTree)
	if err != nil {
		return "", fmt.Errorf("error generating patch: %w", err)
	}

	var buf bytes.Buffer
	for _, change := range changes {
		patch, err := change.Patch()
		if err != nil {
			return "", fmt.Errorf("error creating patch: %w", err)
		}
		if err := patch.Encode(&buf); err != nil {
			return "", fmt.Errorf("error encoding patch: %w", err)
		}
	}

	return buf.String(), nil
}
