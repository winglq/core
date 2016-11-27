package gateway

import (
	"github.com/google/go-github/github"
)

type Gateway struct {
	Client      *github.Client
	UnitTesting bool
}

func (c *Gateway) GetPullRequests(org string, project string) ([]*github.PullRequest, error) {
	pullsOpt := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	pulls := []*github.PullRequest{}
	for {
		pullRequests, resp, err := c.Client.PullRequests.List(org, project, pullsOpt)
		if err != nil {
			return nil, err
		}
		pulls = append(pulls, pullRequests...)

		if resp.NextPage == 0 || c.UnitTesting {
			break
		} else {
			pullsOpt.ListOptions.Page = resp.NextPage
		}
	}
	return pulls, nil
}

func (c *Gateway) GetIssues(org string, project string) ([]*github.Issue, error) {
	// TODO: Handle opened/closed
	issuesOpt := &github.IssueListByRepoOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	filteredIssues := []*github.Issue{}
	for {
		issues, resp, err := c.Client.Issues.ListByRepo(org, project, issuesOpt)
		for i := 0; i < len(issues); i++ {
			if err != nil {
				return nil, err
			}
			if issues[i].PullRequestLinks == nil {
				filteredIssues = append(filteredIssues, issues[i])
			}
		}
		if resp.NextPage == 0 || c.UnitTesting {
			break
		} else {
			issuesOpt.ListOptions.Page = resp.NextPage
		}
	}
	return filteredIssues, nil
}

// TODO: this may not be needed if a better mapping alternative is found
func (c *Gateway) GetPullEvents() ([]*github.PullRequestEvent, error) {
	pullEvents := []*github.PullRequestEvent{}
	return pullEvents, nil
}

// TODO: this may not be needed if a better mapping alternative is found
func (c *Gateway) GetIssueEvents() ([]*github.Event, error) {
	issuesEvents, _, _ := c.Client.Activity.ListIssueEventsForRepository("dotnet", "corefx", nil)
	return issuesEvents, nil
}
