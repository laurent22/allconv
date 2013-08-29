package settings

import (
	"os"
	"os/user"
	"io/ioutil"
	"strings"
	"strconv"
	"time"
)

type Settings struct {
	applicationName_ string
	profileFolder_ string
	loaded_ bool
	dirty_ bool
	inner_ map[string](map[string]string)
	autosave_ bool
}

func New(applicationName string) *Settings {
	output := new(Settings)
	output.applicationName_ = applicationName
	output.loaded_ = false
	output.dirty_ = false
	output.autosave_ = true
	return output
}

func (this *Settings) profileFolder() (string, error) {
	if this.profileFolder_ != "" { return this.profileFolder_, nil }
	currentUser, err := user.Current()
	if err != nil { return "", err }
	folderPath := currentUser.HomeDir + string(os.PathSeparator) + ".config" + string(os.PathSeparator) + this.applicationName_
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil { return "", err }
	this.profileFolder_ = folderPath
	return this.profileFolder_, nil
}

func (this *Settings) profileFile() (string, error) {
	folder, err := this.profileFolder()
	if err != nil { return "", err }
	return folder + string(os.PathSeparator) + "Settings.ini", nil
}

func (this *Settings) Load() error {
	if this.loaded_ { return nil }
	this.inner_ = make(map[string](map[string]string))
	
	profileFilePath, err := this.profileFile()
	if err != nil { return err }
	
	if _, err := os.Stat(profileFilePath); os.IsNotExist(err) {
		this.loaded_ = true
		return nil
	}
	
	content, err := ioutil.ReadFile(profileFilePath)
	if err != nil { return err }
	
	categoryName := ""
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" { continue }
		if line[0] == '[' && line[len(line) - 1] == ']' {
			categoryName = line[1:len(line) - 1]
			this.inner_[categoryName] = make(map[string]string)
		} else {
			equalPos := strings.Index(line, "=")
			if equalPos == -1 { continue }
			pName := line[0:equalPos]
			pValue := line[equalPos+1:]
			this.inner_[categoryName][pName] = pValue
		}
	}
	
	this.loaded_ = true
	this.dirty_ = false
	
	return nil
}

func (this *Settings) Save() error {
	if !this.loaded_ || !this.dirty_ { return nil }
	
	profileFilePath, err := this.profileFile()
	if err != nil { return err }
	
	s := ""
	for category, properties := range this.inner_ {
		s += "[" + category + "]\n"
		for name, value := range properties {
			s += name + "=" + value + "\n"
		}
	}
	
    err = ioutil.WriteFile(profileFilePath, []byte(s), os.ModePerm)
    if err != nil { return err }
	
	this.dirty_ = false
	return nil
}

func (this *Settings) SetValue(category string, name string, value string) error {
	this.Load()
	_, exists := this.inner_[category]
	if !exists {
		this.inner_[category] = make(map[string]string)
	} else {
		current, _ := this.inner_[category][name]
		if current == value { return nil }
	}
	this.inner_[category][name] = value
	this.dirty_ = true
	return this.Save()
}

func (this *Settings) Value(category string, name string, defaultValue string) string {
	this.Load()
	cat, exists := this.inner_[category]
	if !exists { return defaultValue }
	output, exists := cat[name]
	if !exists { return defaultValue }
	return output
}

func (this *Settings) ValueFloat64(category string, name string, defaultValue float64) float64 {
	s := this.Value(category, name, "")
	if s == "" { return defaultValue }
	output, _ := strconv.ParseFloat(s, 64)
	return output
}

func (this *Settings) SetValueFloat64(category string, name string, value float64) error {
	return this.SetValue(category, name, strconv.FormatFloat(value, 'f', 8, 64))
}

func (this *Settings) ValueFloat32(category string, name string, defaultValue float32) float32 {
	s := this.Value(category, name, "")
	if s == "" { return defaultValue }
	output, _ := strconv.ParseFloat(s, 32)
	return float32(output)
}

func (this *Settings) SetValueFloat32(category string, name string, value float32) error {
	return this.SetValue(category, name, strconv.FormatFloat(float64(value), 'f', 8, 32))
}

func (this *Settings) ValueInt(category string, name string, defaultValue int) int {
	s := this.Value(category, name, "")
	if s == "" { return defaultValue }
	output, _ := strconv.Atoi(s)
	return output
}

func (this *Settings) SetValueInt(category string, name string, value int) error {
	return this.SetValue(category, name, strconv.Itoa(value))
}

func (this *Settings) ValueTime(category string, name string, defaultValue time.Time) time.Time {
	s := this.Value(category, name, "")
	if s == "" { return defaultValue }
	output, _ := time.Parse(time.RFC3339, s)
	return output
}

func (this *Settings) SetValueTime(category string, name string, value time.Time) error {
	return this.SetValue(category, name, value.Format(time.RFC3339))
}