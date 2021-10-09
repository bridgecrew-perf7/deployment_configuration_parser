package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type KubeDeployConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Labels struct {
			App string `yaml:"app"`
		} `yaml:"labels"`
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Strategy struct {
			Type string `yaml:"type"`
		} `yaml:"strategy"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `yaml:"app"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers         []map[string]interface{}
				ServiceAccountName string `yaml:"serviceAccountName"`
				Volumes            []map[string]interface{}
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

func main() {

	var resultDeploymentConfiguration map[string]interface{}
	var K KubeDeployConfig

	jsonFile, err := ioutil.ReadFile("deploymentconfig.json")

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	yamlFile, err := ioutil.ReadFile("deploy.yaml")

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	json.Unmarshal([]byte(jsonFile), &resultDeploymentConfiguration)
	yaml.Unmarshal([]byte(yamlFile), &K)

	conConf := resultDeploymentConfiguration["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
	volConf := resultDeploymentConfiguration["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["volumes"]
	metName := resultDeploymentConfiguration["metadata"].(map[string]interface{})["name"]

	mapstructure.Decode(volConf, &K.Spec.Template.Spec.Volumes)
	mapstructure.Decode(conConf, &K.Spec.Template.Spec.Containers)
	mapstructure.Decode(metName, &K.Metadata.Name)
	mapstructure.Decode(metName, &K.Metadata.Labels.App)
	mapstructure.Decode(metName, &K.Spec.Selector.MatchLabels.App)
	mapstructure.Decode(metName, &K.Spec.Template.Metadata.Labels.App)

	out, err := yaml.Marshal(K)
	f, err := os.Create(K.Metadata.Name + ".yaml")
	_, err2 := f.WriteString(string(out))
	if err2 != nil {
		log.Fatalf("err: %v", err)
	}

}
