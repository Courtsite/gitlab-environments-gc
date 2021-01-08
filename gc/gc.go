package gc

import (
	"log"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

var protectedEnvironmentNames []string = []string{
	"main",
	"master",
	"production",
	"protect",
}

func CleanupEnvironments(client *gitlab.Client, projectID interface{}) error {
	log.Printf("🌲 Listing environments for %s\n", projectID)

	opts := &gitlab.ListEnvironmentsOptions{
		PerPage: 200,
	}
	envs, _, err := client.Environments.ListEnvironments(projectID, opts)
	if err != nil {
		log.Fatalln(err)
	}

	skippedEnvs := 0
	deletedEnvs := 0

	defer func() {
		log.Printf("Total environments handled: %d\n", len(envs))
	}()

	for _, env := range envs {
		if isProtectedEnvironment(env) {
			skippedEnvs += 1
			log.Printf("⏭️ Skipping %s (%d) as it is a protected environment\n", env.Name, env.ID)
			continue
		}

		detailedEnv, _, err := client.Environments.GetEnvironment(env.Project.ID, env.ID)
		if err != nil {
			log.Printf("❌ Failed to get environment details %s (%d): %+v\n", env.Name, env.ID, err)
			return err
		}

		if isActiveEnvironment(detailedEnv) {
			skippedEnvs += 1
			log.Printf("⏭️ Skipping %s (%d) as it is an active environment\n", detailedEnv.Name, detailedEnv.ID)
			continue
		}

		log.Printf("🗑️ Deleting environment: %s (%d)\n", detailedEnv.Name, detailedEnv.ID)

		_, err = client.Environments.StopEnvironment(env.Project.ID, detailedEnv.ID)
		if err != nil {
			log.Printf("❌ Failed to stop environment %s (%d): %+v\n", detailedEnv.Name, detailedEnv.ID, err)
			return err
		}

		_, err = client.Environments.DeleteEnvironment(env.Project.ID, detailedEnv.ID)
		if err != nil {
			log.Printf("❌ Failed to delete environment %s (%d): %+v\n", detailedEnv.Name, detailedEnv.ID, err)
			return err
		}

		deletedEnvs += 1

		// This is to avoid hitting any rate limits.
		time.Sleep(time.Second / 2)
	}

	log.Printf("Skipped environments: %d\n", skippedEnvs)
	log.Printf("Deleted environments: %d\n", deletedEnvs)

	return nil
}

func isProtectedEnvironment(env *gitlab.Environment) bool {
	normalisedName := strings.ToLower(env.Name)
	for _, protectedEnvironmentName := range protectedEnvironmentNames {
		if strings.Contains(normalisedName, protectedEnvironmentName) {
			return true
		}
	}
	return false
}

func isActiveEnvironment(env *gitlab.Environment) bool {
	twoWeeks := time.Hour * 24 * 7 * 2
	return env.LastDeployment != nil && env.LastDeployment.UpdatedAt != nil && env.LastDeployment.UpdatedAt.UTC().Add(twoWeeks).After(time.Now().UTC())
}
