package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func clone_template_repos(cmd2 string, wg *sync.WaitGroup) {
	cmd0 := "git"
	cmd1 := "clone"
	cmd := exec.Command(cmd0, cmd1, "https://github.com/Naman1997/"+cmd2)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		defer wg.Done()
		return
	}

	fmt.Println("Finished cloning", cmd2)
	defer wg.Done()
}

func terraform_apply(path string) {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		log.Fatalf("error creating temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpDir, false))
	if err != nil {
		log.Fatalf("error locating Terraform binary: %s", err)
	}

	tf, err := tfexec.NewTerraform(path, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	var wg sync.WaitGroup
	var template bool = false
	var terraform_repo string = "proxmox-terraform-template-k8s"
	var ansible_repo string = "cluster-management"
	fmt.Print("Clone and execute default proxmox template?[Y/N]")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if strings.EqualFold(input.Text(), "Y") {
		terraform_clone_exists, _ := exists("./" + terraform_repo)
		ansible_clone_exists, _ := exists("./" + ansible_repo)
		if !terraform_clone_exists {
			wg.Add(1)
			go clone_template_repos(terraform_repo, &wg)
		} else {
			fmt.Println("Skip: Terraform template repo is already cloned!")
		}
		if !ansible_clone_exists {
			wg.Add(1)
			go clone_template_repos(ansible_repo, &wg)
		} else {
			fmt.Println("Skip: Ansible template repo is already cloned!")
		}
		template = true
		wg.Wait()
	}

	// if template {
	// 	terraform_apply(terraform_repo)
	// }

	if template {
		fmt.Println("Execution completed for template!")
	} else {
		fmt.Println("Sorry, non-template execution is not yet supported")
	}
}
