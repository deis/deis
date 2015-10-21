// Package schema provides access to the fleet API.
//
// See http://github.com/coreos/fleet
//
// Usage example:
//
//   import "github.com/coreos/fleet/Godeps/_workspace/src/google.golang.org/api/schema/v1"
//   ...
//   schemaService, err := schema.New(oauthHttpClient)
package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/api/googleapi"
)

// Always reference these packages, just in case the auto-generated code
// below doesn't.
var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = googleapi.Version
var _ = errors.New
var _ = strings.Replace

const apiId = "fleet:v1"
const apiName = "schema"
const apiVersion = "v1"
const basePath = "$ENDPOINT/fleet/v1/"

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Machines = NewMachinesService(s)
	s.UnitState = NewUnitStateService(s)
	s.Units = NewUnitsService(s)
	return s, nil
}

type Service struct {
	client   *http.Client
	BasePath string // API endpoint base URL

	Machines *MachinesService

	UnitState *UnitStateService

	Units *UnitsService
}

func NewMachinesService(s *Service) *MachinesService {
	rs := &MachinesService{s: s}
	return rs
}

type MachinesService struct {
	s *Service
}

func NewUnitStateService(s *Service) *UnitStateService {
	rs := &UnitStateService{s: s}
	return rs
}

type UnitStateService struct {
	s *Service
}

func NewUnitsService(s *Service) *UnitsService {
	rs := &UnitsService{s: s}
	return rs
}

type UnitsService struct {
	s *Service
}

type Machine struct {
	Id string `json:"id,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`

	PrimaryIP string `json:"primaryIP,omitempty"`
}

type MachinePage struct {
	Machines []*Machine `json:"machines,omitempty"`

	NextPageToken string `json:"nextPageToken,omitempty"`
}

type Unit struct {
	CurrentState string `json:"currentState,omitempty"`

	DesiredState string `json:"desiredState,omitempty"`

	MachineID string `json:"machineID,omitempty"`

	Name string `json:"name,omitempty"`

	Options []*UnitOption `json:"options,omitempty"`
}

type UnitOption struct {
	Name string `json:"name,omitempty"`

	Section string `json:"section,omitempty"`

	Value string `json:"value,omitempty"`
}

type UnitPage struct {
	NextPageToken string `json:"nextPageToken,omitempty"`

	Units []*Unit `json:"units,omitempty"`
}

type UnitState struct {
	Hash string `json:"hash,omitempty"`

	MachineID string `json:"machineID,omitempty"`

	Name string `json:"name,omitempty"`

	SystemdActiveState string `json:"systemdActiveState,omitempty"`

	SystemdLoadState string `json:"systemdLoadState,omitempty"`

	SystemdSubState string `json:"systemdSubState,omitempty"`
}

type UnitStatePage struct {
	NextPageToken string `json:"nextPageToken,omitempty"`

	States []*UnitState `json:"states,omitempty"`
}

// method id "fleet.Machine.List":

type MachinesListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: Retrieve a page of Machine objects.
func (r *MachinesService) List() *MachinesListCall {
	c := &MachinesListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// NextPageToken sets the optional parameter "nextPageToken":
func (c *MachinesListCall) NextPageToken(nextPageToken string) *MachinesListCall {
	c.opt_["nextPageToken"] = nextPageToken
	return c
}

// Fields allows partial responses to be retrieved.
// See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *MachinesListCall) Fields(s ...googleapi.Field) *MachinesListCall {
	c.opt_["fields"] = googleapi.CombineFields(s)
	return c
}

