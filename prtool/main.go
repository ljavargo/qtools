package main

import (
	"github.com/xanzy/go-gitlab"
	"fmt"
	"strings"
	"flag"
)

func main() {
	h := flag.Bool("h",false,"帮助")
	baseUrl := flag.String("u","https://git.example.com","必填，gitlab主页url地址")
	token := flag.String("tk","","必填，private token")
	group := flag.String("g","","group名称")
	project := flag.String("p","","必填，project名称")
	title := flag.String("tt","create pr by prtool","pr title")
	source := flag.String("s","","source branch")
	target := flag.String("t","master","target branch, 支持两种模式，批量：'branch1,branch2,branch3' 链式：'branch1>branch2>branch3' (注意：要加单引号防止>被处理为重定向符号) ")
	closeAll := flag.Bool("closeAll", false, "关闭当前用户指定project的所有PR")

	flag.Parse()
	flag.Usage = usage
	if *h {
		flag.Usage()
		return
	}

	if *baseUrl == "" || *token == "" || *project == "" {
		flag.Usage()
		return
	}

	if *closeAll {
		closePRs(*baseUrl, *token, *group, *project)
	}else {
		createPRs(*baseUrl, *token, *group, *project, *title, *source, *target)
	}
}

func usage() {
	fmt.Printf(`prtool version: 1.0.0
Options:
`)
	flag.PrintDefaults()
}

func closePRs(baseUrl, token, group, project string)  {
	git, projectId := findProjet(baseUrl + "/api/v3", token, group, project)
	if projectId == 0 {
		fmt.Println("not found project")
		return
	}
	me, _, _ := git.Users.CurrentUser()
	openedState := "opened"
	closeState := "close"
	requests, _, e := git.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
		AuthorID: &(me.ID),
		State: &openedState,
	})
	if e != nil {
		fmt.Errorf("ListProjectMergeRequests prs error", e)
		return
	}

	for _,pr := range requests {
		git.MergeRequests.UpdateMergeRequest(projectId, pr.ID, &gitlab.UpdateMergeRequestOptions{StateEvent: &closeState})
		//git.MergeRequests.DeleteMergeRequest(projectId, pr.ID)
		fmt.Printf("close successful: %s\n", pr.WebURL)
	}
}

func findProjet(baseUrl, token, group, project string) (git *gitlab.Client, projectId int) {
	git = gitlab.NewClient(nil, token)
	git.SetBaseURL(baseUrl)
	pros,_,err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{Search:&project})

	if err != nil {
		fmt.Errorf("ListProjects error", err)
		return
	}

	for _,pro:=range pros {
		if pro.Name == project && pro.Namespace.Name == group {
			projectId = pro.ID
		}
	}
	return
}

func createPRs(baseUrl, token, group, project, title, source, target string)  {
	git,projectId := findProjet(baseUrl + "/api/v4", token, group, project)
	sourceSlice := []string{}
	targetSlice := []string{}
	if projectId == 0 {
		fmt.Println("not found project")
		return
	}

	if strings.ContainsRune(target, '>') {
		targetSlice = strings.Split(target, ">")
		sourceSlice = append(sourceSlice, source)
		sourceSlice = append(sourceSlice, targetSlice[0: len(targetSlice)-1]...)

	} else {
		targetSlice = strings.Split(target, ",")
		for i:=0;i<len(targetSlice);i++  {
			sourceSlice = append(sourceSlice, source)
		}
	}

	for i,t:=range targetSlice  {
		createMergeRequest(git, projectId, title, sourceSlice[i], t)
	}

}

func createMergeRequest(git *gitlab.Client, projectId int, title string, sourceBranch string, targetBranch string) {
	request, resp, err := git.MergeRequests.CreateMergeRequest(projectId, &gitlab.CreateMergeRequestOptions{
		Title:           &title,
		SourceBranch:    &sourceBranch,
		TargetBranch:    &targetBranch,
		TargetProjectID: &projectId,
	})
	if err != nil {
		fmt.Errorf("create pr (%s -> %s) resp %v error:",sourceBranch, targetBranch, resp, err)
		return
	}
	fmt.Printf("%s -> %s: %s\n", sourceBranch, targetBranch, request.WebURL)
}