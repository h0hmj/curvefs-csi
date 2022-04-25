/*
Copyright 2022 The Curve Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package curvefsdriver

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os/exec"
	"strconv"
)

const (
	defaultToolExampleConfPath   = "/curvefs/conf/tools.conf"
	defaultClientExampleConfPath = "/curvefs/conf/client.conf"
	toolPath                     = "/curvefs/tools/sbin/curvefs_tool"
	clientPath                   = "/curvefs/client/sbin/curve-fuse"
)

type curvefsTool struct {
	toolParams map[string]string
}

func NewCurvefsTool() *curvefsTool {
	return &curvefsTool{toolParams: map[string]string{}}
}

func (ct *curvefsTool) CreateFs(
	volumeID string,
	capacity int64,
	params map[string]string,
) error {
	ct.validateCommonParams(params)
	ct.validateCreateFsParams(params)
	ct.toolParams["fsName"] = volumeID
	// todo: current capacity is not working
	// call curvefs_tool create-fs to create a fs
	createFsArgs := []string{"create-fs"}
	for k, v := range ct.toolParams {
		arg := fmt.Sprintf("-%s=%s", k, v)
		createFsArgs = append(createFsArgs, arg)
	}
	createFsCmd := exec.Command(toolPath, createFsArgs...)
	output, err := createFsCmd.CombinedOutput()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"Curvefs_tool create-fs failed. cmd: %s %v, output: %s, err: %v",
			toolPath,
			createFsArgs,
			output,
			err,
		)
	}
	return nil
}

func (ct *curvefsTool) DeleteFs(volumeID string, params map[string]string) error {
	ct.validateCommonParams(params)
	ct.toolParams["fsName"] = volumeID
	ct.toolParams["noconfirm"] = "1"
	// call curvefs_tool delete-fs to create a fs
	deleteFsArgs := []string{"delete-fs"}
	for k, v := range ct.toolParams {
		arg := fmt.Sprintf("-%s=%s", k, v)
		deleteFsArgs = append(deleteFsArgs, arg)
	}
	deleteFsCmd := exec.Command(toolPath, deleteFsArgs...)
	output, err := deleteFsCmd.CombinedOutput()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"curvefs_tool delete-fs failed. cmd:%s %v, output: %s, err: %v",
			toolPath,
			deleteFsArgs,
			output,
			err,
		)
	}
	return nil
}

func (ct *curvefsTool) validateCommonParams(params map[string]string) error {
	if mdsAddr, ok := params["mdsAddr"]; ok {
		ct.toolParams["mdsAddr"] = mdsAddr
	} else {
		return status.Error(codes.InvalidArgument, "mdsAddr is missing")
	}
	if confPath, ok := params["toolConfPath"]; ok {
		ct.toolParams["confPath"] = confPath
	} else {
		ct.toolParams["confPath"] = defaultToolExampleConfPath
	}
	return nil
}

func (ct *curvefsTool) validateCreateFsParams(params map[string]string) error {
	if fsType, ok := params["fsType"]; ok {
		ct.toolParams["fsType"] = fsType
		if fsType == "s3" {
			s3Endpoint, ok1 := params["s3Endpoint"]
			s3AccessKey, ok2 := params["s3AccessKey"]
			s3SecretKey, ok3 := params["s3SecretKey"]
			s3Bucket, ok4 := params["s3Bucket"]
			if ok1 && ok2 && ok3 && ok4 {
				ct.toolParams["s3_endpoint"] = s3Endpoint
				ct.toolParams["s3_ak"] = s3AccessKey
				ct.toolParams["s3_sk"] = s3SecretKey
				ct.toolParams["s3_bucket_name"] = s3Bucket
			} else {
				return status.Error(codes.InvalidArgument, "s3Info is incomplete")
			}
		} else if fsType == "volume" {
			if backendVolName, ok := params["backendVolName"]; ok {
				ct.toolParams["volumeName"] = backendVolName
			} else {
				return status.Error(codes.InvalidArgument, "backendVolName is missing")
			}
			if backendVolSizeGB, ok := params["backendVolSizeGB"]; ok {
				backendVolSizeGBInt, err := strconv.ParseInt(backendVolSizeGB, 0, 64)
				if err != nil {
					return status.Error(codes.InvalidArgument, "backendVolSize is not integer")
				}
				if backendVolSizeGBInt < 10 {
					return status.Error(codes.InvalidArgument, "backendVolSize must larger than 10GB")
				}
				ct.toolParams["volumeSize"] = backendVolSizeGB
			} else {
				return status.Error(codes.InvalidArgument, "backendVolSize is missing")
			}
		} else {
			return status.Errorf(codes.InvalidArgument, "unsupported fsType %s", fsType)
		}
	} else {
		return status.Error(codes.InvalidArgument, "fsType is missing")
	}
	return nil
}

type curvefsMounter struct {
	mounterParams map[string]string
}

func NewCurvefsMounter() *curvefsMounter {
	return &curvefsMounter{mounterParams: map[string]string{}}
}

func (cm *curvefsMounter) MountFs(
	fsname string,
	targetPath string,
	params map[string]string,
) error {
	cm.validateMountFsParams(params)
	cm.mounterParams["fsname"] = fsname
	// call curve-fuse -o conf=/etc/curvefs/client.conf -o fsname=testfs \
	//       -o fstype=s3  --mdsAddr=1.1.1.1 <mountpoint>
	var mountFsArgs []string
	doubleDashArgs := map[string]string{"mdsaddr": ""}
	for k, v := range cm.mounterParams {
		if _, ok := doubleDashArgs[k]; ok {
			arg := fmt.Sprintf("--%s=%s", k, v)
			mountFsArgs = append(mountFsArgs, arg)
		} else {
			mountFsArgs = append(mountFsArgs, "-o")
			arg := fmt.Sprintf("%s=%s", k, v)
			mountFsArgs = append(mountFsArgs, arg)
		}
	}
	mountFsArgs = append(mountFsArgs, targetPath)
	mountFsCmd := exec.Command(clientPath, mountFsArgs...)
	output, err := mountFsCmd.CombinedOutput()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"curve-fuse mount failed. cmd: %s %v, output: %s, err: %v",
			clientPath,
			mountFsArgs,
			output,
			err,
		)
	}
	return nil
}

func (cm *curvefsMounter) UmountFs(targetPath string) error {
	umountFsCmd := exec.Command("umount", targetPath)
	output, err := umountFsCmd.CombinedOutput()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"umount %s failed. output: %s, err: %v",
			targetPath,
			output,
			err,
		)
	}
	return nil
}

func (cm *curvefsMounter) validateMountFsParams(params map[string]string) error {
	if mdsAddr, ok := params["mdsAddr"]; ok {
		cm.mounterParams["mdsaddr"] = mdsAddr
	} else {
		return status.Error(codes.InvalidArgument, "mdsAddr is missing")
	}
	if confPath, ok := params["clientConfPath"]; ok {
		cm.mounterParams["conf"] = confPath
	} else {
		cm.mounterParams["conf"] = defaultClientExampleConfPath
	}
	if fsType, ok := params["fsType"]; ok {
		cm.mounterParams["fstype"] = fsType
	} else {
		return status.Error(codes.InvalidArgument, "fsType is missing")
	}
	return nil
}
