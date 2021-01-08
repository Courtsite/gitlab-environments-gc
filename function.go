package function

import (
	"log"
	"net/http"
	"os"

	"github.com/courtsite/gitlab-environments-gc/gc"
	"github.com/xanzy/go-gitlab"
)

func F(w http.ResponseWriter, r *http.Request) {
	gitlabAPIToken := os.Getenv("GITLAB_API_TOKEN")
	if gitlabAPIToken == "" {
		log.Fatalln("`GITLAB_API_TOKEN` is not set in the environment")
	}

	gitlabProjectID := os.Getenv("GITLAB_PROJECT_ID")
	if gitlabProjectID == "" {
		log.Fatalln("`GITLAB_PROJECT_ID` is not set in the environment")
	}

	client, err := gitlab.NewClient(gitlabAPIToken)
	if err != nil {
		log.Fatalln(err)
	}

	_ = gc.CleanupEnvironments(client, gitlabProjectID)
	w.WriteHeader(http.StatusOK)
}
