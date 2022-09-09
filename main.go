package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jeffotoni/k8slog/sdk/config"
	"gitlab.engdb.com.br/gmid/sdk/core/env"
	"gitlab.engdb.com.br/gmid/sdk/core/fmts"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	KUBE_CONFIG = env.GetString("KUBE_CONFIG", "~/.kube/config")
	SHOW_TABLE  = env.GetBool("SHOW_TABLE", true)

	c = config.Config()
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", KUBE_CONFIG)
	if err != nil {
		log.Println(err)
		return
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return
	}

	var totalItem = make(map[string]string)
	var tService int
	var tPods int
	var tCpus int64
	var tMemory int64

	for _, space := range c.Cluster.NameSpace {

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

				mcpu += cpuQuantity
				mmemory += memQuantity

				re := regexp.MustCompile(`([a-zA-Z0-9_-]+)-v([0-9]+)`)
				match := re.FindStringSubmatch(podMetric.ObjectMeta.Name)
				if len(match) > 0 {
					mservice[match[0]] = match[0]
				}
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
	if SHOW_TABLE {
		t.SetStyle(table.StyleColoredBright)
	}
	t.Render()

}
