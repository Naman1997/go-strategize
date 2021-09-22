package services

import (
	"os"
	"os/exec"
)

/*
TerraformInit: Executes terraform init on the path using
the -chdir option for terraform apply
Expects terraform binary to be present
in PATH
*/
func TerraformInit(path string) {
	cmd0 := "terraform"
	cmd1 := "-chdir=" + path
	cmd2 := "init"
	cmd := exec.Command(cmd0, cmd1, cmd2)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		ColorPrint(ERROR, "%v", err)
	}

	ColorPrint(INFO, "Finished executing terraform init")
}

/*
TerraformApply: Executes terraform apply on the path using
the -chdir option for terraform apply
Expects terraform binary to be present
in PATH
*/
func TerraformApply(path string) {
	cmd0 := "terraform"
	cmd1 := "-chdir=" + path
	cmd2 := "apply"
	cmd3 := "-auto-approve"
	cmd := exec.Command(cmd0, cmd1, cmd2, cmd3)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		ColorPrint(ERROR, "%v", err)
	}

	ColorPrint(INFO, "Finished executing terraform apply")
}
