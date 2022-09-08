package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jeffotoni/k8slog/sdk/fmts"
)

// dasboard
// namespace | pods | cpu | mem
// namespace | servico | qnt pod | cpu | mem

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

	mn := map[string]string{
		//"kafka":        "kafka",
		//"kube-system":  "kube-system",
		"m-prd": "m-prd",
		//"rabbitmq-prd": "rabbitmq-prd",
		"s-prd": "s-prd",
		//"velero":       "velero",
		//"default": "default",
		//"oms":          "oms",
		"prd":    "prd",
		"e-prd":  "e-prd",
		"gm-prd": "gm-prd",
		"f-prd":  "f-prd",
		//"log":          "log",
	}

	//var namespacetmp string
	//var i int = 0 Running
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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "NameSpace", "POD"})

	for _, n := range nspaceTotal {
		NamePod(n, t, mnameSpace)
	}

	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", len(mnameSpace), ""})
	t.Render()
	println("")
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "NameSpace", "Pods", "File/CSV"})

	for i, n := range nspaceTotal {
		t.AppendRows([]table.Row{
			{i, n, TotalPods(n, mnameSpace), fmts.ConcatStr("File.", n, ".csv")},
		})
	}

	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", len(mnameSpace), ""})
	t.Render()
}

func TotalPods(space string, spacePods map[string]string) (i int) {
	for _, spacePods := range spacePods {
		v := strings.Split(spacePods, "#")
		nspace := v[0]
		if nspace == space {
			i++
		}
	}
	return
}

func NamePod(space string, t table.Writer, spacePods map[string]string) {
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
