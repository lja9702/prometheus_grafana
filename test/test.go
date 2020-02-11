package main

import(
  "net/http"
  "io/ioutil"
  "fmt"
)
func getContent(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("GET error: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Read body: %v", err)
    }

    return data, nil
}

func main() {
    gitPath :="https://raw.githubusercontent.com/lja9702/prometheus_grafana/master/origin_yaml_list/"
    // var prometheusSpec = PrometheusSpec{
    // 	ScrapeInterv:  "15s",
    // 	NodePort:   "30000",
    // 	NamespaceName: "monitoring123",
    // 	ImgVersion:    "v2.12.0",
    // }

    //CreatePrometheus(prometheusSpec, gitPath)

    ////////test config
  	var grafanaSpec = GrafanaSpec{
  		NamespaceName:  "monitoring123",
  		ImgVersion:     "latest",
  		RequestsMemory: "1Gi",
  		RequestsCpu:    "500m",
  		LimitsMemory:   "2Gi",
  		LimitsCpu:      "1000m",
  		NodePort:	"32000",
  	}
  	///////////
    CreateGrafana(grafanaSpec, gitPath)
}
