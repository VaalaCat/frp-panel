package conf

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"

	"github.com/VaalaCat/frp-panel/pb"
)

var (
	gitVersion = "dev-build"
	gitCommit  = ""
	gitBranch  = ""
	buildDate  = "1970-01-01T00:00:00Z"
)

type VersionInfo struct {
	GitVersion string `json:"gitVersion" yaml:"gitVersion"`
	GitCommit  string `json:"gitCommit" yaml:"gitCommit"`
	GitBranch  string `json:"gitBranch" yaml:"gitBranch"`
	BuildDate  string `json:"buildDate" yaml:"buildDate"`
	GoVersion  string `json:"goVersion" yaml:"goVersion"`
	Compiler   string `json:"compiler" yaml:"compiler"`
	Platform   string `json:"platform" yaml:"platform"`
}

func (v *VersionInfo) String() string {
	tempStr := "BinVersion: {{.GitVersion}}\nGitCommit: {{.GitCommit}}\nBuildDate: {{.BuildDate}}\nGoVersion: {{.GoVersion}}\nCompiler: {{.Compiler}}\nPlatform: {{.Platform}}"
	temp, err := template.New("version").Parse(tempStr)
	if err != nil {
		return ""
	}
	var result bytes.Buffer
	err = temp.Execute(&result, v)
	if err != nil {
		return ""
	}
	return result.String()
}

func (v *VersionInfo) ToProto() *pb.ClientVersion {
	return &pb.ClientVersion{
		GitVersion: v.GitVersion,
		GitCommit:  v.GitCommit,
		GitBranch:  v.GitBranch,
		BuildDate:  v.BuildDate,
		GoVersion:  v.GoVersion,
		Compiler:   v.Compiler,
		Platform:   v.Platform,
	}
}

func GetVersion() *VersionInfo {
	return &VersionInfo{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		GitBranch:  gitBranch,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
