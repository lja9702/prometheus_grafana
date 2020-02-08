package prometheus_grafana

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type GrafanaSpec struct {
	namespaceName  string `json:"namespaceName"`  //namespace 명 (defult: monitoring)
	imgVersion     string `json:"imgVersion"`     //prometheus image version (defalt: latest)
	requestsMemory string `json:"requestsMemory"` //request는 컨테이너가 생성될때 요청하는 리소스 양 (defalt: 1Gi)
	requests_cpu   string `json:"requestsCpu"`    //default: 500m
	limitsMemory   string `json:"limitsMemory"`   //리소스가 더 필요한 경우 추가로 더 사용할 수 있는 부분 (default: 2Gi)
	limitsCpu      string `json:"limitsCpu"`      //default: 1000m
}

func WordbyWordScanGrafana(originPath string, customPath string, fileName string, spec *GrafanaSpec) {
	content, err := ioutil.ReadFile(originPath + fileName)
	if err != nil {
		fmt.Println(err)
	}
	text := string(content)

	// Create replacer with pairs as arguments.
	replacerArg := strings.NewReplacer(
		"{{namespaceName}}", spec.namespaceName,
		"{{imgVersion}}", spec.imgVersion,
		"{{requestsMemory}}", spec.requestsMemory,
		"{{requests_cpu}}", spec.requests_cpu,
		"{{limitsMemory}}", spec.limitsMemory,
		"{{limitsCpu}}", spec.limitsCpu,
	)

	// Replace all pairs.
	newFormatYmlString := replacerArg.Replace(text)
	//fmt.Println(newFormatYmlString)

	if err = ioutil.WriteFile(customPath+fileName, []byte(newFormatYmlString), 0666); err != nil {
		fmt.Println(err)
	}
}

func grafApplyYamlFileCmd(originPath string, customPath string, fileName string, spec *GrafanaSpec, option string) {
	WordbyWordScanGrafana(originPath, customPath, fileName, spec)
	applyYamlFileCmd(customPath, fileName, option, spec.namespaceName)
}

func createGrafana(grafanaSpec GrafanaSpec, originPath string, customPath string) {
	// ////////test config
	// var grafana_spec = GrafanaSpec{
	// 	namespaceName:  "monitoring",
	// 	imgVersion:     "latest",
	// 	requestsMemory: "1Gi",
	// 	requests_cpu:    "500m",
	// 	limitsMemory:   "2Gi",
	// 	limitsCpu:      "1000m",
	// }
	// ///////////

	connectToClusterCmd()
	createNamespaceCmd(grafanaSpec.namespaceName)
	///grafana + prometheus 라면
	grafApplyYamlFileCmd(originPath, customPath, "graf_with_prom_config_map.yaml", &grafanaSpec, "")

	///grafana만이라면 config map 스킵
	//
	grafApplyYamlFileCmd(originPath, customPath, "graf_deployment.yaml", &grafanaSpec, "")

	//////////////////Check deployment file
	cmd := exec.Command("kubectl", "get", "deployments", "--namespace="+grafanaSpec.namespaceName)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
	//////////////////////////////////////////
	grafApplyYamlFileCmd(originPath, customPath, "graf_service.yaml", &grafanaSpec, "--namespace=")
}
