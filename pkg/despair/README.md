# despair

Package despair provides an implemntation of a current Stereoscopic Depth Mapping Algorithm.

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# despair

```go
import "github.com/conneroisu/steroscopic-hardware/pkg/despair"
```

Package despair provides an implemntation of a current Stereoscopic Depth Mapping Algorithm.

## Index

- [func AssembleDisparityMap\(outputChan \<\-chan OutputChunk, dimensions image.Rectangle, chunks int\) \*image.Gray](<#AssembleDisparityMap>)
- [func LoadPNG\(filename string\) \(\*image.Gray, error\)](<#LoadPNG>)
- [func MustLoadPNG\(filename string\) \*image.Gray](<#MustLoadPNG>)
- [func MustSavePNG\(filename string, img image.Image\)](<#MustSavePNG>)
- [func RunSad\(left, right \*image.Gray, blockSize, maxDisparity int\) \*image.Gray](<#RunSad>)
- [func SavePNG\(filename string, img image.Image\) error](<#SavePNG>)
- [func SetupConcurrentSAD\(params \*Parameters, numWorkers int\) \(chan\<\- InputChunk, \<\-chan OutputChunk\)](<#SetupConcurrentSAD>)
- [type InputChunk](<#InputChunk>)
- [type OutputChunk](<#OutputChunk>)
- [type Parameters](<#Parameters>)
  - [func \(p \*Parameters\) Lock\(\)](<#Parameters.Lock>)
  - [func \(p \*Parameters\) Unlock\(\)](<#Parameters.Unlock>)


<a name="AssembleDisparityMap"></a>
## func [AssembleDisparityMap](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/main.go#L168-L172>)

```go
func AssembleDisparityMap(outputChan <-chan OutputChunk, dimensions image.Rectangle, chunks int) *image.Gray
```

AssembleDisparityMap assembles the disparity map from output chunks

<a name="LoadPNG"></a>
## func [LoadPNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L10>)

```go
func LoadPNG(filename string) (*image.Gray, error)
```

LoadPNG loads a PNG image and converts it to grayscale with optimizations

<a name="MustLoadPNG"></a>
## func [MustLoadPNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L44>)

```go
func MustLoadPNG(filename string) *image.Gray
```

MustLoadPNG loads a PNG image and converts it to grayscale with optimizations and panics if an error occurs.

<a name="MustSavePNG"></a>
## func [MustSavePNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L68>)

```go
func MustSavePNG(filename string, img image.Image)
```

MustSavePNG saves a PNG image with optimizations to the given filename and panics if an error occurs.

<a name="RunSad"></a>
## func [RunSad](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/main.go#L116-L119>)

```go
func RunSad(left, right *image.Gray, blockSize, maxDisparity int) *image.Gray
```

RunSad is a convenience function that sets up the pipeline, feeds the images, and assembles the disparity map.

This is not used in the web UI, but is useful for testing.

<a name="SavePNG"></a>
## func [SavePNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L54>)

```go
func SavePNG(filename string, img image.Image) error
```

SavePNG saves a PNG image with optimizations to the given filename and returns an error if one occurs.

<a name="SetupConcurrentSAD"></a>
## func [SetupConcurrentSAD](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/main.go#L28-L31>)

```go
func SetupConcurrentSAD(params *Parameters, numWorkers int) (chan<- InputChunk, <-chan OutputChunk)
```

SetupConcurrentSAD sets up a concurrent SAD processing pipeline It returns an input channel to feed image chunks into and an output channel to receive results from.

If the input channel is closed, the processing pipeline will stop.

<a name="InputChunk"></a>
## type [InputChunk](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/main.go#L12-L15>)

InputChunk represents a portion of the image to process

```go
type InputChunk struct {
    Left, Right *image.Gray
    Region      image.Rectangle
}
```

<a name="OutputChunk"></a>
## type [OutputChunk](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/main.go#L18-L21>)

OutputChunk represents the processed disparity data for a region

```go
type OutputChunk struct {
    DisparityData []uint8
    Region        image.Rectangle
}
```

<a name="Parameters"></a>
## type [Parameters](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/params.go#L9-L13>)

Parameters is a struct that holds the parameters for the stereoscopic image processing.

```go
type Parameters struct {
    BlockSize    int `json:"blockSize"`
    MaxDisparity int `json:"maxDisparity"`
    // contains filtered or unexported fields
}
```

<a name="Parameters.Lock"></a>
### func \(\*Parameters\) [Lock](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/params.go#L16>)

```go
func (p *Parameters) Lock()
```

Lock locks the mutex.

<a name="Parameters.Unlock"></a>
### func \(\*Parameters\) [Unlock](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/params.go#L19>)

```go
func (p *Parameters) Unlock()
```

Unlock unlocks the mutex.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
