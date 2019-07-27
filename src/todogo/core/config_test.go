package core

import (
	"fmt"
	"testing"
)

func TestContexts(t *testing.T) {
	contexts := ContextArray{
		Context{DirPath: "/tmp/titi", Name: "titi"},
		Context{DirPath: "/tmp/toto", Name: "toto"},
		Context{DirPath: "/tmp/tutu", Name: "tutu"},
		Context{DirPath: "/tmp/tata", Name: "tata"},
	}

	pcontext := contexts.GetContext("tutu")
	if pcontext.DirPath != "/tmp/tutu" {
		msg := fmt.Sprintf("DirPath is %s (should be %s)", pcontext.DirPath, "/tmp/tutu")
		t.Error(msg)
	}

	pcontext.DirPath = "/my/path"
	pcontext = contexts.GetContext("tutu")
	if pcontext.DirPath != "/my/path" {
		msg := fmt.Sprintf("DirPath is %s (should be %s)", pcontext.DirPath, "/my/path")
		t.Error(msg)
	}

	contexts.Remove(contexts.IndexFromName("tutu"))
	pcontext = contexts.GetContext("tutu")
	if pcontext != nil {
		msg := fmt.Sprintf("Context %s should not exist", "tutu")
		t.Error(msg)
	}

	initlen := len(contexts)
	context := Context{DirPath: "/tmp/yeye", Name: "yeye"}
	contexts.Append(context)
	if len(contexts) != initlen+1 {
		msg := fmt.Sprintf("Contexts size is %d (should be %d)", len(contexts), initlen)
		t.Error(msg)
	}

	pcontext = contexts.GetContext("yeye")
	if pcontext.DirPath != "/tmp/yeye" {
		msg := fmt.Sprintf("DirPath is %s (should be %s)", pcontext.DirPath, "/tmp/yeye")
		t.Error(msg)
	}
}

func TestConfig(t *testing.T) {
	config := GetTestConfig()
	config.SaveTo("/tmp/toto.json")

	var otherConfig Config
	otherConfig.Load("/tmp/toto.json")
	dirpath := otherConfig.ContextList.GetContext("tutu").DirPath
	if dirpath != "/tmp/tutu" {
		t.Errorf("dirpath is %s (should be %s)", dirpath, "/tmp/tutu")
	}
	contextName := otherConfig.ContextName
	if contextName != "toto" {
		t.Errorf("contextName is %s (should be %s)", contextName, "toto")
	}
}

func TestConfigIO(t *testing.T) {

	config, err := GetConfig()
	if err != nil {
		t.Error(err)
	}

	err = config.ContextList.Append(Context{DirPath: "/tmp/yeye", Name: "yeye"})
	if err != nil {
		t.Error(err)
	}
	err = config.Save()
	if err != nil {
		t.Error(err)
	}

	idx := config.ContextList.IndexFromName("yeye")
	err = config.ContextList.Remove(idx)
	if err != nil {
		fmt.Println(err)
	}
	err = config.Save()
	if err != nil {
		t.Error(err)
	}

}
