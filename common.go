package prometheus_grafana

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"net/http"
	"io"
)

/////////////////////////////////////////Check cmd output///////////
func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

///////////////////////////////////////////////////////////////////////

///////////////////////////////////////////Commend line////////////////
func connectToClusterCmd() {
	//Connect to the Cluster
	cmd := exec.Command("bash", "-c", "ACCOUNT=$(gcloud info --format='value(config.account)')")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)

	setKubectlAcountCmd := exec.Command("kubectl", "create", "clusterrolebinding", "owner-cluster-admin-binding",
		"--clusterrole", "cluster-admin", "--user", "$ACCOUNT")

	printCommand(setKubectlAcountCmd)
	output, setKubectlAcountErr := setKubectlAcountCmd.CombinedOutput()
	printError(setKubectlAcountErr)
	printOutput(output)
}

func createNamespaceCmd(namespace_name string) {
	cmd := exec.Command("kubectl", "create", "namespace", namespace_name)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
}

func applyYamlFileCmd(fileName string, option string, namespace_name string) {
	var cmd *exec.Cmd
	if strings.Compare(option, "--namespace=") == 0 {
		cmd = exec.Command("kubectl", "apply", "-f", fileName, option+namespace_name)
	} else {
		cmd = exec.Command("kubectl", "apply", "-f", fileName)
	}

	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
}

///////////////////////////////////////////////////////////////////////


///////////////Download yaml file
func DownloadFile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) { return }

	fmt.Println("==> done deleting file")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
