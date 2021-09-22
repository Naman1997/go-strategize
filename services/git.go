package services

import (
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	ip              string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	urlSchema       string = `((ftp|https?):\/\/)`
	urlUsername     string = `(\S+(:\S*)?@)`
	urlPath         string = `((\/|\?|#)[^\s]*)`
	urlPort         string = `(:(\d{1,5}))`
	urlIP           string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	urlSubdomain    string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	urlString              = `^` + urlSchema + `?` + urlUsername + `?` + `((` + urlIP + `|(\[` + ip + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + urlSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + urlPort + `?` + urlPath + `?$`
	maxURLRuneCount        = 2083
	minURLRuneCount        = 3
)

var (
	rxURL = regexp.MustCompile(urlString)
)

/*
CloneRepos clones repos if URLs provided are valid
Also makes sure that the repo is not already
cloned by checking the current dir for folder
whose name matches the repo name
*/
func CloneRepos(terraformRepo string, ansibleRepo string, homedir string) {
	var wg sync.WaitGroup
	_, terraformExists := os.Stat(FormatRepo(terraformRepo))
	_, ansibleExists := os.Stat(FormatRepo(ansibleRepo))
	if terraformExists != nil && len(terraformRepo) > 0 {
		wg.Add(1)
		go cloneTemplateRepos(terraformRepo, &wg)
	} else {
		ColorPrint(WARN, "Terraform template repo is already cloned!")

	}
	if ansibleExists != nil && len(ansibleRepo) > 0 {
		wg.Add(1)
		go cloneTemplateRepos(ansibleRepo, &wg)
	} else {
		ColorPrint(WARN, "Ansible template repo is already cloned!")
	}
	wg.Wait()
}

//FormatRepo returns folder name for a given URL
func FormatRepo(repo string) string {
	repo = repo[strings.LastIndex(repo, "/")+1:]
	return strings.Replace(repo, ".git", "", -1) + "/"
}

func cloneTemplateRepos(path string, wg *sync.WaitGroup) {
	cmd0 := "git"
	cmd1 := "clone"
	cmd := exec.Command(cmd0, cmd1, path)
	_, err := cmd.Output()

	if err != nil {
		defer wg.Done()
		ColorPrint(ERROR, "%v", err)
	}

	ColorPrint(INFO, "Finished cloning "+path)
	defer wg.Done()
}

/*
IsURL checks if the provided URL is valid
function from asaskevich/govalidator
*/
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
