package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _proto_micro_mall_pay_proto_pay_business_pay_business_swagger_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\x41\x6f\xdb\xb6\x17\xbf\xfb\x53\x10\xfa\xff\x8f\x41\xd3\x75\xc3\x0e\xbd\xa9\xb2\x9a\x18\x4d\x6c\x43\xb6\x87\x06\x43\x21\x30\x12\xed\xb0\x90\x48\x95\xa2\xd2\x1a\x43\x80\x5e\x0a\xf4\xb2\xed\xb0\x6e\x28\xb0\xc3\x2e\xc3\xae\xbd\x0c\x6b\xd1\xee\xe3\xc4\x59\x3f\xc6\x40\xd9\xb2\x45\x8a\xb2\x65\xd9\x59\x5d\x20\x02\x02\x44\x24\xdf\xe3\x7b\x7c\xef\xf7\xf8\xde\x93\xbf\x6b\x00\x60\xc4\x4f\xe1\x68\x84\x98\x71\x17\x18\x77\x6e\xdd\x36\xf6\xc4\x18\x26\x43\x6a\xdc\x05\x62\x1e\x00\x83\x63\x1e\x20\x31\x1f\x31\xca\xe9\x7e\x88\x3d\x46\xdd\x10\x06\x81\x1b\xc1\xb1\x3b\x1d\x14\xff\x9d\x26\x31\x26\x28\x8e\xa5\x97\x5b\xe9\x7c\xca\x16\x00\xe3\x1c\xb1\x18\x53\x22\x98\xcd\xfe\x05\x84\x72\x10\x23\x6e\x34\x00\xb8\x48\x37\xf7\x28\x89\x93\x10\xc5\xc6\x5d\xf0\xed\x94\x0a\x46\x51\x80\x3d\xc8\x31\x25\xfb\x8f\x63\x4a\xc4\xda\x47\xe9\xda\x88\x51\x3f\xf1\x2a\xae\x85\xfc\x2c\x5e\x68\xb5\x7f\xfe\xc5\x3e\xf4\x3c\x9a\x10\x3e\x1f\x04\xc0\x18\xa1\xfc\xab\x38\x9f\x24\x0c\x21\x1b\x0b\x91\x3f\xfe\xf0\x76\xf2\xe3\x2f\x1f\xff\xfc\xe3\xea\xe5\xdb\x99\x46\xe9\x12\x1a\x21\x96\xee\xd8\xf2\xc5\xb2\x2e\x1c\xdf\x9b\x69\xdf\x43\xec\x1c\x7b\xc8\xbd\x8f\x89\x6f\xce\x36\xcb\x11\x32\x14\x47\x94\xc4\x28\x96\xb6\x04\xc0\xb8\x73\xfb\xb6\x32\x04\x80\xe1\xa3\xd8\x63\x38\xe2\xb3\x03\x34\x41\x9c\x78\x1e\x8a\xe3\x61\x12\x80\x8c\xd3\xad\x1c\xfb\xa9\xf8\xde\x19\x0a\x61\x81\x19\x00\xc6\xff\x19\x1a\x0a\x3e\xff\xdb\xf7\xd1\x10\x13\x2c\xf8\xca\xb6\xcb\x49\xed\xcc\xf8\x1b\x12\x97\x8b\xdc\xdb\x45\x7e\x63\xc3\x47\x43\x98\x04\x7c\xb5\x12\x04\x24\x04\x3d\x8b\x90\xc7\x91\x0f\x10\x63\x94\x6d\x4f\x17\x96\x10\x8e\x43\x64\x0b\xae\x4b\x04\x6f\x68\x54\x30\x22\xc8\x60\x88\x38\x62\x0b\xdf\x9a\x3e\x8a\x3e\x04\x86\x29\x36\xe8\x53\x82\x98\x2a\x30\x4e\x75\x7c\x92\x20\x36\x56\xa7\x18\x7a\x92\x60\x86\x84\xc3\x0c\x61\x10\x23\x65\x9a\x8f\xa3\x94\x2d\x64\x0c\x16\x68\x31\x47\xa1\xea\x32\x12\x55\xcc\x19\x26\x23\x45\x65\x85\x89\x47\x83\x00\x79\xe2\xa0\xee\x53\x16\x42\x61\x2a\x23\x4c\x02\x8e\x8d\x32\xa3\x96\x68\x3e\x03\x91\x9b\x6e\xbe\xf5\x03\x98\xa9\xa2\xcc\x22\x92\x84\x8a\x5d\xd2\xf1\x2e\x62\x02\xf3\x7b\xea\xb8\x45\xc3\x08\x12\x55\x06\x00\x8c\xde\x38\xe6\x28\x94\x4f\xea\xd1\x9e\xea\xb2\x99\x2f\x67\xfc\xd7\x3c\x21\x8f\x62\xb2\x13\xc7\x63\xb5\x4f\x8a\x47\x30\xe8\x35\x2b\xeb\x2f\x18\x68\xa1\x93\xa3\x31\x38\x1c\xa9\xa0\xd1\xc4\xc4\x05\x9f\x47\x8d\x3c\xb7\xd9\x81\xe6\xe3\xf3\xbe\x77\x06\xd9\x08\xe5\xc3\x74\x44\xe3\x25\x71\x3a\x8d\xd0\x93\x17\x2f\x26\xcf\xff\x5e\x2f\x4e\xcf\xa2\x9d\x35\xdd\xef\xf3\x89\xd4\x92\xdc\x37\xb1\x3a\x7b\x4a\xf0\x78\x4a\xfd\x62\x4c\x25\x65\x33\x39\x24\x72\x96\xa8\x40\xdc\xa2\xd5\x9e\x24\x28\xe6\x55\x74\xbf\x5e\xb0\x09\x51\x2b\x43\x6d\xf2\xf2\xd7\xc9\x87\xf7\x75\x52\x22\x8b\x21\xc8\xd1\xe7\x97\x14\x49\x72\xdf\x40\x2d\x7b\x76\x1c\x6a\x8a\xd5\x3e\x31\xd4\x38\x83\x3e\x12\xf2\x55\xc6\xd9\x3f\x1f\x7e\xbb\x7c\xf7\xfc\xea\xd5\x5f\x93\xef\x7f\xbe\x7a\xf5\xe6\xf2\xc3\xeb\xf5\xd0\xd6\x17\x3b\x76\xa5\x54\x72\xd7\x81\x96\x89\x7c\x83\xb1\xec\xd9\x71\x8c\x2d\x0c\xb6\x13\xf0\x4a\x12\xec\xaf\x57\xd9\x5f\xbe\xff\xfd\xea\xf5\x4f\x93\x57\x6f\x2e\xdf\x3d\x6f\x35\xd7\x03\xd8\x01\xe2\xa9\xfe\x83\x81\x4c\xb9\xe3\x20\xcb\x8b\x7d\x03\xb4\xec\x29\x01\x5a\xea\x57\x3b\x51\xc8\x75\x4d\x4d\x21\x67\x1d\x9a\xce\x81\x5d\xbd\x96\x35\x4f\xd6\x2e\x64\x03\x8c\x08\x77\xb1\x7f\x6d\xfa\xaf\x29\x50\x72\x8d\xa2\x28\xb3\xc3\x79\x97\x04\x13\xfe\xf5\x57\xd7\x5e\x06\xcf\xfb\xa1\x39\x6f\x5f\x74\x2f\x57\x96\x11\xb9\xc0\x97\xe9\x45\x4f\x1f\x23\x6f\x91\x6a\x1b\x11\x13\x61\x8d\x63\x25\x36\xcd\x5a\x59\x4a\xb8\x2a\x6f\x48\x95\xb4\xa3\x96\x59\x56\x8b\x46\xa9\x8f\xa4\xec\x5e\xb9\x96\xea\x0b\x62\x2d\xfb\x45\x13\x66\x6d\xde\x16\xc5\xa4\x9c\x31\x0c\x95\x1e\xf2\x52\xed\xf3\x94\x34\xe1\xee\x34\xa6\x10\x5a\x8b\x3e\x72\x43\xc4\xd5\x60\x59\x41\xa1\x4e\x76\x9d\x1d\x0b\xf2\x86\x6a\x17\xf9\x46\x5d\xdd\x67\xd8\xc0\xd7\x3c\x1a\x86\xe9\x2d\xb1\xb6\x49\x04\x5d\xf1\xc6\xaa\xac\x82\x4d\x38\x1b\x5f\x1b\x4a\x96\x59\xed\x73\xf5\xf3\x53\x18\x40\xe2\x15\xd8\x96\xaa\x5c\xd5\x12\x7d\x59\xd6\xb2\x30\x5c\xbc\x0b\x8b\xdd\x5e\x4d\x9f\x57\xe9\xf0\xce\x03\x74\x59\x57\x57\x27\xe8\xfc\x5c\x6a\x49\x29\x37\x5d\xf3\xed\x56\xad\x30\xf3\x16\xab\x5e\x12\xc9\xed\x37\x02\x9e\x5f\xc3\x43\x1c\xc4\x2d\x41\xa8\x75\x90\x30\x1e\x6d\xcf\x39\xb4\x05\xfb\x0d\x5c\x3f\x21\x5c\xf5\x8d\xaf\x5d\x8a\xfd\x1a\xab\xe9\xdc\xbc\xfe\x11\xe8\x3e\x88\xee\xe8\x01\x20\x71\xbf\xb9\x01\x8e\x4b\x33\x93\xea\x99\x5c\x65\x37\x9f\x5e\xaa\xfa\x4c\x6f\xe5\xe1\x6a\x4b\xd1\x1d\x3d\x5d\xa5\xb7\xb0\xa1\x5b\xc9\xe9\xd8\x26\x41\x2e\x72\x57\x0b\x26\x59\xbc\xb4\xaa\x51\xf2\x4c\x1c\xd5\x4c\x50\xa3\x00\x72\xb1\x49\x4d\x72\x1f\xa5\x45\xd2\xd6\x8e\x3a\xbb\xc0\x6a\xdd\xe4\xbd\x81\x65\xd9\xbd\x5e\xfe\x36\xb7\x1d\xa7\xe3\xc8\xd7\xbb\xed\xb8\xed\x4e\xdf\xb5\x1f\xb6\x7a\xfd\xc2\x4c\x61\xf4\xd8\x76\xac\x43\xb3\xdd\xd7\xd3\xcc\x67\x0b\x33\xbd\xc3\x4e\x57\x4f\x93\xce\x14\x47\x1f\x0c\x4a\x96\x3f\x18\xe8\x07\xcd\xe3\xce\x20\x93\xab\xdd\x19\x1c\x1c\x16\x94\xb9\x67\x1e\x99\x6d\xcb\x2e\x59\x32\x97\x7d\xf9\x32\xd3\xb2\xd2\x7d\x8e\x3a\xd6\x03\x2d\x79\xb6\x60\xf9\x01\x65\xab\x7a\x7d\xb3\x6f\x17\x98\xa5\xe2\x2e\x65\x24\xad\xd0\x33\x69\xda\x56\xeb\xd8\x3c\x72\xbb\xa6\xd3\xb3\x5d\xdb\x91\xec\xde\x77\xcc\x76\xcf\xb4\xfa\xad\x4e\xdb\xbd\x6f\xb6\x8e\xec\xa6\x4e\xc9\xc2\xb6\x7d\xc7\x6c\xda\x6e\xd7\x3c\x71\x9d\x41\x5b\x3f\xa1\xf1\xba\xc5\xa4\xfd\xb0\xdb\x72\xec\xe2\x5c\xc7\x69\xce\xfc\xf0\xd8\xec\x5b\x87\xae\x50\x6f\x85\xba\xad\xf6\x37\xe6\x91\xdc\xbc\x4c\x57\x4d\xdb\x4a\xae\x63\x5b\x1d\xa7\x59\xa6\x81\x08\xda\xae\x7d\xdc\xed\x9f\x2c\x4d\x6e\x33\x5d\x4a\xf1\x99\xde\x00\x07\x94\xfa\x71\x13\x71\x88\x83\x4d\xa2\x61\xfd\xc2\x9c\x21\x3f\xf1\x66\xed\xcb\x2d\x05\x9e\xac\x39\xbe\x79\xdd\xb9\x61\xd3\x80\xe3\x10\xb9\xe8\x59\x84\x59\xf5\xb0\x9a\x23\x27\x94\xe3\xe1\xd8\x4d\x58\x50\x87\x5a\xee\x0c\xaf\x4d\x1e\x22\xe6\x9d\xc1\x7a\x36\x85\x9c\x43\xef\xac\x9e\xd0\x8a\x27\x82\x6a\x89\x44\xc1\x97\x6b\xb8\xcb\x16\x2a\x9f\xe2\x8f\x1c\x2b\x6b\x7e\x6d\xc5\xc9\x7f\x9b\x9d\xca\xd8\xd3\xa7\xa7\x7b\x3b\x9c\x45\x09\xc4\x3f\xdb\x72\x29\x53\xf8\xba\xba\xa3\x99\xf6\x34\xd4\x6d\x33\xdb\x4e\x35\xdf\xa0\xe9\x24\x7d\x7a\x51\x3e\xba\xe8\xdb\x4b\xa6\xda\xd1\x61\x94\xd3\xd3\x64\x68\x92\x8d\x6e\x02\xb1\xbe\x6e\x20\x3e\x87\x41\xb2\xca\x9b\x4a\x5c\xfc\x74\xcc\x57\x36\x5c\xa5\xaf\x72\x1b\xa8\x88\x14\x06\x95\xf5\x5b\x06\x16\x4c\x38\x1a\x29\x3f\xd6\x95\x21\xfc\xe5\x9d\xb2\xfb\x27\x8e\xe1\xa8\xd6\xbd\x39\xbd\x44\x0a\x5f\x81\xb7\x16\xee\x72\x2e\xa5\x0f\x70\xb2\x9d\x1a\xe2\xef\xa2\xf1\x6f\x00\x00\x00\xff\xff\x5d\xa6\x95\x93\x02\x30\x00\x00")

func proto_micro_mall_pay_proto_pay_business_pay_business_swagger_json() ([]byte, error) {
	return bindata_read(
		_proto_micro_mall_pay_proto_pay_business_pay_business_swagger_json,
		"proto/micro_mall_pay_proto/pay_business/pay_business.swagger.json",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"proto/micro_mall_pay_proto/pay_business/pay_business.swagger.json": proto_micro_mall_pay_proto_pay_business_pay_business_swagger_json,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"proto": &_bintree_t{nil, map[string]*_bintree_t{
		"micro_mall_pay_proto": &_bintree_t{nil, map[string]*_bintree_t{
			"pay_business": &_bintree_t{nil, map[string]*_bintree_t{
				"pay_business.swagger.json": &_bintree_t{proto_micro_mall_pay_proto_pay_business_pay_business_swagger_json, map[string]*_bintree_t{
				}},
			}},
		}},
	}},
}}
