package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jeffotoni/k8slog/sdk/config"
	"github.com/jeffotoni/k8slog/sdk/fmts"
)

// dasboard
// namespace | pods | cpu | mem
// namespace | servico | qnt pod | cpu | mem
// kubectl get hpa -n <namespace> <pod>
// kubectl top pod <pod> --use-protocol-buffers -n <namespace>
var (
	c = config.Config()
)

func main() {

	// log
	flog := "./kubectl.log"
	f, err := os.Open(flog)
	if err != nil {
		fmt.Printf("\n%s", err.Error())
		return
	}

	nspaceTotal := make(map[string]string)
	mnameSpace := make(map[string]string)
	//mname := make(map[string]string)

	fscan := bufio.NewScanner(f)
	fscan.Split(bufio.ScanLines)

	//fmt.Printf("\n%t\n", c.Cluster.NameSpace)
	//return

	var mn = make(map[string]string)
	for _, v := range c.Cluster.NameSpace {
		mn[v] = v
	}

	// NAMESPACE | NAME | READY | STATUS | RESTARTS | AGE
	for fscan.Scan() {
		line := fscan.Text()
		line = strings.TrimSpace(line)
		lineV := strings.Split(line, " ")

		inspace := 0
		iname := 0
		iready := 0
		istatus := 0
		irestarts := 0
		iage := 0

		var j int = 0
		for l, v := range lineV {
			if len(v) > 0 {
				switch j {
				case 0:
					inspace = l
				case 1:
					iname = l
				case 2:
					iready = l
				case 3:
					istatus = l
				case 4:
					irestarts = l
				case 5:
					iage = l
				}
				//println("l:", l, " - v:", v)
				j++
			}
		}

		//var pv []Pods
		pns := strings.TrimSpace(lineV[inspace])
		_, ok := mn[pns]
		if ok {
			pname := strings.TrimSpace(lineV[iname])
			piready := strings.TrimSpace(lineV[iready])
			pistatus := strings.TrimSpace(lineV[istatus])
			pirestarts := strings.TrimSpace(lineV[irestarts])
			piage := strings.TrimSpace(lineV[iage])

			mnameSpace[fmts.ConcatStr(pns, "#", pname)] = fmts.ConcatStr(pns, "#", pname, "#", piready, "#", pistatus, "#", pirestarts, "#", piage)
			nspaceTotal[pns] = pns

		}
	}

	f.Close()

	//ShowNameSpacePods(nspaceTotal, mnameSpace)
	ShowTotalServicePods(nspaceTotal, mnameSpace)
	// progress.Show()
	println("")

}

func ShowTotalServicePods(nspaceTotal, mnameSpace map[string]string) {
	t := table.NewWriter()
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "NameSpace", "Service", "Pods", "CPU", "MEM"})

	TotalPods(t, nspaceTotal, mnameSpace)

	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", len(mnameSpace), ""})
	t.SetStyle(table.StyleColoredBright)
	t.Render()

}

func ShowNameSpacePods(nspaceTotal, mnameSpace map[string]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "NameSpace", "POD"})

	NamePod(t, nspaceTotal, mnameSpace)

	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", len(mnameSpace), ""})
	t.Render()
}

func TotalPods(t table.Writer, nspaceTotal, mnameSpace map[string]string) {
	for i, n := range nspaceTotal {
		cpu, memory := sumCpu(n, mnameSpace)
		t.AppendRows([]table.Row{
			{i, n, sumService(n, mnameSpace), sumPods(n, mnameSpace), cpu, memory},
		})
	}
}

func sumCpu(space string, spacePods map[string]string) (scpu, smemory string) {
	//var services = make(map[string]string)
	var cpu, memory int
	for _, spacePods := range spacePods {
		v := strings.Split(spacePods, "#")
		nspace := v[0]
		status := strings.ToLower(v[3])
		if nspace == space && status == "running" {
			npod := strings.TrimSpace(v[1])
			runCMD(nspace, npod)
			// execut kubectl get pod
			cpu = cpu + 1
			memory = memory + 1
		}
	}

	scpu = strconv.Itoa(cpu)
	scpu = fmts.ConcatStr(scpu, "m")

	smemory = strconv.Itoa(memory)
	smemory = fmts.ConcatStr(smemory, "Mi")
	return
}

func runCMD(namespace, pod string) (sout string) {
	//var stdout, stderr bytes.Buffer
	command := fmts.ConcatStr("kubectl top pod ", pod, " --use-protocol-buffers -n ", namespace)
	out, err := exec.Command(command).Output()
	if err != nil {
		//fmt.Printf("error::::::", err.Error())
		return
	}
	sout = string(out)
	return
	//println("stdout.String():", string(out))
}

func sumService(space string, spacePods map[string]string) (i int) {
	var services = make(map[string]string)
	for _, spacePods := range spacePods {
		v := strings.Split(spacePods, "#")
		nspace := v[0]
		status := strings.ToLower(v[3])
		if nspace == space && status == "running" {
			i++
			npod := strings.TrimSpace(v[1])
			re := regexp.MustCompile(`([a-zA-Z0-9_-]+)-v([0-9]+)`)
			match := re.FindStringSubmatch(npod)
			if len(match) > 0 {
				//println("pod:", match[0])
				//services = append(services, match[0])
				services[match[0]] = match[0]
			}
		}
	}
	i = len(services)
	return
}

func sumPods(space string, spacePods map[string]string) (i int) {
	for _, spacePods := range spacePods {
		v := strings.Split(spacePods, "#")
		nspace := v[0]
		status := strings.ToLower(v[3])
		if nspace == space && status == "running" {
			i++
		}
	}
	return
}

func NamePod(t table.Writer, nspaceTotal, spacePods map[string]string) {
	for _, space := range nspaceTotal {
		var j int
		for _, spacePods := range spacePods {
			v := strings.Split(spacePods, "#")
			nspace := v[0]
			if nspace == space {
				j++
				npod := v[1]
				t.AppendRows([]table.Row{
					{j, nspace, npod},
				})
			}
		}
	}
}
