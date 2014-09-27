// Package update provides access to the .
//
// Usage example:
//
//   import "code.google.com/p/google-api-go-client/update/v1"
//   ...
//   updateService, err := update.New(oauthHttpClient)
package update

import (
	"bytes"
	"github.com/coreos/updateservicectl/Godeps/_workspace/src/code.google.com/p/google-api-go-client/googleapi"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

const apiId = "update:v1"
const apiName = "update"
const apiVersion = "v1"
const basePath = "http://internal/_ah/api/update/v1/"

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Admin = NewAdminService(s)
	s.App = NewAppService(s)
	s.Appversion = NewAppversionService(s)
	s.Channel = NewChannelService(s)
	s.Client = NewClientService(s)
	s.Clientupdate = NewClientupdateService(s)
	s.Group = NewGroupService(s)
	s.Upstream = NewUpstreamService(s)
	s.Util = NewUtilService(s)
	return s, nil
}

type Service struct {
	client   *http.Client
	BasePath string // API endpoint base URL

	Admin *AdminService

	App *AppService

	Appversion *AppversionService

	Channel *ChannelService

	Client *ClientService

	Clientupdate *ClientupdateService

	Group *GroupService

	Upstream *UpstreamService

	Util *UtilService
}

func NewAdminService(s *Service) *AdminService {
	rs := &AdminService{s: s}
	return rs
}

type AdminService struct {
	s *Service
}

func NewAppService(s *Service) *AppService {
	rs := &AppService{s: s}
	rs.Package = NewAppPackageService(s)
	return rs
}

type AppService struct {
	s *Service

	Package *AppPackageService
}

func NewAppPackageService(s *Service) *AppPackageService {
	rs := &AppPackageService{s: s}
	return rs
}

type AppPackageService struct {
	s *Service
}

func NewAppversionService(s *Service) *AppversionService {
	rs := &AppversionService{s: s}
	return rs
}

type AppversionService struct {
	s *Service
}

func NewChannelService(s *Service) *ChannelService {
	rs := &ChannelService{s: s}
	return rs
}

type ChannelService struct {
	s *Service
}

func NewClientService(s *Service) *ClientService {
	rs := &ClientService{s: s}
	return rs
}

type ClientService struct {
	s *Service
}

func NewClientupdateService(s *Service) *ClientupdateService {
	rs := &ClientupdateService{s: s}
	return rs
}

type ClientupdateService struct {
	s *Service
}

func NewGroupService(s *Service) *GroupService {
	rs := &GroupService{s: s}
	rs.Requests = NewGroupRequestsService(s)
	return rs
}

type GroupService struct {
	s *Service

	Requests *GroupRequestsService
}

func NewGroupRequestsService(s *Service) *GroupRequestsService {
	rs := &GroupRequestsService{s: s}
	rs.Events = NewGroupRequestsEventsService(s)
	rs.Versions = NewGroupRequestsVersionsService(s)
	return rs
}

type GroupRequestsService struct {
	s *Service

	Events *GroupRequestsEventsService

	Versions *GroupRequestsVersionsService
}

func NewGroupRequestsEventsService(s *Service) *GroupRequestsEventsService {
	rs := &GroupRequestsEventsService{s: s}
	return rs
}

type GroupRequestsEventsService struct {
	s *Service
}

func NewGroupRequestsVersionsService(s *Service) *GroupRequestsVersionsService {
	rs := &GroupRequestsVersionsService{s: s}
	return rs
}

type GroupRequestsVersionsService struct {
	s *Service
}

func NewUpstreamService(s *Service) *UpstreamService {
	rs := &UpstreamService{s: s}
	return rs
}

type UpstreamService struct {
	s *Service
}

func NewUtilService(s *Service) *UtilService {
	rs := &UtilService{s: s}
	return rs
}

type UtilService struct {
	s *Service
}

type AdminListUsersResp struct {
	Users []*AdminUser `json:"users,omitempty"`
}

type AdminUser struct {
	Token string `json:"token,omitempty"`

	User string `json:"user,omitempty"`
}

type AdminUserReq struct {
	UserName string `json:"userName,omitempty"`
}

type App struct {
	Description string `json:"description,omitempty"`

	Id string `json:"id,omitempty"`

	Label string `json:"label,omitempty"`
}

type AppChannel struct {
	AppId string `json:"appId,omitempty"`

	Label string `json:"label,omitempty"`

	Publish bool `json:"publish,omitempty"`

	Upstream string `json:"upstream,omitempty"`

	Version string `json:"version,omitempty"`
}

type AppInsertReq struct {
	Description string `json:"description,omitempty"`

	Id string `json:"id,omitempty"`

	Label string `json:"label,omitempty"`
}

type AppListResp struct {
	Items []*App `json:"items,omitempty"`
}

type AppUpdateReq struct {
	Description string `json:"description,omitempty"`

	Id string `json:"id,omitempty"`

	Label string `json:"label,omitempty"`
}

type AppVersionItem struct {
	AppId string `json:"appId,omitempty"`

	Count int64 `json:"count,omitempty"`

	GroupId string `json:"groupId,omitempty"`

	Version string `json:"version,omitempty"`
}

type AppVersionList struct {
	Items []*AppVersionItem `json:"items,omitempty"`
}

type ChannelListResp struct {
	Items []*AppChannel `json:"items,omitempty"`
}

type ChannelRequest struct {
	AppId string `json:"appId,omitempty"`

	Label string `json:"label,omitempty"`

	Publish bool `json:"publish,omitempty"`

	Version string `json:"version,omitempty"`
}

type ClientCountResp struct {
	Count int64 `json:"count,omitempty"`
}

type ClientHistoryItem struct {
	DateTime int64 `json:"dateTime,omitempty,string"`

	ErrorCode string `json:"errorCode,omitempty"`

	EventResult string `json:"eventResult,omitempty"`

	EventType string `json:"eventType,omitempty"`

	GroupId string `json:"groupId,omitempty"`

	Version string `json:"version,omitempty"`
}

type ClientHistoryResp struct {
	Items []*ClientHistoryItem `json:"items,omitempty"`
}

