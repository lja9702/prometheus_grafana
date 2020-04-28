# prometheus_grafana
kubernetes상에서 prometheus와 grafana pod를 자동으로 배포해주는 모듈로
service, deployment, clusterRole, config map이 한번에 자동 배포되며,
grafana의 경우 prometheus와 연결될 수 있는 configMap.yaml을 사용하였습니다.

JSON으로 namespace, nodeport, docker image version 등을 넘겨주어 자동 배포할 수 있습니다.

# Getting Started
- go version: 1.13
- common.go: prometheus.go와 grafana.go에서 모두 사용하는 모듈을 모아놓은 파일
- prometheus.go: prometheus배포를 위한 사용자정의 json 스펙 정의 및 배포를 위한 메인 모듈을 모아놓은 파일
- grafana.go: grafana배포를 위한 사용자정의 json 스펙 정의 및 배포를 위한 메인 모듈을 모아놓은 파일

- module 추가
```
pg "github.com/lja9702/prometheus_grafana"
```
## prometheus 배포
* JSON 형식
```
type PrometheusSpec struct {
	NamespaceName string `json:"namespaceName"` //namespace 명 (defult: monitoring)
	ImgVersion    string `json:"imgVersion"`    //prometheus image version (defalt: latest)
	ScrapeInterv  string `json:"scrapeInterv"`  //prometheus가 스크랩을 요청하는 시간 간격(default: 15s)
	NodePort   string `json:"nodePort"`   //(default: 30000)
}
```
* 배포 모듈 사용
```
var prometheusSpec = PrometheusSpec{  //예시
    ScrapeInterv:  "15s",
    NodePort:   "30000",
    NamespaceName: "monitoring123",
    ImgVersion:    "v2.12.0",
}
//custom yaml 파일을 만들기 위한 base yaml 파일 읽기
gitPath := "https://raw.githubusercontent.com/lja9702/prometheus_grafana/master/origin_yaml_list/"
//Prometheus pod 생성
pg.CreatePrometheus(prometheusSpec, gitPath)
```
## grafana 배포
* JSON 형식
```
type GrafanaSpec struct {
	NamespaceName  string `json:"namespaceName"`  //namespace 명 (defult: monitoring)
	ImgVersion     string `json:"imgVersion"`     //prometheus image version (defalt: latest)
	RequestsMemory string `json:"requestsMemory"` //request는 컨테이너가 생성될때 요청하는 리소스 양 (defalt: 1Gi)
	RequestsCpu   string `json:"requestsCpu"`    //default: 500m
	LimitsMemory   string `json:"limitsMemory"`   //리소스가 더 필요한 경우 추가로 더 사용할 수 있는 부분 (default: 2Gi)
	LimitsCpu      string `json:"limitsCpu"`      //default: 1000m
	NodePort	string `json:"nodePort"`	//default: 32000
}
```
* 배포 모듈 사용
```
var grafanaSpec = GrafanaSpec{   //예시
  		NamespaceName:  "monitoring123",
  		ImgVersion:     "latest",
  		RequestsMemory: "1Gi",
  		RequestsCpu:    "500m",
  		LimitsMemory:   "2Gi",
  		LimitsCpu:      "1000m",
  		NodePort:	"32000",
}

//custom yaml 파일을 만들기 위한 base yaml 파일 읽기
gitPath := "https://raw.githubusercontent.com/lja9702/prometheus_grafana/master/origin_yaml_list/"
//Grafana pod 생성
pg.CreateGrafana(grafanaSpec, gitPath)
```

**참고: default라고 적힌 값은 진행한 프로젝트 frontend에서 지정한 임의의 값이며, 현 모듈에서는 적용되지 않습니다.**
