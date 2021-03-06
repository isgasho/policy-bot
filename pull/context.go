// Copyright 2018 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pull

import (
	"time"
)

// MembershipContext defines methods to get information
// about about user membership in Github organizations and teams.
type MembershipContext interface {
	// IsTeamMember returns true if the user is a member of the given team.
	// Teams are specified as "org-name/team-name".
	IsTeamMember(team, user string) (bool, error)

	// IsOrgMember returns true if the user is a member of the given organzation.
	IsOrgMember(org, user string) (bool, error)

	// IsCollaborator returns true if the user meets the desiredPerm of the given organzation's repository.
	IsCollaborator(org, repo, user, desiredPerm string) (bool, error)
}

// Context is the context for a pull request. It defines methods to get
// information about the pull request and the VCS system containing the pull
// request (e.g. GitHub).
//
// A new Context should be created for each request, so implementations are not
// required to be thread-safe.
type Context interface {
	MembershipContext

	// Locator returns a locator string for the pull request. The locator
	// string is formated as "<owner>/<repository>#<number>"
	Locator() string

	// RepositoryOwner returns the owner of the repo that the pull request targets.
	RepositoryOwner() string

	// RepositoryName returns the repo that the pull request targets.
	RepositoryName() string

	// Author returns the username of the user who opened the pull request.
	Author() (string, error)

	// ChangedFiles returns the files that were changed in this pull request.
	ChangedFiles() ([]*File, error)

	// Commits returns the commits that are part of this pull request. The
	// commit order is implementation dependent.
	Commits() ([]*Commit, error)

	// Comments lists all comments on a Pull Request. The comment order is
	// implementation dependent.
	Comments() ([]*Comment, error)

	// Reviews lists all reviews on a Pull Request. The review order is
	// implementation dependent.
	Reviews() ([]*Review, error)

	// Branches returns the base (also known as target) and head branch names
	// of this pull request. Branches in this repository have no prefix, while
	// branches in forks are prefixed with the owner of the fork and a colon.
	// The base branch will always be unprefixed.
	Branches() (base string, head string, err error)

	// TargetCommits returns recent commits on the target branch of the pull
	// request. The exact number of commits is an implementation detail.
	TargetCommits() ([]*Commit, error)
}

type FileStatus int

const (
	FileModified FileStatus = iota
	FileAdded
	FileDeleted
)

type File struct {
	Filename  string
	Status    FileStatus
	Additions int
	Deletions int
}

type Commit struct {
	CreatedAt       time.Time
	SHA             string
	Parents         []string
	CommittedViaWeb bool

	// Author is the login name of the author. It is empty if the author is not
	// a real user.
	Author string

	// Commiter is the login name of the committer. It is empty if the
	// committer is not a real user.
	Committer string
}

// Users returns the login names of the users associated with this commit.
func (c *Commit) Users() []string {
	var users []string
	if c.Author != "" {
		users = append(users, c.Author)
	}
	if c.Committer != "" {
		users = append(users, c.Committer)
	}
	return users
}

type CommitsByCreationTime []*Commit

func (cs CommitsByCreationTime) Len() int      { return len(cs) }
func (cs CommitsByCreationTime) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs CommitsByCreationTime) Less(i, j int) bool {
	return cs[i].CreatedAt.Before(cs[j].CreatedAt)
}

type Comment struct {
	CreatedAt time.Time
	Author    string
	Body      string
}

type ReviewState string

const (
	ReviewApproved         ReviewState = "approved"
	ReviewChangesRequested ReviewState = "changes_requested"
	ReviewCommented        ReviewState = "commented"
	ReviewDismissed        ReviewState = "dismissed"
	ReviewPending          ReviewState = "pending"
)

type Review struct {
	CreatedAt time.Time
	Author    string
	State     ReviewState
	Body      string

	// ID is the GitHub node ID of the review, used to resolve dismissals
	ID string
}
