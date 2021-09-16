package main

import (
	// "bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"log"
	"sync"

	ssh "github.com/Naman1997/go-stratergize/services"
)

func main() {
	// terraform_vars_file := flag.String("var-file", "", "Path to .tfvars file")
	ssh_username := flag.String("ssh-user", "naman", "Username for SSH")
	ssh_key := flag.String("ssh-key", "/home/naman/.ssh/id_rsa", "Path of SSH public key")
	flag.Parse()
	conn, err := ssh.Connect(*ssh_username, *ssh_key)
	if err != nil {
		log.Fatal(err)
	}

	output, err := conn.SendCommands("sleep 2", "hostname")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

	// var wg sync.WaitGroup
	// var template bool = false
	// var terraform_repo string = "proxmox-terraform-template-k8s"
	// var ansible_repo string = "cluster-management"
	// fmt.Print("Clone and execute default proxmox template?[Y/N]")
	// input := bufio.NewScanner(os.Stdin)
	// input.Scan()
	// dir, _ := os.Getwd()

	// //Clone template repos
	// if strings.EqualFold(input.Text(), "Y") {
	// 	terraform_clone_exists, _ := exists(terraform_repo + "/")
	// 	ansible_clone_exists, _ := exists(ansible_repo + "/")
	// 	if !terraform_clone_exists {
	// 		wg.Add(1)
	// 		go clone_template_repos(terraform_repo, &wg)
	// 	} else {
	// 		fmt.Println("Skip: Terraform template repo is already cloned!")
	// 	}
	// 	if !ansible_clone_exists {
	// 		wg.Add(1)
	// 		go clone_template_repos(ansible_repo, &wg)
	// 	} else {
	// 		fmt.Println("Skip: Ansible template repo is already cloned!")
	// 	}
	// 	template = true
	// 	wg.Wait()
	// }

	// //Copy over .tfvars file if specified
	// if len(*terraform_vars_file) > 0 {
	// 	tfvars_exists, _ := exists(*terraform_vars_file)
	// 	fmt.Println(tfvars_exists)
	// 	if tfvars_exists {
	// 		fmt.Println("Copying .tfvars file over to the required folder")
	// 		cpCmd := exec.Command("cp", *terraform_vars_file, dir+"/"+terraform_repo+"/terraform.tfvars")
	// 		_ = cpCmd.Run()
	// 	} else {
	// 		fmt.Println("Error: Provided path not found!")
	// 		os.Exit(1)
	// 	}
	// }

	// //Initialize and apply with terraform
	// if template {
	// 	terraform_init(terraform_repo, dir)
	// 	terraform_apply(terraform_repo, dir)
	// }

	// if template {
	// 	fmt.Println("Execution completed for template!")
	// } else {
	// 	fmt.Println("Sorry, non-template execution is not yet supported")
	// 	os.Exit(1)
	// }
}

func clone_template_repos(path string, wg *sync.WaitGroup) {
	cmd0 := "git"
	cmd1 := "clone"
	cmd := exec.Command(cmd0, cmd1, "https://github.com/Naman1997/"+path)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		defer wg.Done()
		return
	}

	fmt.Println("Finished cloning", path)
	defer wg.Done()
}

func terraform_init(path string, dir string) {
	cmd0 := "terraform"
	cmd1 := "-chdir=" + dir + "/" + path
	cmd2 := "init"
	cmd := exec.Command(cmd0, cmd1, cmd2)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: Unable to execute init with terraform!")
		os.Exit(1)
	}

	fmt.Println("INFO: Finished executing init stage", cmd.Stdin)
}

func terraform_apply(path string, dir string) {
	cmd0 := "terraform"
	cmd1 := "-chdir=" + dir + "/" + path
	cmd2 := "apply"
	cmd3 := "-auto-approve"
	cmd := exec.Command(cmd0, cmd1, cmd2, cmd3)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: Unable to create resources with terraform!")
		os.Exit(1)
	}

	fmt.Println("INFO: Finished executing apply stage", cmd.Stdin)
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
