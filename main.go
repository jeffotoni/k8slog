package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeffotoni/k8slog/sdk/fmts"
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

	mn := map[string]string{
		//"kafka":        "kafka",
		//"kube-system":  "kube-system",
		"m-prd": "m-prd",
		//"rabbitmq-prd": "rabbitmq-prd",
		"s-prd": "s-prd",
		//"velero":       "velero",
		//"default": "default",
		//"oms":          "oms",
		"prd": "prd",
		//"e-prd":        "e-prd",
		//"gm-prd":       "gm-prd",
		//"f-prd":        "f-prd",
		//"log":          "log",
	}

	//var namespacetmp string
	//var i int = 0
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
			// fmt.Println("v:", pns)
			// mname[pname] = pname
			// fmt.Println("name:", pname)
			// i++
			// if i == 5 {
			// 	break
			// }

			//namespacetmp = pname
		}
	}

	f.Close()

	// data := [][]string{}
	for _, n := range nspaceTotal {
		// data = [][]string{
		// 	[]string{n, TotalPods(n, mnameSpace)},
		// }

		//fmt.Println("namespace:", is, " total pods:", len(v))
		fmt.Println("namespace:", n, " pods:", TotalPods(n, mnameSpace))
	}

	// table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	// table.SetFooter([]string{"", "", "Total NameSpace", TotalPods(n, mnameSpace)}) // Add Footer
	// table.SetBorder(false)                                                         // Set Border to false
	// table.AppendBulk(data)                                                         // Add Bulk Data
	// table.Render()

	fmt.Println("TOTAL NAMESPACE:", len(mnameSpace))
	//fmt.Println("NAME/PODS:", len(mname))

	// NAMESPACE   |       POD                        | READY   |  STATUS            | RESTARTS  | AGE
	// -----------+--------------------------+-------+----------------------------------------------------
	//  m-prd | r-product-details-v2-58b44bf4bb-v8tvc |  1/1     | Running           |  0        |  8h
	//  m-prd | r-product-details-v2-58b44bf4bb-v8tvc |  1/1     | Running           |  0        |  8h
	// -----------+--------------------------+-------+----------------------------------------------------
	// 										TOTAL NAMESPACE | 2025
	// 									  --------+----------

}

func TotalPods(space string, spacePods map[string]string) (i int) {
	for _, spacePods := range spacePods {
		v := strings.Split(spacePods, "#")
		for _, nspace := range v {
			if nspace == space {
				i++
			}
		}
	}
	return
}
