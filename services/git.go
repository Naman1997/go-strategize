package services

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

//fixme: Improve log messages. Say repo not provided if len(repo) == 0
func CloneRepos(terraform_repo string, ansible_repo string) {
	var wg sync.WaitGroup
	terraform_exists, _ := Exists(FormatRepo(terraform_repo))
	ansible_exists, _ := Exists(FormatRepo(ansible_repo))
	if !terraform_exists && len(terraform_repo) > 0 {
		wg.Add(1)
		go clone_template_repos(terraform_repo, &wg)
	} else {
		fmt.Println("SKIP: Terraform template repo is already cloned!")
	}
	if !ansible_exists && len(ansible_repo) > 0 {
		wg.Add(1)
		go clone_template_repos(ansible_repo, &wg)
	} else {
		fmt.Println("SKIP: Ansible template repo is already cloned!")
	}
	wg.Wait()
}

func FormatRepo(repo string) string {
	repo = repo[strings.LastIndex(repo, "/")+1:]
	return strings.Replace(repo, ".git", "", -1) + "/"
}

func clone_template_repos(path string, wg *sync.WaitGroup) {
	cmd0 := "git"
	cmd1 := "clone"
	cmd := exec.Command(cmd0, cmd1, path)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		defer wg.Done()
		return
	}

	fmt.Println("Finished cloning", path)
	defer wg.Done()
}
