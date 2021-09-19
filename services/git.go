package services

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	IP              string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLSchema       string = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername     string = `(\S+(:\S*)?@)`
	URLPath         string = `((\/|\?|#)[^\s]*)`
	URLPort         string = `(:(\d{1,5}))`
	URLIP           string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	URLSubdomain    string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URL                    = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	maxURLRuneCount        = 2083
	minURLRuneCount        = 3
)

var (
	rxURL = regexp.MustCompile(URL)
)

//fixme: Improve log messages. Say repo not provided if len(repo) == 0
func CloneRepos(terraform_repo string, ansible_repo string, homedir string) {
	var wg sync.WaitGroup
	_, terraform_exists := os.Stat(FormatRepo(terraform_repo))
	_, ansible_exists := os.Stat(FormatRepo(ansible_repo))
	if terraform_exists != nil && len(terraform_repo) > 0 {
		wg.Add(1)
		go clone_template_repos(terraform_repo, &wg)
	} else {
		fmt.Println("[SKIP] Terraform template repo is already cloned!")
	}
	if ansible_exists != nil && len(ansible_repo) > 0 {
		wg.Add(1)
		go clone_template_repos(ansible_repo, &wg)
	} else {
		fmt.Println("[SKIP] Ansible template repo is already cloned!")
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
		defer wg.Done()
		log.Fatalf("[ERROR] %v", err)
	}

	fmt.Println("[INFO] Finished cloning", path)
	defer wg.Done()
}

//Function from asaskevich/govalidator
func IsURL(str string) bool {
	if str == "" || utf8.RuneCountInString(str) >= maxURLRuneCount || len(str) <= minURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	strTemp := str
	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
		// support no indicated urlscheme but with colon for port number
		// http:// is appended so url.Parse will succeed, strTemp used so it does not impact rxURL.MatchString
		strTemp = "http://" + str
	}
	u, err := url.Parse(strTemp)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)
}
