package handlers

import (
	rsp "github.com/SyntSugar/ss-infra-go/api/response"
	"github.com/SyntSugar/ss-infra-go/whoami"
	"github.com/gin-gonic/gin"
)

type build struct {
	Number   string `json:"number"`
	Datetime string `json:"datetime"`
}

type commit struct {
	Hash   string `json:"hash"`
	Branch string `json:"branch"`
}

type buildInfo struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Build   build  `json:"build"`
	Commit  commit `json:"commit"`
}

func Whoami(ctx *gin.Context) {
	rsp.ResponseWithOK(ctx, buildInfo{
		Service: whoami.Name(),
		Version: whoami.Version(),
		Build: build{
			Number:   whoami.Number(),
			Datetime: whoami.BuildAt(),
		},
		Commit: commit{
			Hash:   whoami.CommitHash(),
			Branch: whoami.GitBranch(),
		},
	})
}
