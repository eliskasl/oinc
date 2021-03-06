package steps

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mfojtik/oinc/pkg/log"
	"github.com/mfojtik/oinc/pkg/util"
)

type InstallRouterStep struct {
	DefaultStep
}

func (*InstallRouterStep) String() string { return "install-router" }

func addRouterUser() error {
	out, err := util.RunAdminOc("get", "scc", "privileged", "-o", "json")
	if err != nil {
		return err
	}
	result := ""
	for _, line := range strings.Split(string(out), "\n") {
		result += line + "\n"
		if strings.Contains(line, `"users":`) {
			result += `"system:serviceaccount:default:router",` + "\n"
		}
	}
	f, err := ioutil.TempFile("", "scc")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	ioutil.WriteFile(f.Name(), []byte(result), 0600)
	_, err = util.RunAdminOc("replace", "scc", "privileged", "-f", f.Name())
	return err
}

func (*InstallRouterStep) Execute() error {
	log.Info("Installing HAProxy router ....")
	if err := addRouterUser(); err != nil {
		return err
	}
	_, err := util.RunOAdm("router",
		"--create",
		"--credentials", filepath.Join(util.MasterConfigPath, "openshift-router.kubeconfig"),
		"--service-account=router",
	)
	return err
}