type ClientUpdate struct {
	AppId string `json:"appId,omitempty"`

	ClientId string `json:"clientId,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	EventResult string `json:"eventResult,omitempty"`

	EventType string `json:"eventType,omitempty"`

	GroupId string `json:"groupId,omitempty"`

	LastSeen string `json:"lastSeen,omitempty"`

	Oem string `json:"oem,omitempty"`

	Version string `json:"version,omitempty"`
}

type ClientUpdateList struct {
	Items []*ClientUpdate `json:"items,omitempty"`
}

type GenerateUuidResp struct {
	Uuid string `json:"uuid,omitempty"`
}

type Group struct {
	AppId string `json:"appId,omitempty"`

	ChannelId string `json:"channelId,omitempty"`

	Id string `json:"id,omitempty"`

	Label string `json:"label,omitempty"`

	UpdateCount int64 `json:"updateCount,omitempty"`

	UpdateInterval int64 `json:"updateInterval,omitempty"`

	UpdatePooling bool `json:"updatePooling,omitempty"`

	UpdatesPaused bool `json:"updatesPaused,omitempty"`
}

type GroupList struct {
	Items []*Group `json:"items,omitempty"`
}

type GroupRequestsItem struct {
	Result string `json:"result,omitempty"`

	Type string `json:"type,omitempty"`

	Values []*GroupRequestsValues `json:"values,omitempty"`

	Version string `json:"version,omitempty"`
}

type GroupRequestsRollup struct {
	Items []*GroupRequestsItem `json:"items,omitempty"`
}

type GroupRequestsValues struct {
	Count int64 `json:"count,omitempty,string"`

	Timestamp int64 `json:"timestamp,omitempty,string"`
}

type Package struct {
	AppId string `json:"appId,omitempty"`

	DateCreated string `json:"dateCreated,omitempty"`

	MetadataSignatureRsa string `json:"metadataSignatureRsa,omitempty"`

	MetadataSize string `json:"metadataSize,omitempty"`

	ReleaseNotes string `json:"releaseNotes,omitempty"`

	Required bool `json:"required,omitempty"`

	Sha1Sum string `json:"sha1Sum,omitempty"`

	Sha256Sum string `json:"sha256Sum,omitempty"`

	Size string `json:"size,omitempty"`

	Url string `json:"url,omitempty"`

	Version string `json:"version,omitempty"`
}

type PackageList struct {
	Items []*Package `json:"items,omitempty"`

	Total int64 `json:"total,omitempty"`
}

type PublicPackageItem struct {
	AppId string `json:"AppId,omitempty"`

	Packages []*Package `json:"packages,omitempty"`
}

type PublicPackageList struct {
	Items []*PublicPackageItem `json:"items,omitempty"`
}

type Upstream struct {
	Id int64 `json:"id,omitempty"`

	Label string `json:"label,omitempty"`

	Url string `json:"url,omitempty"`
}

type UpstreamListResp struct {
	Items []*Upstream `json:"items,omitempty"`
}

type UpstreamSyncResp struct {
	Detail string `json:"detail,omitempty"`

	Status string `json:"status,omitempty"`
}

// method id "update.admin.createUser":

type AdminCreateUserCall struct {
	s            *Service
	adminuserreq *AdminUserReq
	opt_         map[string]interface{}
}

// CreateUser: Create a new user.
func (r *AdminService) CreateUser(adminuserreq *AdminUserReq) *AdminCreateUserCall {
	c := &AdminCreateUserCall{s: r.s, opt_: make(map[string]interface{})}
	c.adminuserreq = adminuserreq
	return c
}

func (c *AdminCreateUserCall) Do() (*AdminUser, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.adminuserreq)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "admin/user")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(AdminUser)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Create a new user.",
	//   "httpMethod": "POST",
	//   "id": "update.admin.createUser",
	//   "path": "admin/user",
	//   "request": {
	//     "$ref": "AdminUserReq",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "AdminUser"
	//   }
	// }

}

// method id "update.admin.deleteUser":

type AdminDeleteUserCall struct {
	s        *Service
	userName string
	opt_     map[string]interface{}
}

// DeleteUser: Delete a user.
func (r *AdminService) DeleteUser(userName string) *AdminDeleteUserCall {
	c := &AdminDeleteUserCall{s: r.s, opt_: make(map[string]interface{})}
	c.userName = userName
	return c
}

func (c *AdminDeleteUserCall) Do() (*AdminUser, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "admin/user/{userName}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{userName}", url.QueryEscape(c.userName), 1)
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
	ret := new(AdminUser)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Delete a user.",
	//   "httpMethod": "DELETE",
	//   "id": "update.admin.deleteUser",
	//   "parameterOrder": [
	//     "userName"
	//   ],
	//   "parameters": {
	//     "userName": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "admin/user/{userName}",
	//   "response": {
	//     "$ref": "AdminUser"
	//   }
	// }

}

// method id "update.admin.genToken":

type AdminGenTokenCall struct {
	s            *Service
	userName     string
	adminuserreq *AdminUserReq
	opt_         map[string]interface{}
}

// GenToken: Generate a new token.
func (r *AdminService) GenToken(userName string, adminuserreq *AdminUserReq) *AdminGenTokenCall {
	c := &AdminGenTokenCall{s: r.s, opt_: make(map[string]interface{})}
	c.userName = userName
	c.adminuserreq = adminuserreq
	return c
}

func (c *AdminGenTokenCall) Do() (*AdminUser, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.adminuserreq)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "admin/user/{userName}/token/new")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{userName}", url.QueryEscape(c.userName), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(AdminUser)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Generate a new token.",
	//   "httpMethod": "PUT",
	//   "id": "update.admin.genToken",
	//   "parameterOrder": [
	//     "userName"
	//   ],
	//   "parameters": {
	//     "userName": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "admin/user/{userName}/token/new",
	//   "request": {
	//     "$ref": "AdminUserReq",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "AdminUser"
	//   }
	// }

}

// method id "update.admin.getUser":

type AdminGetUserCall struct {
	s        *Service
	userName string
	opt_     map[string]interface{}
}

// GetUser: Get a user.
func (r *AdminService) GetUser(userName string) *AdminGetUserCall {
	c := &AdminGetUserCall{s: r.s, opt_: make(map[string]interface{})}
	c.userName = userName
	return c
}

func (c *AdminGetUserCall) Do() (*AdminUser, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "admin/user/{userName}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{userName}", url.QueryEscape(c.userName), 1)
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
	ret := new(AdminUser)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Get a user.",
	//   "httpMethod": "GET",
	//   "id": "update.admin.getUser",
	//   "parameterOrder": [
	//     "userName"
	//   ],
	//   "parameters": {
	//     "userName": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "admin/user/{userName}",
	//   "response": {
	//     "$ref": "AdminUser"
	//   }
	// }

}

// method id "update.admin.listUsers":

type AdminListUsersCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// ListUsers: List Users.
func (r *AdminService) ListUsers() *AdminListUsersCall {
	c := &AdminListUsersCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AdminListUsersCall) Do() (*AdminListUsersResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "admin/user")
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
	ret := new(AdminListUsersResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List Users.",
	//   "httpMethod": "GET",
	//   "id": "update.admin.listUsers",
	//   "path": "admin/user",
	//   "response": {
	//     "$ref": "AdminListUsersResp"
	//   }
	// }

}

// method id "update.app.delete":

type AppDeleteCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

// Delete: Delete an application.
func (r *AppService) Delete(id string) *AppDeleteCall {
	c := &AppDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *AppDeleteCall) Do() (*App, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
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
	ret := new(App)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Delete an application.",
	//   "httpMethod": "DELETE",
	//   "id": "update.app.delete",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{id}",
	//   "response": {
	//     "$ref": "App"
	//   }
	// }

}

// method id "update.app.get":

type AppGetCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

// Get: Get an application.
func (r *AppService) Get(id string) *AppGetCall {
	c := &AppGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *AppGetCall) Do() (*App, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
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
	ret := new(App)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Get an application.",
	//   "httpMethod": "GET",
	//   "id": "update.app.get",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{id}",
	//   "response": {
	//     "$ref": "App"
	//   }
	// }

}

// method id "update.app.insert":

type AppInsertCall struct {
	s            *Service
	appinsertreq *AppInsertReq
	opt_         map[string]interface{}
}

// Insert: Insert an application.
func (r *AppService) Insert(appinsertreq *AppInsertReq) *AppInsertCall {
	c := &AppInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.appinsertreq = appinsertreq
	return c
}

func (c *AppInsertCall) Do() (*App, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.appinsertreq)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(App)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Insert an application.",
	//   "httpMethod": "POST",
	//   "id": "update.app.insert",
	//   "path": "apps",
	//   "request": {
	//     "$ref": "AppInsertReq",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "App"
	//   }
	// }

}

// method id "update.app.list":

type AppListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: List all application.
func (r *AppService) List() *AppListCall {
	c := &AppListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AppListCall) Do() (*AppListResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps")
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
	ret := new(AppListResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all application.",
	//   "httpMethod": "GET",
	//   "id": "update.app.list",
	//   "path": "apps",
	//   "response": {
	//     "$ref": "AppListResp"
	//   }
	// }

}

// method id "update.app.patch":

type AppPatchCall struct {
	s            *Service
	id           string
	appupdatereq *AppUpdateReq
	opt_         map[string]interface{}
}

// Patch: Update an application. This method supports patch semantics.
func (r *AppService) Patch(id string, appupdatereq *AppUpdateReq) *AppPatchCall {
	c := &AppPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	c.appupdatereq = appupdatereq
	return c
}

func (c *AppPatchCall) Do() (*App, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.appupdatereq)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(App)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Update an application. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "update.app.patch",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{id}",
	//   "request": {
	//     "$ref": "AppUpdateReq"
	//   },
	//   "response": {
	//     "$ref": "App"
	//   }
	// }

}

// method id "update.app.update":

type AppUpdateCall struct {
	s            *Service
	id           string
	appupdatereq *AppUpdateReq
	opt_         map[string]interface{}
}

// Update: Update an application.
func (r *AppService) Update(id string, appupdatereq *AppUpdateReq) *AppUpdateCall {
	c := &AppUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	c.appupdatereq = appupdatereq
	return c
}

func (c *AppUpdateCall) Do() (*App, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.appupdatereq)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(App)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Update an application.",
	//   "httpMethod": "PATCH",
	//   "id": "update.app.update",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{id}",
	//   "request": {
	//     "$ref": "AppUpdateReq",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "App"
	//   }
	// }

}

// method id "update.app.package.delete":

type AppPackageDeleteCall struct {
	s       *Service
	appId   string
	version string
	opt_    map[string]interface{}
}

// Delete: Delete an package.
func (r *AppPackageService) Delete(appId string, version string) *AppPackageDeleteCall {
	c := &AppPackageDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.version = version
	return c
}

// MetadataSignatureRsa sets the optional parameter
// "metadataSignatureRsa":
func (c *AppPackageDeleteCall) MetadataSignatureRsa(metadataSignatureRsa string) *AppPackageDeleteCall {
	c.opt_["metadataSignatureRsa"] = metadataSignatureRsa
	return c
}

// MetadataSize sets the optional parameter "metadataSize":
func (c *AppPackageDeleteCall) MetadataSize(metadataSize string) *AppPackageDeleteCall {
	c.opt_["metadataSize"] = metadataSize
	return c
}

// ReleaseNotes sets the optional parameter "releaseNotes":
func (c *AppPackageDeleteCall) ReleaseNotes(releaseNotes string) *AppPackageDeleteCall {
	c.opt_["releaseNotes"] = releaseNotes
	return c
}

// Required sets the optional parameter "required":
func (c *AppPackageDeleteCall) Required(required bool) *AppPackageDeleteCall {
	c.opt_["required"] = required
	return c
}

// Sha1Sum sets the optional parameter "sha1Sum":
func (c *AppPackageDeleteCall) Sha1Sum(sha1Sum string) *AppPackageDeleteCall {
	c.opt_["sha1Sum"] = sha1Sum
	return c
}

// Sha256Sum sets the optional parameter "sha256Sum":
func (c *AppPackageDeleteCall) Sha256Sum(sha256Sum string) *AppPackageDeleteCall {
	c.opt_["sha256Sum"] = sha256Sum
	return c
}

// Size sets the optional parameter "size":
func (c *AppPackageDeleteCall) Size(size string) *AppPackageDeleteCall {
	c.opt_["size"] = size
	return c
}

// Url sets the optional parameter "url":
func (c *AppPackageDeleteCall) Url(url string) *AppPackageDeleteCall {
	c.opt_["url"] = url
	return c
}

func (c *AppPackageDeleteCall) Do() (*Package, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["metadataSignatureRsa"]; ok {
		params.Set("metadataSignatureRsa", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["metadataSize"]; ok {
		params.Set("metadataSize", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["releaseNotes"]; ok {
		params.Set("releaseNotes", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["required"]; ok {
		params.Set("required", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["sha1Sum"]; ok {
		params.Set("sha1Sum", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["sha256Sum"]; ok {
		params.Set("sha256Sum", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["size"]; ok {
		params.Set("size", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["url"]; ok {
		params.Set("url", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/packages/{version}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{version}", url.QueryEscape(c.version), 1)
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
	ret := new(Package)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Delete an package.",
	//   "httpMethod": "DELETE",
	//   "id": "update.app.package.delete",
	//   "parameterOrder": [
	//     "appId",
	//     "version"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "metadataSignatureRsa": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "metadataSize": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "releaseNotes": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "required": {
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "sha1Sum": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "sha256Sum": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "size": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "url": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "version": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/packages/{version}",
	//   "response": {
	//     "$ref": "Package"
	//   }
	// }

}

// method id "update.app.package.insert":

type AppPackageInsertCall struct {
	s        *Service
	appId    string
	version  string
	package_ *Package
	opt_     map[string]interface{}
}

// Insert: Insert a new package version.
func (r *AppPackageService) Insert(appId string, version string, package_ *Package) *AppPackageInsertCall {
	c := &AppPackageInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.version = version
	c.package_ = package_
	return c
}

func (c *AppPackageInsertCall) Do() (*Package, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.package_)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/packages/{version}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{version}", url.QueryEscape(c.version), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Package)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Insert a new package version.",
	//   "httpMethod": "POST",
	//   "id": "update.app.package.insert",
	//   "parameterOrder": [
	//     "appId",
	//     "version"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "version": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/packages/{version}",
	//   "request": {
	//     "$ref": "Package",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "Package"
	//   }
	// }

}

// method id "update.app.package.list":

type AppPackageListCall struct {
	s     *Service
	appId string
	opt_  map[string]interface{}
}

// List: List all of the package versions.
func (r *AppPackageService) List(appId string) *AppPackageListCall {
	c := &AppPackageListCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	return c
}

// Limit sets the optional parameter "limit":
func (c *AppPackageListCall) Limit(limit int64) *AppPackageListCall {
	c.opt_["limit"] = limit
	return c
}

// Skip sets the optional parameter "skip":
func (c *AppPackageListCall) Skip(skip int64) *AppPackageListCall {
	c.opt_["skip"] = skip
	return c
}

// Version sets the optional parameter "version":
func (c *AppPackageListCall) Version(version string) *AppPackageListCall {
	c.opt_["version"] = version
	return c
}

func (c *AppPackageListCall) Do() (*PackageList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["limit"]; ok {
		params.Set("limit", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["skip"]; ok {
		params.Set("skip", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["version"]; ok {
		params.Set("version", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/packages")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
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
	ret := new(PackageList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all of the package versions.",
	//   "httpMethod": "GET",
	//   "id": "update.app.package.list",
	//   "parameterOrder": [
	//     "appId"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "limit": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "skip": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "version": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/packages",
	//   "response": {
	//     "$ref": "PackageList"
	//   }
	// }

}

// method id "update.app.package.publicList":

type AppPackagePublicListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// PublicList: List all of the publicly available published packages.
func (r *AppPackageService) PublicList() *AppPackagePublicListCall {
	c := &AppPackagePublicListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AppPackagePublicListCall) Do() (*PublicPackageList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "public/packages")
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
	ret := new(PublicPackageList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all of the publicly available published packages.",
	//   "httpMethod": "GET",
	//   "id": "update.app.package.publicList",
	//   "path": "public/packages",
	//   "response": {
	//     "$ref": "PublicPackageList"
	//   }
	// }

}

// method id "update.appversion.list":

type AppversionListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: List Client updates grouped by app/version.
func (r *AppversionService) List() *AppversionListCall {
	c := &AppversionListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// AppId sets the optional parameter "appId":
func (c *AppversionListCall) AppId(appId string) *AppversionListCall {
	c.opt_["appId"] = appId
	return c
}

// DateEnd sets the optional parameter "dateEnd":
func (c *AppversionListCall) DateEnd(dateEnd int64) *AppversionListCall {
	c.opt_["dateEnd"] = dateEnd
	return c
}

// DateStart sets the optional parameter "dateStart":
func (c *AppversionListCall) DateStart(dateStart int64) *AppversionListCall {
	c.opt_["dateStart"] = dateStart
	return c
}

// EventResult sets the optional parameter "eventResult":
func (c *AppversionListCall) EventResult(eventResult string) *AppversionListCall {
	c.opt_["eventResult"] = eventResult
	return c
}

// EventType sets the optional parameter "eventType":
func (c *AppversionListCall) EventType(eventType string) *AppversionListCall {
	c.opt_["eventType"] = eventType
	return c
}

// GroupId sets the optional parameter "groupId":
func (c *AppversionListCall) GroupId(groupId string) *AppversionListCall {
	c.opt_["groupId"] = groupId
	return c
}

// Oem sets the optional parameter "oem":
func (c *AppversionListCall) Oem(oem string) *AppversionListCall {
	c.opt_["oem"] = oem
	return c
}

// Version sets the optional parameter "version":
func (c *AppversionListCall) Version(version string) *AppversionListCall {
	c.opt_["version"] = version
	return c
}

func (c *AppversionListCall) Do() (*AppVersionList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["appId"]; ok {
		params.Set("appId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["dateEnd"]; ok {
		params.Set("dateEnd", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["dateStart"]; ok {
		params.Set("dateStart", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["eventResult"]; ok {
		params.Set("eventResult", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["eventType"]; ok {
		params.Set("eventType", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["groupId"]; ok {
		params.Set("groupId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["oem"]; ok {
		params.Set("oem", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["version"]; ok {
		params.Set("version", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "appversions")
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
	ret := new(AppVersionList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List Client updates grouped by app/version.",
	//   "httpMethod": "GET",
	//   "id": "update.appversion.list",
	//   "parameters": {
	//     "appId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dateEnd": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dateStart": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "eventResult": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "eventType": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "groupId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "oem": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "version": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "appversions",
	//   "response": {
	//     "$ref": "AppVersionList"
	//   }
	// }

}

// method id "update.channel.delete":

type ChannelDeleteCall struct {
	s     *Service
	appId string
	label string
	opt_  map[string]interface{}
}

// Delete: Delete a channel.
func (r *ChannelService) Delete(appId string, label string) *ChannelDeleteCall {
	c := &ChannelDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.label = label
	return c
}

// Publish sets the optional parameter "publish":
func (c *ChannelDeleteCall) Publish(publish bool) *ChannelDeleteCall {
	c.opt_["publish"] = publish
	return c
}

// Version sets the optional parameter "version":
func (c *ChannelDeleteCall) Version(version string) *ChannelDeleteCall {
	c.opt_["version"] = version
	return c
}

func (c *ChannelDeleteCall) Do() (*ChannelRequest, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["publish"]; ok {
		params.Set("publish", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["version"]; ok {
		params.Set("version", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/channels/{label}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{label}", url.QueryEscape(c.label), 1)
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
	ret := new(ChannelRequest)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Delete a channel.",
	//   "httpMethod": "DELETE",
	//   "id": "update.channel.delete",
	//   "parameterOrder": [
	//     "appId",
	//     "label"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "label": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "publish": {
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "version": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/channels/{label}",
	//   "response": {
	//     "$ref": "ChannelRequest"
	//   }
	// }

}

// method id "update.channel.insert":

type ChannelInsertCall struct {
	s              *Service
	appId          string
	channelrequest *ChannelRequest
	opt_           map[string]interface{}
}

// Insert: Insert a channel.
func (r *ChannelService) Insert(appId string, channelrequest *ChannelRequest) *ChannelInsertCall {
	c := &ChannelInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.channelrequest = channelrequest
	return c
}

func (c *ChannelInsertCall) Do() (*AppChannel, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.channelrequest)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/channels")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(AppChannel)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Insert a channel.",
	//   "httpMethod": "POST",
	//   "id": "update.channel.insert",
	//   "parameterOrder": [
	//     "appId"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/channels",
	//   "request": {
	//     "$ref": "ChannelRequest",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "AppChannel"
	//   }
	// }

}

// method id "update.channel.list":

type ChannelListCall struct {
	s     *Service
	appId string
	opt_  map[string]interface{}
}

// List: List channels.
func (r *ChannelService) List(appId string) *ChannelListCall {
	c := &ChannelListCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	return c
}

func (c *ChannelListCall) Do() (*ChannelListResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/channels")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
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
	ret := new(ChannelListResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List channels.",
	//   "httpMethod": "GET",
	//   "id": "update.channel.list",
	//   "parameterOrder": [
	//     "appId"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/channels",
	//   "response": {
	//     "$ref": "ChannelListResp"
	//   }
	// }

}

// method id "update.channel.publicList":

type ChannelPublicListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// PublicList: List all publicly available published channels.
func (r *ChannelService) PublicList() *ChannelPublicListCall {
	c := &ChannelPublicListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *ChannelPublicListCall) Do() (*ChannelListResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "public/channels")
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
	ret := new(ChannelListResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all publicly available published channels.",
	//   "httpMethod": "GET",
	//   "id": "update.channel.publicList",
	//   "path": "public/channels",
	//   "response": {
	//     "$ref": "ChannelListResp"
	//   }
	// }

}

// method id "update.channel.update":

type ChannelUpdateCall struct {
	s              *Service
	appId          string
	label          string
	channelrequest *ChannelRequest
	opt_           map[string]interface{}
}

// Update: Update a channel.
func (r *ChannelService) Update(appId string, label string, channelrequest *ChannelRequest) *ChannelUpdateCall {
	c := &ChannelUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.label = label
	c.channelrequest = channelrequest
	return c
}

func (c *ChannelUpdateCall) Do() (*AppChannel, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.channelrequest)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/channels/{label}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{label}", url.QueryEscape(c.label), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(AppChannel)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Update a channel.",
	//   "httpMethod": "PATCH",
	//   "id": "update.channel.update",
	//   "parameterOrder": [
	//     "appId",
	//     "label"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "label": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/channels/{label}",
	//   "request": {
	//     "$ref": "ChannelRequest",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "AppChannel"
	//   }
	// }

}

// method id "update.client.history":

type ClientHistoryCall struct {
	s        *Service
	clientId string
	opt_     map[string]interface{}
}

// History: Get the update history of a single client.
func (r *ClientService) History(clientId string) *ClientHistoryCall {
	c := &ClientHistoryCall{s: r.s, opt_: make(map[string]interface{})}
	c.clientId = clientId
	return c
}

func (c *ClientHistoryCall) Do() (*ClientHistoryResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	params.Set("clientId", fmt.Sprintf("%v", c.clientId))
	urls := googleapi.ResolveRelative(c.s.BasePath, "client/history")
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
	ret := new(ClientHistoryResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Get the update history of a single client.",
	//   "httpMethod": "GET",
	//   "id": "update.client.history",
	//   "parameterOrder": [
	//     "clientId"
	//   ],
	//   "parameters": {
	//     "clientId": {
	//       "location": "query",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "client/history",
	//   "response": {
	//     "$ref": "ClientHistoryResp"
	//   }
	// }

}

// method id "update.clientupdate.count":

type ClientupdateCountCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// Count: Get client count for criteria.
func (r *ClientupdateService) Count() *ClientupdateCountCall {
	c := &ClientupdateCountCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// AppId sets the optional parameter "appId":
func (c *ClientupdateCountCall) AppId(appId string) *ClientupdateCountCall {
	c.opt_["appId"] = appId
	return c
}

// DateEnd sets the optional parameter "dateEnd":
func (c *ClientupdateCountCall) DateEnd(dateEnd int64) *ClientupdateCountCall {
	c.opt_["dateEnd"] = dateEnd
	return c
}

// DateStart sets the optional parameter "dateStart":
func (c *ClientupdateCountCall) DateStart(dateStart int64) *ClientupdateCountCall {
	c.opt_["dateStart"] = dateStart
	return c
}

// EventResult sets the optional parameter "eventResult":
func (c *ClientupdateCountCall) EventResult(eventResult string) *ClientupdateCountCall {
	c.opt_["eventResult"] = eventResult
	return c
}

// EventType sets the optional parameter "eventType":
func (c *ClientupdateCountCall) EventType(eventType string) *ClientupdateCountCall {
	c.opt_["eventType"] = eventType
	return c
}

// GroupId sets the optional parameter "groupId":
func (c *ClientupdateCountCall) GroupId(groupId string) *ClientupdateCountCall {
	c.opt_["groupId"] = groupId
	return c
}

// Oem sets the optional parameter "oem":
func (c *ClientupdateCountCall) Oem(oem string) *ClientupdateCountCall {
	c.opt_["oem"] = oem
	return c
}

// Version sets the optional parameter "version":
func (c *ClientupdateCountCall) Version(version string) *ClientupdateCountCall {
	c.opt_["version"] = version
	return c
}

func (c *ClientupdateCountCall) Do() (*ClientCountResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["appId"]; ok {
		params.Set("appId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["dateEnd"]; ok {
		params.Set("dateEnd", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["dateStart"]; ok {
		params.Set("dateStart", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["eventResult"]; ok {
		params.Set("eventResult", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["eventType"]; ok {
		params.Set("eventType", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["groupId"]; ok {
		params.Set("groupId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["oem"]; ok {
		params.Set("oem", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["version"]; ok {
		params.Set("version", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "clientupdatecount")
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
	ret := new(ClientCountResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Get client count for criteria.",
	//   "httpMethod": "GET",
	//   "id": "update.clientupdate.count",
	//   "parameters": {
	//     "appId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dateEnd": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dateStart": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "eventResult": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "eventType": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "groupId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "oem": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "version": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "clientupdatecount",
	//   "response": {
	//     "$ref": "ClientCountResp"
	//   }
	// }

}

// method id "update.clientupdate.list":

type ClientupdateListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: List all client updates.
func (r *ClientupdateService) List() *ClientupdateListCall {
	c := &ClientupdateListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// AppId sets the optional parameter "appId":
func (c *ClientupdateListCall) AppId(appId string) *ClientupdateListCall {
	c.opt_["appId"] = appId
	return c
}

// ClientId sets the optional parameter "clientId":
func (c *ClientupdateListCall) ClientId(clientId string) *ClientupdateListCall {
	c.opt_["clientId"] = clientId
	return c
}

// DateEnd sets the optional parameter "dateEnd":
func (c *ClientupdateListCall) DateEnd(dateEnd int64) *ClientupdateListCall {
	c.opt_["dateEnd"] = dateEnd
	return c
}

// DateStart sets the optional parameter "dateStart":
func (c *ClientupdateListCall) DateStart(dateStart int64) *ClientupdateListCall {
	c.opt_["dateStart"] = dateStart
	return c
}

// EventResult sets the optional parameter "eventResult":
func (c *ClientupdateListCall) EventResult(eventResult string) *ClientupdateListCall {
	c.opt_["eventResult"] = eventResult
	return c
}

// EventType sets the optional parameter "eventType":
func (c *ClientupdateListCall) EventType(eventType string) *ClientupdateListCall {
	c.opt_["eventType"] = eventType
	return c
}

// GroupId sets the optional parameter "groupId":
func (c *ClientupdateListCall) GroupId(groupId string) *ClientupdateListCall {
	c.opt_["groupId"] = groupId
	return c
}

// Limit sets the optional parameter "limit":
func (c *ClientupdateListCall) Limit(limit int64) *ClientupdateListCall {
	c.opt_["limit"] = limit
	return c
}

// Oem sets the optional parameter "oem":
func (c *ClientupdateListCall) Oem(oem string) *ClientupdateListCall {
	c.opt_["oem"] = oem
	return c
}

// Skip sets the optional parameter "skip":
func (c *ClientupdateListCall) Skip(skip int64) *ClientupdateListCall {
	c.opt_["skip"] = skip
	return c
}

// Version sets the optional parameter "version":
func (c *ClientupdateListCall) Version(version string) *ClientupdateListCall {
	c.opt_["version"] = version
	return c
}

func (c *ClientupdateListCall) Do() (*ClientUpdateList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["appId"]; ok {
		params.Set("appId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["clientId"]; ok {
		params.Set("clientId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["dateEnd"]; ok {
		params.Set("dateEnd", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["dateStart"]; ok {
		params.Set("dateStart", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["eventResult"]; ok {
		params.Set("eventResult", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["eventType"]; ok {
		params.Set("eventType", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["groupId"]; ok {
		params.Set("groupId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["limit"]; ok {
		params.Set("limit", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["oem"]; ok {
		params.Set("oem", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["skip"]; ok {
		params.Set("skip", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["version"]; ok {
		params.Set("version", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "clientupdates")
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
	ret := new(ClientUpdateList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all client updates.",
	//   "httpMethod": "GET",
	//   "id": "update.clientupdate.list",
	//   "parameters": {
	//     "appId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "clientId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dateEnd": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "dateStart": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "eventResult": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "eventType": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "groupId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "limit": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "oem": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "skip": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "version": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "clientupdates",
	//   "response": {
	//     "$ref": "ClientUpdateList"
	//   }
	// }

}

// method id "update.group.delete":

type GroupDeleteCall struct {
	s     *Service
	appId string
	id    string
	opt_  map[string]interface{}
}

// Delete: Delete a group.
func (r *GroupService) Delete(appId string, id string) *GroupDeleteCall {
	c := &GroupDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.id = id
	return c
}

// ChannelId sets the optional parameter "channelId":
func (c *GroupDeleteCall) ChannelId(channelId string) *GroupDeleteCall {
	c.opt_["channelId"] = channelId
	return c
}

// Label sets the optional parameter "label":
func (c *GroupDeleteCall) Label(label string) *GroupDeleteCall {
	c.opt_["label"] = label
	return c
}

// UpdateCount sets the optional parameter "updateCount":
func (c *GroupDeleteCall) UpdateCount(updateCount int64) *GroupDeleteCall {
	c.opt_["updateCount"] = updateCount
	return c
}

// UpdateInterval sets the optional parameter "updateInterval":
func (c *GroupDeleteCall) UpdateInterval(updateInterval int64) *GroupDeleteCall {
	c.opt_["updateInterval"] = updateInterval
	return c
}

// UpdatePooling sets the optional parameter "updatePooling":
func (c *GroupDeleteCall) UpdatePooling(updatePooling bool) *GroupDeleteCall {
	c.opt_["updatePooling"] = updatePooling
	return c
}

// UpdatesPaused sets the optional parameter "updatesPaused":
func (c *GroupDeleteCall) UpdatesPaused(updatesPaused bool) *GroupDeleteCall {
	c.opt_["updatesPaused"] = updatesPaused
	return c
}

func (c *GroupDeleteCall) Do() (*Group, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["channelId"]; ok {
		params.Set("channelId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["label"]; ok {
		params.Set("label", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateCount"]; ok {
		params.Set("updateCount", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateInterval"]; ok {
		params.Set("updateInterval", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updatePooling"]; ok {
		params.Set("updatePooling", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updatesPaused"]; ok {
		params.Set("updatesPaused", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
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
	ret := new(Group)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Delete a group.",
	//   "httpMethod": "DELETE",
	//   "id": "update.group.delete",
	//   "parameterOrder": [
	//     "appId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "channelId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "label": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "updateCount": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "updateInterval": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "updatePooling": {
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "updatesPaused": {
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "apps/{appId}/groups/{id}",
	//   "response": {
	//     "$ref": "Group"
	//   }
	// }

}

// method id "update.group.get":

type GroupGetCall struct {
	s     *Service
	appId string
	id    string
	opt_  map[string]interface{}
}

// Get: Get a group.
func (r *GroupService) Get(appId string, id string) *GroupGetCall {
	c := &GroupGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.id = id
	return c
}

// ChannelId sets the optional parameter "channelId":
func (c *GroupGetCall) ChannelId(channelId string) *GroupGetCall {
	c.opt_["channelId"] = channelId
	return c
}

// Label sets the optional parameter "label":
func (c *GroupGetCall) Label(label string) *GroupGetCall {
	c.opt_["label"] = label
	return c
}

// UpdateCount sets the optional parameter "updateCount":
func (c *GroupGetCall) UpdateCount(updateCount int64) *GroupGetCall {
	c.opt_["updateCount"] = updateCount
	return c
}

// UpdateInterval sets the optional parameter "updateInterval":
func (c *GroupGetCall) UpdateInterval(updateInterval int64) *GroupGetCall {
	c.opt_["updateInterval"] = updateInterval
	return c
}

// UpdatePooling sets the optional parameter "updatePooling":
func (c *GroupGetCall) UpdatePooling(updatePooling bool) *GroupGetCall {
	c.opt_["updatePooling"] = updatePooling
	return c
}

// UpdatesPaused sets the optional parameter "updatesPaused":
func (c *GroupGetCall) UpdatesPaused(updatesPaused bool) *GroupGetCall {
	c.opt_["updatesPaused"] = updatesPaused
	return c
}

func (c *GroupGetCall) Do() (*Group, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["channelId"]; ok {
		params.Set("channelId", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["label"]; ok {
		params.Set("label", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateCount"]; ok {
		params.Set("updateCount", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateInterval"]; ok {
		params.Set("updateInterval", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updatePooling"]; ok {
		params.Set("updatePooling", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updatesPaused"]; ok {
		params.Set("updatesPaused", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
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
	ret := new(Group)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Get a group.",
	//   "httpMethod": "GET",
	//   "id": "update.group.get",
	//   "parameterOrder": [
	//     "appId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "channelId": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "label": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "updateCount": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "updateInterval": {
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "updatePooling": {
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "updatesPaused": {
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "apps/{appId}/groups/{id}",
	//   "response": {
	//     "$ref": "Group"
	//   }
	// }

}

// method id "update.group.insert":

type GroupInsertCall struct {
	s     *Service
	appId string
	group *Group
	opt_  map[string]interface{}
}

// Insert: Create a new group.
func (r *GroupService) Insert(appId string, group *Group) *GroupInsertCall {
	c := &GroupInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.group = group
	return c
}

func (c *GroupInsertCall) Do() (*Group, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.group)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Group)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Create a new group.",
	//   "httpMethod": "POST",
	//   "id": "update.group.insert",
	//   "parameterOrder": [
	//     "appId"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/groups",
	//   "request": {
	//     "$ref": "Group",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "Group"
	//   }
	// }

}

// method id "update.group.list":

type GroupListCall struct {
	s     *Service
	appId string
	opt_  map[string]interface{}
}

// List: List all of the groups.
func (r *GroupService) List(appId string) *GroupListCall {
	c := &GroupListCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	return c
}

// Limit sets the optional parameter "limit":
func (c *GroupListCall) Limit(limit int64) *GroupListCall {
	c.opt_["limit"] = limit
	return c
}

func (c *GroupListCall) Do() (*GroupList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["limit"]; ok {
		params.Set("limit", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
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
	ret := new(GroupList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all of the groups.",
	//   "httpMethod": "GET",
	//   "id": "update.group.list",
	//   "parameterOrder": [
	//     "appId"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "limit": {
	//       "default": "10",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     }
	//   },
	//   "path": "apps/{appId}/groups",
	//   "response": {
	//     "$ref": "GroupList"
	//   }
	// }

}

// method id "update.group.patch":

type GroupPatchCall struct {
	s     *Service
	appId string
	id    string
	group *Group
	opt_  map[string]interface{}
}

// Patch: Patch a group. This method supports patch semantics.
func (r *GroupService) Patch(appId string, id string, group *Group) *GroupPatchCall {
	c := &GroupPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.id = id
	c.group = group
	return c
}

func (c *GroupPatchCall) Do() (*Group, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.group)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Group)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Patch a group. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "update.group.patch",
	//   "parameterOrder": [
	//     "appId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/groups/{id}",
	//   "request": {
	//     "$ref": "Group"
	//   },
	//   "response": {
	//     "$ref": "Group"
	//   }
	// }

}

// method id "update.group.update":

type GroupUpdateCall struct {
	s     *Service
	appId string
	id    string
	group *Group
	opt_  map[string]interface{}
}

// Update: Patch a group.
func (r *GroupService) Update(appId string, id string, group *Group) *GroupUpdateCall {
	c := &GroupUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.id = id
	c.group = group
	return c
}

func (c *GroupUpdateCall) Do() (*Group, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.group)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", url.QueryEscape(c.id), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Group)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Patch a group.",
	//   "httpMethod": "PATCH",
	//   "id": "update.group.update",
	//   "parameterOrder": [
	//     "appId",
	//     "id"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "id": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/groups/{id}",
	//   "request": {
	//     "$ref": "Group",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "Group"
	//   }
	// }

}

// method id "update.group.requests.events.rollup":

type GroupRequestsEventsRollupCall struct {
	s         *Service
	appId     string
	groupId   string
	dateStart int64
	dateEnd   int64
	opt_      map[string]interface{}
}

// Rollup: Rollup all client requests by event for this group.
func (r *GroupRequestsEventsService) Rollup(appId string, groupId string, dateStart int64, dateEnd int64) *GroupRequestsEventsRollupCall {
	c := &GroupRequestsEventsRollupCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.groupId = groupId
	c.dateStart = dateStart
	c.dateEnd = dateEnd
	return c
}

// Resolution sets the optional parameter "resolution":
func (c *GroupRequestsEventsRollupCall) Resolution(resolution int64) *GroupRequestsEventsRollupCall {
	c.opt_["resolution"] = resolution
	return c
}

// Versions sets the optional parameter "versions":
func (c *GroupRequestsEventsRollupCall) Versions(versions string) *GroupRequestsEventsRollupCall {
	c.opt_["versions"] = versions
	return c
}

func (c *GroupRequestsEventsRollupCall) Do() (*GroupRequestsRollup, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["resolution"]; ok {
		params.Set("resolution", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["versions"]; ok {
		params.Set("versions", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups/{groupId}/requests/events/{dateStart}/{dateEnd}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{groupId}", url.QueryEscape(c.groupId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{dateStart}", strconv.FormatInt(c.dateStart, 10), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{dateEnd}", strconv.FormatInt(c.dateEnd, 10), 1)
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
	ret := new(GroupRequestsRollup)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Rollup all client requests by event for this group.",
	//   "httpMethod": "GET",
	//   "id": "update.group.requests.events.rollup",
	//   "parameterOrder": [
	//     "appId",
	//     "groupId",
	//     "dateStart",
	//     "dateEnd"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "dateEnd": {
	//       "format": "int64",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "dateStart": {
	//       "format": "int64",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "groupId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "resolution": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "versions": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/groups/{groupId}/requests/events/{dateStart}/{dateEnd}",
	//   "response": {
	//     "$ref": "GroupRequestsRollup"
	//   }
	// }

}

// method id "update.group.requests.versions.rollup":

type GroupRequestsVersionsRollupCall struct {
	s         *Service
	appId     string
	groupId   string
	dateStart int64
	dateEnd   int64
	opt_      map[string]interface{}
}

// Rollup: Rollup all clients requests by version for this group.
func (r *GroupRequestsVersionsService) Rollup(appId string, groupId string, dateStart int64, dateEnd int64) *GroupRequestsVersionsRollupCall {
	c := &GroupRequestsVersionsRollupCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	c.groupId = groupId
	c.dateStart = dateStart
	c.dateEnd = dateEnd
	return c
}

// Resolution sets the optional parameter "resolution":
func (c *GroupRequestsVersionsRollupCall) Resolution(resolution int64) *GroupRequestsVersionsRollupCall {
	c.opt_["resolution"] = resolution
	return c
}

// Versions sets the optional parameter "versions":
func (c *GroupRequestsVersionsRollupCall) Versions(versions string) *GroupRequestsVersionsRollupCall {
	c.opt_["versions"] = versions
	return c
}

func (c *GroupRequestsVersionsRollupCall) Do() (*GroupRequestsRollup, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["resolution"]; ok {
		params.Set("resolution", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["versions"]; ok {
		params.Set("versions", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "apps/{appId}/groups/{groupId}/requests/versions/{dateStart}/{dateEnd}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{appId}", url.QueryEscape(c.appId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{groupId}", url.QueryEscape(c.groupId), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{dateStart}", strconv.FormatInt(c.dateStart, 10), 1)
	req.URL.Path = strings.Replace(req.URL.Path, "{dateEnd}", strconv.FormatInt(c.dateEnd, 10), 1)
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
	ret := new(GroupRequestsRollup)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Rollup all clients requests by version for this group.",
	//   "httpMethod": "GET",
	//   "id": "update.group.requests.versions.rollup",
	//   "parameterOrder": [
	//     "appId",
	//     "groupId",
	//     "dateStart",
	//     "dateEnd"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "dateEnd": {
	//       "format": "int64",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "dateStart": {
	//       "format": "int64",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "groupId": {
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "resolution": {
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "versions": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}/groups/{groupId}/requests/versions/{dateStart}/{dateEnd}",
	//   "response": {
	//     "$ref": "GroupRequestsRollup"
	//   }
	// }

}

// method id "update.upstream.delete":

type UpstreamDeleteCall struct {
	s    *Service
	id   int64
	opt_ map[string]interface{}
}

// Delete: Delete an upstream.
func (r *UpstreamService) Delete(id int64) *UpstreamDeleteCall {
	c := &UpstreamDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

// Label sets the optional parameter "label":
func (c *UpstreamDeleteCall) Label(label string) *UpstreamDeleteCall {
	c.opt_["label"] = label
	return c
}

// Url sets the optional parameter "url":
func (c *UpstreamDeleteCall) Url(url string) *UpstreamDeleteCall {
	c.opt_["url"] = url
	return c
}

func (c *UpstreamDeleteCall) Do() (*Upstream, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["label"]; ok {
		params.Set("label", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["url"]; ok {
		params.Set("url", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative(c.s.BasePath, "upstream/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", strconv.FormatInt(c.id, 10), 1)
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
	ret := new(Upstream)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Delete an upstream.",
	//   "httpMethod": "DELETE",
	//   "id": "update.upstream.delete",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "format": "int32",
	//       "location": "path",
	//       "required": true,
	//       "type": "integer"
	//     },
	//     "label": {
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "url": {
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "upstream/{id}",
	//   "response": {
	//     "$ref": "Upstream"
	//   }
	// }

}

// method id "update.upstream.insert":

type UpstreamInsertCall struct {
	s        *Service
	upstream *Upstream
	opt_     map[string]interface{}
}

// Insert: Insert an upstream.
func (r *UpstreamService) Insert(upstream *Upstream) *UpstreamInsertCall {
	c := &UpstreamInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.upstream = upstream
	return c
}

func (c *UpstreamInsertCall) Do() (*Upstream, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.upstream)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "upstream")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Upstream)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Insert an upstream.",
	//   "httpMethod": "POST",
	//   "id": "update.upstream.insert",
	//   "path": "upstream",
	//   "request": {
	//     "$ref": "Upstream",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "Upstream"
	//   }
	// }

}

// method id "update.upstream.list":

type UpstreamListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: List all upstreams.
func (r *UpstreamService) List() *UpstreamListCall {
	c := &UpstreamListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UpstreamListCall) Do() (*UpstreamListResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "upstream")
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
	ret := new(UpstreamListResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "List all upstreams.",
	//   "httpMethod": "GET",
	//   "id": "update.upstream.list",
	//   "path": "upstream",
	//   "response": {
	//     "$ref": "UpstreamListResp"
	//   }
	// }

}

// method id "update.upstream.sync":

type UpstreamSyncCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// Sync: Synchronize all upstreams.
func (r *UpstreamService) Sync() *UpstreamSyncCall {
	c := &UpstreamSyncCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UpstreamSyncCall) Do() (*UpstreamSyncResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "upstream/sync")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
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
	ret := new(UpstreamSyncResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Synchronize all upstreams.",
	//   "httpMethod": "POST",
	//   "id": "update.upstream.sync",
	//   "path": "upstream/sync",
	//   "response": {
	//     "$ref": "UpstreamSyncResp"
	//   }
	// }

}

// method id "update.upstream.update":

type UpstreamUpdateCall struct {
	s        *Service
	id       int64
	upstream *Upstream
	opt_     map[string]interface{}
}

// Update: Update an upstream.
func (r *UpstreamService) Update(id int64, upstream *Upstream) *UpstreamUpdateCall {
	c := &UpstreamUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	c.upstream = upstream
	return c
}

func (c *UpstreamUpdateCall) Do() (*Upstream, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.upstream)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "upstream/{id}")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.URL.Path = strings.Replace(req.URL.Path, "{id}", strconv.FormatInt(c.id, 10), 1)
	googleapi.SetOpaque(req.URL)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Upstream)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Update an upstream.",
	//   "httpMethod": "PUT",
	//   "id": "update.upstream.update",
	//   "parameterOrder": [
	//     "id"
	//   ],
	//   "parameters": {
	//     "id": {
	//       "format": "int32",
	//       "location": "path",
	//       "required": true,
	//       "type": "integer"
	//     }
	//   },
	//   "path": "upstream/{id}",
	//   "request": {
	//     "$ref": "Upstream",
	//     "parameterName": "resource"
	//   },
	//   "response": {
	//     "$ref": "Upstream"
	//   }
	// }

}

// method id "update.util.uuid":

type UtilUuidCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// Uuid: Generate a new UUID.
func (r *UtilService) Uuid() *UtilUuidCall {
	c := &UtilUuidCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UtilUuidCall) Do() (*GenerateUuidResp, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative(c.s.BasePath, "util/uuid")
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
	ret := new(GenerateUuidResp)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Generate a new UUID.",
	//   "httpMethod": "GET",
	//   "id": "update.util.uuid",
	//   "path": "util/uuid",
	//   "response": {
	//     "$ref": "GenerateUuidResp"
	//   }
	// }

}
