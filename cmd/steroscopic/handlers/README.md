# handlers

Http handlers for the web ui can be found here.

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# handlers

```go
import "github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/handlers"
```

Package handlers contains functions for handling API requests.

## Index

- [func Make\(fn APIFn\) http.HandlerFunc](<#Make>)
- [func MorphableHandler\(wrapper func\(templ.Component\) templ.Component, morph templ.Component\) http.HandlerFunc](<#MorphableHandler>)
- [func PreviewSeqHandler\(w http.ResponseWriter, r \*http.Request\)](<#PreviewSeqHandler>)
- [type APIFn](<#APIFn>)
  - [func ConfigureCamera\(logger \*logger.Logger, params \*despair.Parameters, leftStream, rightStream, outputStream \*camera.StreamManager, isLeft bool\) APIFn](<#ConfigureCamera>)
  - [func GetPorts\(logger \*logger.Logger\) APIFn](<#GetPorts>)
  - [func ManualCalcDepthMapHandler\(logger \*logger.Logger\) APIFn](<#ManualCalcDepthMapHandler>)
  - [func ParametersHandler\(logger \*logger.Logger, params \*despair.Parameters\) APIFn](<#ParametersHandler>)
  - [func StreamHandlerFn\(manager \*camera.StreamManager\) APIFn](<#StreamHandlerFn>)


<a name="Make"></a>
## func [Make](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/api.go#L15>)

```go
func Make(fn APIFn) http.HandlerFunc
```

Make returns a function that can be used as an http.HandlerFunc.

<a name="MorphableHandler"></a>
## func [MorphableHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/api.go#L34-L37>)

```go
func MorphableHandler(wrapper func(templ.Component) templ.Component, morph templ.Component) http.HandlerFunc
```

MorphableHandler returns a handler that checks for the presence of the hx\-trigger header and serves either the full or morphed view.

<a name="PreviewSeqHandler"></a>
## func [PreviewSeqHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/delim.go#L15>)

```go
func PreviewSeqHandler(w http.ResponseWriter, r *http.Request)
```

PreviewSeqHandler handles requests to preview sequences in different formats

<a name="APIFn"></a>
## type [APIFn](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/api.go#L12>)

APIFn is a function that handles an API request.

```go
type APIFn func(w http.ResponseWriter, r *http.Request) error
```

<a name="ConfigureCamera"></a>
### func [ConfigureCamera](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/configure.go#L58-L63>)

```go
func ConfigureCamera(logger *logger.Logger, params *despair.Parameters, leftStream, rightStream, outputStream *camera.StreamManager, isLeft bool) APIFn
```

ConfigureCamera handles client requests to configure all camera parameters at once.

<a name="GetPorts"></a>
### func [GetPorts](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/targets.go#L14-L16>)

```go
func GetPorts(logger *logger.Logger) APIFn
```

GetPorts handles client requests to configure the camera.

<a name="ManualCalcDepthMapHandler"></a>
### func [ManualCalcDepthMapHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/manual.go#L10-L12>)

```go
func ManualCalcDepthMapHandler(logger *logger.Logger) APIFn
```

ManualCalcDepthMapHandler is a handler for the manual depth map calculation endpoint.

<a name="ParametersHandler"></a>
### func [ParametersHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/configure.go#L15>)

```go
func ParametersHandler(logger *logger.Logger, params *despair.Parameters) APIFn
```

ParametersHandler handles client requests to change the parameters of the desparity map generator.

<a name="StreamHandlerFn"></a>
### func [StreamHandlerFn](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/steroscopic/handlers/stream.go#L17>)

```go
func StreamHandlerFn(manager *camera.StreamManager) APIFn
```

StreamHandlerFn returns a handler for streaming camera images to multiple clients

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
