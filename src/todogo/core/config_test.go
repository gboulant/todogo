package core

import "testing"

func TestConfig(t *testing.T) {
	config := Config{
		ContextName: "toto",
		ContextList: ContextArray{
			Context{DirPath: "/tmp/toto", Name: "toto"},
			Context{DirPath: "/tmp/tutu", Name: "tutu"},
			Context{DirPath: "/tmp/titi", Name: "titi"},
			Context{DirPath: "/tmp/tata", Name: "tata"},
		},
		Parameters: Parameters{
			DefaultCommand: "board",
		},
	}

	config.Save("/tmp/toto.json")

	var otherConfig Config
	otherConfig.Load("/tmp/toto.json")
	dirpath := otherConfig.ContextList[1].DirPath
	if dirpath != "/tmp/tutu" {
		t.Errorf("dirpath is %s (should be %s)", dirpath, "/tmp/tutu")
	}
	contextName := otherConfig.ContextName
	if contextName != "toto" {
		t.Errorf("contextName is %s (should be %s)", contextName, "toto")
	}
}

func TestConfigHandler(t *testing.T) {
	var handler ConfigHandler
	err := handler.Load()
	if err != nil {
		t.Error(err)
	}

	context := Context{DirPath: "/tmp/toto", Name: "toto"}
	handler.AddContext(&context)

	pcontext, err := handler.GetContext("toto")
	if err != nil {
		t.Error(err)
	}

	pcontext.DirPath = "/tmp/rototo"
	pcontext, _ = handler.GetContext("toto")
	if pcontext.DirPath != "/tmp/rototo" {
		t.Errorf("dirpath is %s (should be %s)", pcontext.DirPath, "/tmp/rototo")
	}
}
