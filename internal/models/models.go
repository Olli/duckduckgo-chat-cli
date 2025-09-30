package models

import (
	"os"
	"reflect"
	"strings"
	//"fmt"

	"duckduckgo-chat-cli/internal/ui"

	"github.com/AlecAivazis/survey/v2"
)

const (
	StatusURL     = "https://duckduckgo.com/duckchat/v1/status"
	ChatURL       = "https://duckduckgo.com/duckchat/v1/chat"
	StatusHeaders = "1"
)

type (
	Model      string
	ModelAlias string
)

type ModelStruct struct {
	modelName        Model
	modelAlias       ModelAlias
	modelShortName   string
	modelDisplayName string
}

type ModelSlice []ModelStruct

const (
	GPT4Mini Model = "gpt-4o-mini"
	GPT5Mini Model = "gpt-5-mini"
	GPTOSS   Model = "openai/gpt-oss-120b"
	Claude3  Model = "claude-3-haiku-20240307"
	Llama    Model = "meta-llama/Llama-3.3-70B-Instruct-Turbo"
	Mixtral  Model = "mistralai/Mistral-Small-24B-Instruct-2501"

	GPT4MiniAlias ModelAlias = "gpt-4o-mini"
	GPT5MiniAlias ModelAlias = "gpt-5-mini"
	GPTOSSAlias   ModelAlias = "gpt-oss"
	Claude3Alias  ModelAlias = "claude-3-haiku"
	LlamaAlias    ModelAlias = "llama"
	MixtralAlias  ModelAlias = "mixtral"
)

var Models ModelSlice = []ModelStruct{
	{
		modelName:        GPT4Mini,
		modelAlias:       GPT4MiniAlias,
		modelShortName:   string(GPT4MiniAlias),
		modelDisplayName: "GPT-4o-mini",
	},
	{
		modelName: GPT5Mini, 
		modelAlias: GPT5MiniAlias,
		modelShortName: string(GPT5MiniAlias),
		modelDisplayName: "GPT-5-mini",

	},
	{
		modelName: GPTOSS, 
		modelAlias: GPTOSSAlias, 
		modelShortName: string(GPTOSSAlias),
		modelDisplayName: "GPT-OSS",
	},
	{
		modelName: Claude3, 
		modelAlias: Claude3Alias, 
		modelShortName: string(Claude3Alias),
		modelDisplayName: "Claude-3-haiku",
	},
	{
		modelName: Llama, 
		modelAlias: LlamaAlias, 
		modelShortName: string(LlamaAlias),
		modelDisplayName: "Llama 3.3",
	},
	{
		modelName: Mixtral, 
		modelAlias: MixtralAlias, 
		modelShortName: string(MixtralAlias),
		modelDisplayName: "Mistral Small 3",
	},
}

var modelMap = map[ModelAlias]Model{
	GPT4MiniAlias: GPT4Mini,
	GPT5MiniAlias: GPT5Mini,
	Claude3Alias:  Claude3,
	LlamaAlias:    Llama,
	MixtralAlias:  Mixtral,
}

var modelDisplayMap = map[Model]string{
	GPT4Mini: "GPT-4o-mini",
	GPT5Mini: "GPT-5-mini",
	GPTOSS:   "GPT-OSS",
	Claude3:  "Claude-3-haiku",
	Llama:    "Llama 3.3",
	Mixtral:  "Mistral Small 3",
}

var DefaultModel = Models.defaultModel()

func (m ModelSlice) defaultModel() Model {
	model, _ :=  m.GetModel(GPT4Mini)
	return model
}

func (m ModelSlice) defaultModelStruct() ModelStruct {
	model, _ :=  m.getModelStruct(GPT4Mini)
	return model
}


// GetDisplayNames is a getter for display names of models 
func (m ModelSlice) GetDisplayNames() []string {
	displayNames  := make([]string,0,len(m)) 
	for _,actualModel := range m { 
		displayNames = append(displayNames,actualModel.modelDisplayName)	
	}
	return displayNames 
}

