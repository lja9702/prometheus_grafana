package prometheus_grafana

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type PrometheusSpec struct {
	NamespaceName string `json:"namespaceName"` //namespace 명 (defult: monitoring)
	ImgVersion    string `json:"imgVersion"`    //prometheus image version (defalt: latest)
	ScrapeInterv  string `json:"scrapeInterv"`  //prometheus가 스크랩을 요청하는 시간 간격(default: 15s)
	NodePort   string `json:"nodePort"`   //(default: 30000)
}

func WordbyWordScanPrometheus(fileName string, spec *PrometheusSpec) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	text := string(content)

	// Create replacer with pairs as arguments.
	replacerArg := strings.NewReplacer(
		"{{namespaceName}}", spec.NamespaceName,
		"{{nodePort}}", spec.NodePort,
		"{{scrapeInterv}}", spec.ScrapeInterv,
		"{{imgVersion}}", spec.ImgVersion,
	)

	// Replace all pairs.
	newFormatYmlString := replacerArg.Replace(text)
	//fmt.Println(newFormatYmlString)

	if err = ioutil.WriteFile("custom_"+fileName, []byte(newFormatYmlString), 0666); err != nil {
		fmt.Println(err)
	}
}

func promApplyYamlFileCmd(gitPath string, fileName string, spec *PrometheusSpec, option string) {
	if non_err := DownloadFile(fileName, gitPath); non_err == nil {
			WordbyWordScanPrometheus(fileName, spec)
			applyYamlFileCmd(fileName, option, spec.NamespaceName)
			deleteFile(fileName)
	} else{
		panic(err)
	}
}

func CreatePrometheus(prometheusSpec PrometheusSpec, gitPath string) {
	// /////////test config
	// var Prometheus_spec = PrometheusSpec{
	// 	scrapeInterv:  "15s",
	// 	nodePort:   "30000",
	// 	namespaceName: "monitoring",
	// 	imgVersion:    "v2.12.0",
	// }
	// ///////////////////

	//connectToClusterCmd()

	createNamespaceCmd(prometheusSpec.NamespaceName)
	promApplyYamlFileCmd(gitPath, "prom_clusterRole.yaml", &prometheusSpec, "")
	promApplyYamlFileCmd(gitPath, "prom_config_map.yaml", &prometheusSpec, "")
	promApplyYamlFileCmd(gitPath, "prom_deployment.yaml", &prometheusSpec, "")

	//////////////////Check deployment file
	cmd := exec.Command("kubectl", "get", "deployments", "--namespace="+prometheusSpec.NamespaceName)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
	//////////////////////////////////////////
	promApplyYamlFileCmd(gitPath, "prom_service.yaml", &prometheusSpec, "--namespace=")
	return true
}
