package models

import (
	"os"
	"reflect"
	"strings"

	"duckduckgo-chat-cli/internal/ui"

	"github.com/AlecAivazis/survey/v2"
)

const (
	StatusURL        = "https://duckduckgo.com/duckchat/v1/status"
	ChatURL          = "https://duckduckgo.com/duckchat/v1/chat"
	StatusHeaders    = "1"
)

type Model string
type ModelAlias string

const (
	GPT4Mini Model = "gpt-4o-mini"
	GPT5Mini Model = "gpt-5-mini"
	Claude3  Model = "claude-3-haiku-20240307"
	Llama    Model = "meta-llama/Llama-3.3-70B-Instruct-Turbo"
	Mixtral  Model = "mistralai/Mistral-Small-24B-Instruct-2501"
	o4mini   Model = "o4-mini"

	GPT4MiniAlias ModelAlias = "gpt-4o-mini"
	GPT5MiniAlias ModelAlias = "gpt-5-mini"
	Claude3Alias  ModelAlias = "claude-3-haiku"
	LlamaAlias    ModelAlias = "llama"
	MixtralAlias  ModelAlias = "mixtral"
	o4miniAlias   ModelAlias = "o4mini"
)

var modelMap = map[ModelAlias]Model{
	GPT4MiniAlias: GPT4Mini,
	GPT5MiniAlias: GPT5Mini,
	Claude3Alias:  Claude3,
	LlamaAlias:    Llama,
	MixtralAlias:  Mixtral,
	o4miniAlias:   o4mini,
}

var modelDisplayMap = map[Model]string{
	GPT4Mini: "GPT-4o-mini",
	GPT5Mini: "GPT-5-mini",
	Claude3:  "Claude-3-haiku",
	Llama:    "Llama 3.3",
	Mixtral:  "Mistral Small 3",
	o4mini:   "o4-mini",
}

var DefaultModel = GPT4Mini

func GetModelFromAlias(alias string) Model {
	if model, ok := modelMap[ModelAlias(alias)]; ok {
		return model
	}
	return GPT4Mini // default model
}

func HandleModelChange(chat interface{}, modelArg string) ModelAlias {
	// If a model argument is provided, try to use it directly
	if modelArg != "" {
		for alias, model := range modelMap {
			if strings.EqualFold(modelArg, string(alias)) || strings.EqualFold(modelArg, string(model)) {
				return alias
			}
		}
		ui.Errorln("Invalid model choice: %s", modelArg)
		return ""
	}

	// Show an interactive menu if no argument is provided
	modelOptions := []string{
		"GPT-4o-mini",
		"GPT-5-mini",
		"Claude-3-haiku",
		"Llama 3.3",
		"Mistral Small 3",
		"o4-mini",
		"Cancel",
	}

	currentModel := GetCurrentModel(chat)
	defaultModel, ok := modelDisplayMap[currentModel]
	if !ok {
		defaultModel = modelDisplayMap[DefaultModel] // Fallback
	}

	var choice string
	prompt := &survey.Select{
		Message: "Choose a new model:",
		Options: modelOptions,
		Default: defaultModel,
	}
	err := survey.AskOne(prompt, &choice, survey.WithStdio(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		// Gracefully exit on any error, including Ctrl+C (Interrupt)
		return ""
	}

	switch strings.ToLower(choice) {
	case "gpt-4o-mini":
		return GPT4MiniAlias
	case "gpt-5-mini":
		return GPT5MiniAlias
	case "claude-3-haiku":
		return Claude3Alias
	case "llama 3.3":
		return LlamaAlias
	case "mistral small 3":
		return MixtralAlias
	case "o4mini":
		return o4miniAlias
	case "cancel":
		ui.Warningln("Model change canceled")
		return ""
	default:
		return ""
	}
}

func GetCurrentModel(chat interface{}) Model {
	v := reflect.ValueOf(chat)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if modelField := v.FieldByName("Model"); modelField.IsValid() {
		return Model(modelField.String())
	}
	return ""
}
