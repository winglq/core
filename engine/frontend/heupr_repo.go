package frontend

import (
	// "coralreefci/engine/gateway/conflation"
	"coralreefci/models"

	"github.com/google/go-github/github"
)

type HeuprModel struct {
	Model *models.Model
	// Benchmark        Benchmark // TODO: Struct to build
	// Scenarios        []conflation.Scenario
	// PilotScenarios   []conflation.Scenario
	// LearnedScenarios []conflation.Scenario
	// StrategyParams   StrategyParams TODO: Baseline with self evolving
	//                                       parameters (Tossing Graph?,
	//                                       Conflation Scenarios, etc.)
}

type HeuprHive struct {
	Models []*HeuprModel
	// PilotModels        []HeuprModel
	// ModelBlender       ModelBlender // TODO: Struct to build
	// TossingGraph       TossingGraphAlgorithm // TODO: Struct to build
	// StrategyParams     StrategyParams // TODO: Struct to build
	// AggregateBenchmark Benchmark
}

type HeuprRepo struct {
	Repo   *github.Repository
	Hive   *HeuprHive
	Client *github.Client
}

func (h *HeuprServer) NewHeuprRepo(repos []*github.Repository, client *github.Client) {
	for _, repo := range repos {
		h.Repos[*repo.ID] = &HeuprRepo{
			Repo:   repo,
			Client: client,
		}
	}
}

// TODO: This method will place the desired model(s) for a given repository ID
//       into the respective field e.g.:
//       (h.Hive.Models = append(h.Hive.Models, exampleModel))
//       Additionally, this method will need to be able to properly call in the
//       necessary logic for defining how each model(s) is trained (although it
//       is quite likely that they will all train on all available information
//       which will be filtered down by the conflation scenarios).
func (h *HeuprServer) initModels(id int) {

}

/*
// TODO: This method will need to change substantially in the switch to gob.
func (h *HeuprServer) InitHeuprRepos(path ...string) {
	defer h.CloseDB()
	h.OpenDB()
	if path == nil {
		// tokens, err := h.Database.retrieveBulk()
		// if err != nil {
        //
		// }
	} else {

	}

	// TODO:
	// if string is nil
	// - default to opening the storage.db file
	// - load in info from file
	// - boot up models / repos via that information
	// else
	// - open the provided file
	// - boot up from those specifications
}
*/

func (h *HeuprServer) GetIssues(id int) ([]*github.Issue, error) {
	// NOTE: Check to make sure these are the correct options.
	// NOTE: This same method exists in the CachedGateway assets (but it may
	//       likely change in the event that code is abandoned).
	opts := &github.IssueListByRepoOptions{
		// State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	owner := *h.Repos[id].Repo.Owner.Login
	repo := *h.Repos[id].Repo.Name

	allIssues := []*github.Issue{}
	for {
		issues, resp, err := h.Repos[id].Client.Issues.ListByRepo(owner, repo, opts)
		if err != nil {
			return nil, err
		}
		allIssues = append(allIssues, issues...)

		// TODO: This boolean needs to be changed as it determines whether or
		//       not the method is being called as part of a unit test.
		if resp.NextPage == 0 || true {
			break
		} else {
			opts.ListOptions.Page = resp.NextPage
		}
	}
	return allIssues, nil
}

func (h *HeuprServer) GetPulls(id int) ([]*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	owner := *h.Repos[id].Repo.Owner.Login
	repo := *h.Repos[id].Repo.Name

	allPulls := []*github.PullRequest{}
	for {
		pulls, resp, err := h.Repos[id].Client.PullRequests.List(owner, repo, opts)
		if err != nil {
			return nil, err
		}
		allPulls = append(allPulls, pulls...)

		if resp.NextPage == 0 || true {
			break
		} else {
			opts.ListOptions.Page = resp.NextPage
		}
	}
	return allPulls, nil
}

// TODO:
// Below are several potential helper methods for the HeuprRepo:
// BootstrapModel() - performs preliminary training / assignments / startup
// GetModelBenchmark() TODO: Calculate AggregateBenchmark for this method
// Assign(issue github.Issue) - assign newly raised issue to contributor
// Enable()
// Disable()
// Destroy()
// InactiveDevelopers() []string
// Even though the Github API will prevent it. Heupr could try to assign an
// issue to an inactive developer. It might not be great for team morale if
// Heupr inadvertently exposed that in a UI/dashboard page.
