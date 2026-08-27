package initiator

import (
	"aidoc/internal/entity"
	"aidoc/platform/ocr/postprocessing"

	"github.com/spf13/viper"
)

func initPostProcessing() postprocessing.PostProcessing {

	rules := []entity.CorrectionRule{}

	for _, rule := range viper.GetStringSlice("ocr.rules") {
		rules = append(rules, entity.CorrectionRule{
			Pattern: rule,
		})
	}


	cfg := entity.PostProcessConfig{
		CorrectionRules: rules,
		HeaderPatterns:  viper.GetStringSlice("ocr.headers"),
		FooterPatterns:  viper.GetStringSlice("ocr.footers"),
		Dictionary:      viper.GetStringSlice("ocr.dictionary"),
	}


	return postprocessing.NewPostProcessor(cfg)
}