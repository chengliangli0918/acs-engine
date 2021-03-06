package acsengine

import (
	"path"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/leonelquinteros/gotext"
)

func TestAssignKubernetesParameters(t *testing.T) {
	// Initialize locale for translation
	locale := gotext.NewLocale(path.Join("..", "..", "translations"), "en_US")
	i18n.Initialize(locale)

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: locale,
		},
	}
	// iterate the test data directory
	apiModelTestFiles := &[]APIModelTestFile{}
	if e := IterateTestFilesDirectory(TestDataDir, apiModelTestFiles); e != nil {
		t.Error(e.Error())
		return
	}

	for _, tuple := range *apiModelTestFiles {
		containerService, _, err := apiloader.LoadContainerServiceFromFile(tuple.APIModelFilename, true, false, nil)
		if err != nil {
			t.Errorf("Loading file %s got error: %s", tuple.APIModelFilename, err.Error())
			continue
		}

		parametersMap := paramsMap{}
		containerService.Location = "eatsus"
		cloudSpecConfig := containerService.GetCloudSpecConfig()
		assignKubernetesParameters(containerService.Properties, parametersMap, cloudSpecConfig, DefaultGeneratorCode)
		for k, v := range parametersMap {
			switch val := v.(paramsMap)["value"].(type) {
			case *bool:
				t.Errorf("got a pointer to bool in paramsMap value, this is dangerous!: %s: %v", k, val)
			}
		}
	}
}

func isKeyPresent(key string, paramMap map[string]interface{}) bool {
	_, ok := paramMap[key]
	return ok
}
