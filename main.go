package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"gitlab.engdb.com.br/gmid/sdk/core/fmts"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	config, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBE_CONFIG"))
	mc, _ := metrics.NewForConfig(config)

	var nameSpace = []string{
		"f-prd",
		"e-prd",
		"m-prd",
		"s-prd",
		"gm-prd",
		"prd",
	}

	var totalItem = make(map[string]string)
	var tService int
	var tPods int
	var tCpus int64
	var tMemory int64

	for _, space := range nameSpace {

		var mservice = make(map[string]string)
		var mcpu int64
		var mmemory int64

		podMetrics, _ := mc.MetricsV1beta1().PodMetricses(space).List(context.TODO(), v1.ListOptions{})
		for _, podMetric := range podMetrics.Items {

			podContainers := podMetric.Containers
			for _, container := range podContainers {
				cpuQuantity := container.Usage.Cpu().MilliValue()
				// if !ok {
				// 	continue
				// }
				memQuantity, ok := container.Usage.Memory().AsInt64()
				if !ok {
					continue
				}

				//memQuantity = memQuantity / 1024 / 1024
				// msg := fmt.Sprintf("NameSpace: %s - Container Name: %s \n CPU usage: %d \n Memory usage: %d Mi",
				// 	podMetric.ObjectMeta.Namespace,
				// 	podMetric.ObjectMeta.Name,
				// 	cpuQuantity,
				// 	memQuantity)

				mcpu += cpuQuantity
				mmemory += memQuantity

				re := regexp.MustCompile(`([a-zA-Z0-9_-]+)-v([0-9]+)`)
				match := re.FindStringSubmatch(podMetric.ObjectMeta.Name)
				if len(match) > 0 {
					mservice[match[0]] = match[0]
				}

				// fmt.Println(msg)
			}

			mmemoryon := mmemory / 1024 / 1024
			strTotal := fmt.Sprintf("%s#%d#%d#%d#%d", podMetric.Namespace, len(podMetrics.Items), len(mservice), mcpu, mmemoryon)
			totalItem[podMetric.Namespace] = strTotal
		}
	}

	t := table.NewWriter()
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "NAMESPACE", "SERVICE", "PODS", "CPU", "MEM"})

	var i int
	for _, v := range totalItem {
		fields := strings.Split(v, "#")
		if len(fields) == 5 {
			i++
			namespace := fields[0]
			totalPods := fields[1]
			totalService := fields[2]
			totalCpu := fmts.ConcatStr(fields[3], "m")
			totalMemory := fmts.ConcatStr(fields[4], "Mi")

			// convert
			ts, _ := strconv.Atoi(totalService)
			tp, _ := strconv.Atoi(totalPods)
			tc, _ := strconv.ParseInt(fields[3], 10, 64)
			tm, _ := strconv.ParseInt(fields[4], 10, 64)
			tService += ts
			tPods += tp
			tCpus += tc
			tMemory += tm

			t.AppendRows([]table.Row{
				{i, namespace, totalService, totalPods, totalCpu, totalMemory},
			})

		}
	}

	tCpusStr := strconv.FormatInt(tCpus, 10)
	tCpusStr = fmts.ConcatStr(tCpusStr, "m")

	tMemoryStr := strconv.FormatInt(tMemory, 10)
	tMemoryStr = fmts.ConcatStr(tMemoryStr, "Mi")

	//t.AppendSeparator()
	t.AppendFooter(table.Row{"", "TOTAL", tService, tPods, tCpusStr, tMemoryStr})
	t.SetStyle(table.StyleColoredBright)
	t.Render()

}