func (c *MachinesListCall) Do() (*MachinePage, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["nextPageToken"]; ok {
		params.Set("nextPageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["fields"]; ok {
		params.Set("fields", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "machines")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	var ret *MachinePage
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Retrieve a page of Machine objects.",
	//   "httpMethod": "GET",
	//   "id": "fleet.Machine.List",
	//   "parameters": {
	//     "nextPageToken": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "machines",
	//   "response": {
	//     "$ref": "MachinePage"
	//   }
	// }

}

// method id "fleet.UnitState.List":

type UnitStateListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: Retrieve a page of UnitState objects.
func (r *UnitStateService) List() *UnitStateListCall {
	c := &UnitStateListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// MachineID sets the optional parameter "machineID":
func (c *UnitStateListCall) MachineID(machineID string) *UnitStateListCall {
	c.opt_["machineID"] = machineID
	return c
}

// NextPageToken sets the optional parameter "nextPageToken":
func (c *UnitStateListCall) NextPageToken(nextPageToken string) *UnitStateListCall {
	c.opt_["nextPageToken"] = nextPageToken
	return c
}

// UnitName sets the optional parameter "unitName":
func (c *UnitStateListCall) UnitName(unitName string) *UnitStateListCall {
	c.opt_["unitName"] = unitName
	return c
}

// Fields allows partial responses to be retrieved.
// See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UnitStateListCall) Fields(s ...googleapi.Field) *UnitStateListCall {
	c.opt_["fields"] = googleapi.CombineFields(s)
	return c
}

func (c *UnitStateListCall) Do() (*UnitStatePage, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["machineID"]; ok {
		params.Set("machineID", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["nextPageToken"]; ok {
		params.Set("nextPageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["unitName"]; ok {
		params.Set("unitName", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["fields"]; ok {
		params.Set("fields", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "state")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	var ret *UnitStatePage
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Retrieve a page of UnitState objects.",
	//   "httpMethod": "GET",
	//   "id": "fleet.UnitState.List",
	//   "parameters": {
	//     "machineID": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "nextPageToken": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "unitName": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "state",
	//   "response": {
	//     "$ref": "UnitStatePage"
	//   }
	// }

}

// method id "fleet.Unit.Delete":

type UnitsDeleteCall struct {
	s        *Service
	unitName string
	opt_     map[string]interface{}
}

// Delete: Delete the referenced Unit object.
func (r *UnitsService) Delete(unitName string) *UnitsDeleteCall {
	c := &UnitsDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.unitName = unitName
	return c
}

// Fields allows partial responses to be retrieved.
// See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UnitsDeleteCall) Fields(s ...googleapi.Field) *UnitsDeleteCall {
	c.opt_["fields"] = googleapi.CombineFields(s)
	return c
}

func (c *UnitsDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["fields"]; ok {
		params.Set("fields", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "units/{unitName}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	googleapi.Expand(req.URL, map[string]string{
		"unitName": c.unitName,
	})
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Delete the referenced Unit object.",
	//   "httpMethod": "DELETE",
	//   "id": "fleet.Unit.Delete",
	//   "parameterOrder": [
	//     "unitName"
	//   ],
	//   "parameters": {
	//     "unitName": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "units/{unitName}"
	// }

}

// method id "fleet.Unit.Get":

type UnitsGetCall struct {
	s        *Service
	unitName string
	opt_     map[string]interface{}
}

// Get: Retrieve a single Unit object.
func (r *UnitsService) Get(unitName string) *UnitsGetCall {
	c := &UnitsGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.unitName = unitName
	return c
}

// Fields allows partial responses to be retrieved.
// See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UnitsGetCall) Fields(s ...googleapi.Field) *UnitsGetCall {
	c.opt_["fields"] = googleapi.CombineFields(s)
	return c
}

func (c *UnitsGetCall) Do() (*Unit, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["fields"]; ok {
		params.Set("fields", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "units/{unitName}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	googleapi.Expand(req.URL, map[string]string{
		"unitName": c.unitName,
	})
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	var ret *Unit
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Retrieve a single Unit object.",
	//   "httpMethod": "GET",
	//   "id": "fleet.Unit.Get",
	//   "parameterOrder": [
	//     "unitName"
	//   ],
	//   "parameters": {
	//     "unitName": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "units/{unitName}",
	//   "response": {
	//     "$ref": "Unit"
	//   }
	// }

}

// method id "fleet.Unit.List":

type UnitsListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: Retrieve a page of Unit objects.
func (r *UnitsService) List() *UnitsListCall {
	c := &UnitsListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// NextPageToken sets the optional parameter "nextPageToken":
func (c *UnitsListCall) NextPageToken(nextPageToken string) *UnitsListCall {
	c.opt_["nextPageToken"] = nextPageToken
	return c
}

// Fields allows partial responses to be retrieved.
// See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UnitsListCall) Fields(s ...googleapi.Field) *UnitsListCall {
	c.opt_["fields"] = googleapi.CombineFields(s)
	return c
}

func (c *UnitsListCall) Do() (*UnitPage, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["nextPageToken"]; ok {
		params.Set("nextPageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["fields"]; ok {
		params.Set("fields", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "units")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	var ret *UnitPage
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Retrieve a page of Unit objects.",
	//   "httpMethod": "GET",
	//   "id": "fleet.Unit.List",
	//   "parameters": {
	//     "nextPageToken": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "units",
	//   "response": {
	//     "$ref": "UnitPage"
	//   }
	// }

}

// method id "fleet.Unit.Set":

type UnitsSetCall struct {
	s        *Service
	unitName string
	unit     *Unit
	opt_     map[string]interface{}
}

// Set: Create or update a Unit.
func (r *UnitsService) Set(unitName string, unit *Unit) *UnitsSetCall {
	c := &UnitsSetCall{s: r.s, opt_: make(map[string]interface{})}
	c.unitName = unitName
	c.unit = unit
	return c
}

// Fields allows partial responses to be retrieved.
// See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UnitsSetCall) Fields(s ...googleapi.Field) *UnitsSetCall {
	c.opt_["fields"] = googleapi.CombineFields(s)
	return c
}

func (c *UnitsSetCall) Do() error {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.unit)
	if err != nil {
		return err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["fields"]; ok {
		params.Set("fields", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "units/{unitName}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	googleapi.Expand(req.URL, map[string]string{
		"unitName": c.unitName,
	})
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Create or update a Unit.",
	//   "httpMethod": "PUT",
	//   "id": "fleet.Unit.Set",
	//   "parameterOrder": [
	//     "unitName"
	//   ],
	//   "parameters": {
	//     "unitName": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "units/{unitName}",
	//   "request": {
	//     "$ref": "Unit"
	//   }
	// }

}