func (m ModelSlice) GetShortName(searchedModel Model) string {
	model, found := m.getModelStruct(searchedModel)
	if found {
		if model.modelShortName != ""{
			return model.modelShortName 
		} else {
			return model.modelDisplayName
		}
	}
	return "model not found"
}

func (m ModelSlice) GetModelFromAlias(alias string) Model {
	for _,actualModel := range m {
		if strings.EqualFold(string(actualModel.modelAlias), alias ){
			return actualModel.modelName
		}
		
	}
	return DefaultModel 
}

// GetModel gets the modelName from the ModelStruct of the found model
func (m ModelSlice) GetModel(searchedModel Model) (foundModel Model, ok bool) {
	actualModel, ok := m.getModelStruct(searchedModel) 
	if ok {
		return actualModel.modelName, ok
	}
	return
}

func (m ModelSlice) getModelStruct(searchedModel Model) (foundModel ModelStruct, ok bool) {
	for _,actualModel := range m {
		if actualModel.modelName == searchedModel {
			return actualModel, true
		}
		
	}
	return 
}


func (m ModelSlice) getModelByDisplayName(searchedDisplayName string) (foundModel Model, ok bool) {
	for _, actualModel := range m {
		if strings.EqualFold(strings.ToLower(actualModel.modelDisplayName), searchedDisplayName) {
			return actualModel.modelName, true
		}
	}
	return 
}




func GetModelFromAlias(alias string) Model {
	if model, ok := modelMap[ModelAlias(alias)]; ok {
		return model
	}
	return GPT4Mini // default model
}

func HandleModelChange(chat interface{}, modelArg string) Model {

	// If a model argument is provided, try to use it directly
	if modelArg != "" {
		for _, model := range Models {
			if strings.EqualFold(modelArg, string(model.modelAlias)) || strings.EqualFold(modelArg, string(model.modelName)) {
				return model.modelName
			}
		}
		ui.Errorln("Invalid model choice: %s", modelArg)
		return ""
	}

	modelOptions := Models.GetDisplayNames()
	modelOptions = append(modelOptions,"Cancel")

	currentModel := GetCurrentModel(chat)
	defaultModel , ok := Models.getModelStruct(currentModel) 

	if !ok {
		defaultModel =  Models.defaultModelStruct() // Fallback
	}

	var choice string
	prompt := &survey.Select{
		Message: "Choose a new model:",
		Options: modelOptions,
		Default: defaultModel.modelDisplayName,
	}

	
	err := survey.AskOne(prompt, &choice, survey.WithStdio(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		// Gracefully exit on any error, including Ctrl+C (Interrupt)
		return ""
	}

	// search for the the model
	choosenModel, modelFound := Models.getModelByDisplayName(strings.ToLower(choice)) 
	// use default if nothing found
	if !modelFound { 
		choosenModel = Models.defaultModel() 

	} else if strings.EqualFold(strings.ToLower(choice), "cancel") { 
		ui.Warningln("Model change canceled")
		return ""
 	}

 	return choosenModel

	// switch strings.ToLower(choice) {
	// case "gpt-4o-mini":
	// 	return GPT4MiniAlias
	// case "gpt-5-mini":
	// 	return GPT5MiniAlias
	// case "claude-3-haiku":
	// 	return Claude3Alias
	// case "llama 3.3":
	// 	return LlamaAlias
	// case "mistral small 3":
	// 	return MixtralAlias
	// case "cancel":
	// 	ui.Warningln("Model change canceled")
	// 	return ""
	// default:
	// 	return ""
	// }
	// return
}

func GetCurrentModel(chat interface{}) Model {
	v := reflect.ValueOf(chat)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if modelField := v.FieldByName("Model"); modelField.IsValid() {
		// return Model(modelField.String())
		return Models.GetModelFromAlias(modelField.String())
	}
	return ""
}
