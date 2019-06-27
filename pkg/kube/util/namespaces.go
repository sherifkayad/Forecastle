package util

import (
	"github.com/stakater/Forecastle/pkg/config"
	"github.com/stakater/Forecastle/pkg/kube"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func PopulateNamespaceList(namespaceSelector config.NamespaceSelector) ([]string, error) {
	if namespaceSelector.Any {
		return []string{metav1.NamespaceAll}, nil
	}

	namespaces := []string{}

	if namespaceSelector.LabelSelector != nil && (len(namespaceSelector.LabelSelector.MatchLabels) != 0 || len(namespaceSelector.LabelSelector.MatchExpressions) != 0) {
		var labelsMap map[string]string
		var err error
		if labelsMap, err = labelSelectorAsMap(namespaceSelector.LabelSelector); err != nil {
			return nil, err
		}

		set := labels.Set(labelsMap)
		nsOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
		nsList, err := kube.GetClient().CoreV1().Namespaces().List(nsOptions)
		if err != nil {
			return nil, err
		}

		for _, ns := range nsList.Items {
			namespaces = append(namespaces, ns.Name)
		}
	}

	return removeDuplicates(append(namespaces, namespaceSelector.MatchNames...)), nil
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}