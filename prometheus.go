package prometheus_grafana

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type PrometheusSpec struct {
	namespaceName string `json:"namespaceName"` //namespace 명 (defult: monitoring)
	imgVersion    string `json:"imgVersion"`    //prometheus image version (defalt: latest)
	scrapeInterv  string `json:"scrapeInterv"`  //prometheus가 스크랩을 요청하는 시간 간격(default: 15s)
	targetHosts   string `json:"targetHosts"`   //(default: ['localhost:9090'])
}

func WordbyWordScanPrometheus(originPath string, customPath string, fileName string, spec *PrometheusSpec) {
	content, err := ioutil.ReadFile(originPath + fileName)
	if err != nil {
		fmt.Println(err)
	}
	text := string(content)

	// Create replacer with pairs as arguments.
	replacerArg := strings.NewReplacer(
		"{{namespaceName}}", spec.namespaceName,
		"{{targetHosts}}", spec.targetHosts,
		"{{scrapeInterv}}", spec.scrapeInterv,
		"{{imgVersion}}", spec.imgVersion,
	)

	// Replace all pairs.
	newFormatYmlString := replacerArg.Replace(text)
	//fmt.Println(newFormatYmlString)

	if err = ioutil.WriteFile(customPath+fileName, []byte(newFormatYmlString), 0666); err != nil {
		fmt.Println(err)
	}
}

func promApplyYamlFileCmd(originPath string, customPath string, fileName string, spec *PrometheusSpec, option string) {
	WordbyWordScanPrometheus(originPath, customPath, fileName, spec)
	applyYamlFileCmd(customPath, fileName, option, spec.namespaceName)
}

func createPrometheus(prometheusSpec PrometheusSpec) {
	// /////////test config
	// var Prometheus_spec = PrometheusSpec{
	// 	scrapeInterv:  "15s",
	// 	targetHosts:   "['localhost:9090']",
	// 	namespaceName: "monitoring",
	// 	imgVersion:    "v2.12.0",
	// }
	// ///////////////////

	connectToClusterCmd()
	createNamespaceCmd(prometheusSpec.namespaceName)
	promApplyYamlFileCmd("origin_yaml_list/", "custom_yaml_list/", "prom_clusterRole.yaml", &prometheusSpec, "")
	promApplyYamlFileCmd("origin_yaml_list/", "custom_yaml_list/", "prom_config_map.yaml", &prometheusSpec, "")
	promApplyYamlFileCmd("origin_yaml_list/", "custom_yaml_list/", "prom_deployment.yaml", &prometheusSpec, "")

	//////////////////Check deployment file
	cmd := exec.Command("kubectl", "get", "deployments", "--namespace="+prometheusSpec.namespaceName)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
	//////////////////////////////////////////
	promApplyYamlFileCmd("origin_yaml_list/", "custom_yaml_list/", "prom_service.yaml", &prometheusSpec, "--namespace=")
}
