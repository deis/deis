package omaha

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestOmahaRequestUpdateCheck(t *testing.T) {
	file, err := os.Open("../fixtures/update-engine/update/request.xml")
	if err != nil {
		t.Error(err)
	}
	fix, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	v := Request{}
	xml.Unmarshal(fix, &v)

	if v.Os.Version != "Indy" {
		t.Error("Unexpected version", v.Os.Version)
	}

	if v.Apps[0].Id != "{87efface-864d-49a5-9bb3-4b050a7c227a}" {
		t.Error("Expected an App Id")
	}

	if v.Apps[0].BootId != "{7D52A1CC-7066-40F0-91C7-7CB6A871BFDE}" {
		t.Error("Expected a Boot Id")
	}

	if v.Apps[0].MachineID != "{8BDE4C4D-9083-4D61-B41C-3253212C0C37}" {
		t.Error("Expected a MachineId")
	}

	if v.Apps[0].OEM != "ec3000" {
		t.Error("Expected an OEM")
	}

	if v.Apps[0].UpdateCheck == nil {
		t.Error("Expected an UpdateCheck")
	}

	if v.Apps[0].Version != "ForcedUpdate" {
		t.Error("Verison is ForcedUpdate")
	}

	if v.Apps[0].FromTrack != "developer-build" {
		t.Error("developer-build")
	}

	if v.Apps[0].Track != "dev-channel" {
		t.Error("dev-channel")
	}

	if v.Apps[0].Events[0].Type != "3" {
		t.Error("developer-build")
	}
}

func ExampleOmaha_NewResponse() {
	response := NewResponse("unit-test")
	app := response.AddApp("{52F1B9BC-D31A-4D86-9276-CBC256AADF9A}")
	app.Status = "ok"
	p := app.AddPing()
	p.Status = "ok"
	u := app.AddUpdateCheck()
	u.Status = "ok"
	u.AddUrl("http://localhost/updates")
	m := u.AddManifest("9999.0.0")
	m.AddPackage("+LXvjiaPkeYDLHoNKlf9qbJwvnk=", "update.gz", "67546213", true)
	a := m.AddAction("postinstall")
	a.ChromeOSVersion = "9999.0.0"
	a.Sha256 = "0VAlQW3RE99SGtSB5R4m08antAHO8XDoBMKDyxQT/Mg="
	a.NeedsAdmin = false
	a.IsDelta = true
	a.DisablePayloadBackoff = true

	if raw, err := xml.MarshalIndent(response, "", " "); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("%s%s\n", xml.Header, raw)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <response protocol="3.0" server="unit-test">
	//  <daystart elapsed_seconds="0"></daystart>
	//  <app appid="{52F1B9BC-D31A-4D86-9276-CBC256AADF9A}" status="ok">
	//   <ping status="ok"></ping>
	//   <updatecheck status="ok">
	//    <urls>
	//     <url codebase="http://localhost/updates"></url>
	//    </urls>
	//    <manifest version="9999.0.0">
	//     <packages>
	//      <package hash="+LXvjiaPkeYDLHoNKlf9qbJwvnk=" name="update.gz" size="67546213" required="true"></package>
	//     </packages>
	//     <actions>
	//      <action event="postinstall" ChromeOSVersion="9999.0.0" sha256="0VAlQW3RE99SGtSB5R4m08antAHO8XDoBMKDyxQT/Mg=" needsadmin="false" IsDelta="true" DisablePayloadBackoff="true"></action>
	//     </actions>
	//    </manifest>
	//   </updatecheck>
	//  </app>
	// </response>
}

func ExampleOmaha_NewRequest() {
	request := NewRequest("Indy", "Chrome OS", "ForcedUpdate_x86_64", "")
	app := request.AddApp("{27BD862E-8AE8-4886-A055-F7F1A6460627}", "1.0.0.0")
	app.AddUpdateCheck()

	event := app.AddEvent()
	event.Type = "1"
	event.Result = "0"

	if raw, err := xml.MarshalIndent(request, "", " "); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("%s%s\n", xml.Header, raw)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <request protocol="3.0">
	//  <os platform="Chrome OS" version="Indy" sp="ForcedUpdate_x86_64"></os>
	//  <app appid="{27BD862E-8AE8-4886-A055-F7F1A6460627}" version="1.0.0.0">
	//   <updatecheck></updatecheck>
	//   <event eventtype="1" eventresult="0"></event>
	//  </app>
	// </request>
}
