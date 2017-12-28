package main

import (
	"k8s.io/client-go/util/homedir"
	"flag"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"strings"
	"os"
	"bufio"
	"log"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	var kubeconfig *string

	if home:= homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	var namespace = flag.String("ns","","Namespace")
	var labelPrefix = flag.String("lb","","Label prefix")

	var directory = flag.String("dir","","Directory to save files")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	check(err)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	check(err)

	pods, err := clientset.CoreV1().
			ConfigMaps(*namespace).
			List(metav1.ListOptions{})
	check(err)

	fmt.Println("Succes connect to cluster. Try to get information...")

	header, _ := template.New("header").Parse("{{- printf \"%-20.20s\" \"NAME\" }}{{- printf \"%-30.30s\" \"FILE TO TAKE\"  }}{{- printf \"%-20s\" \"FILE TO SAVE\"  }} \n")
	line, _ := template.New("header").Parse("{{- printf \"%-20.20s\" .Name }}{{- printf \"%-30.30s\" .FileToTake  }}{{- printf \"%-20s\" .FileToSave  }} \n")
	header.Execute(os.Stdout, "")

	for _, configMap := range pods.Items {

		for labelName, labelValue := range configMap.Labels{

			if strings.HasPrefix(labelName, *labelPrefix) {

				var filename = strings.Replace(labelName, *labelPrefix,"", -1)

				line.Execute(
					os.Stdout,
					map[string]string{
						"Name": configMap.Name,
						"FileToTake": filename,
						"FileToSave": labelValue,
					})


				if fileData, ok := configMap.Data[filename]; ok {
					fo, err := os.Create(*directory+labelValue)
					check(err)

					w := bufio.NewWriter(fo)

					_, err = w.Write([]byte(fileData));
					check(err)
					check(w.Flush())
					check(fo.Close())
				}else {
					log.New(os.Stderr, "", 0).Printf(
						"Cant find file name: %s in ConfigMap: %s\n",
						filename, configMap.Name,
					)
					log.New(os.Stdout, "", 0).Printf(
						"***************** Cant find file name: %s in ConfigMap: %s ******************\n",
						filename, configMap.Name,
					)
				}

			}
		}
	}


}