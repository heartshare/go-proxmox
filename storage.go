package proxmox

import (
	"fmt"
	"net/url"

	_ "github.com/davecgh/go-spew/spew"
)

type Storage struct {
	StorageType string
	Active      float64
	Total       float64
	Content     string
	Shared      float64
	Storage     string
	Used        float64
	Avail       float64
	node        Node
}

type StorageList map[string]Storage

func (storage Storage) CreateVolume(FileName string, DiskSize string, VmId string) (map[string]interface{}, error) {
	var form url.Values
	var err error
	var data map[string]interface{}
	var target string

	fmt.Println("!CreateVolume")

	form = url.Values{
		"filename": {FileName},
		//		"node":     {storage.node.Node},
		"size":   {DiskSize},
		"format": {"qcow2"},
		"vmid":   {VmId},
	}

	target = "nodes/" + storage.node.Node + "/storage/" + storage.Storage + "/content"
	data, err = storage.node.proxmox.PostForm(target, form)
	if err != nil {
		fmt.Println("Error!!!")
		return nil, err
	}
	fmt.Println("Storage created")
	return data, err
}

func (storage Storage) Volumes() (VolumeList, error) {
	var err error
	var target string
	var data map[string]interface{}
	var list VolumeList
	var volume Volume
	var results []interface{}

	fmt.Println("!Volumes")

	target = "nodes/" + storage.node.Node + "/storage/" + storage.Storage + "/content"
	data, err = storage.node.proxmox.Get(target)
	if err != nil {
		return nil, err
	}

	list = make(VolumeList)
	results = data["data"].([]interface{})
	for _, v0 := range results {
		v := v0.(map[string]interface{})
		volume = Volume{
			Size:    v["size"].(float64),
			VolId:   v["volid"].(string),
			VmId:    v["vmid"].(string),
			Format:  v["format"].(string),
			Content: v["content"].(string),
			Used:    v["used"].(float64),
			storage: storage,
		}
		list[volume.VolId] = volume
	}

	return list, nil
}
