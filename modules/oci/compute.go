package oci

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

// DeleteImage deletes a custom image with given OCID.
func DeleteImage(t *testing.T, ocid string) {
	err := DeleteImageE(t, ocid)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteImageE deletes a custom image with given OCID.
func DeleteImageE(t *testing.T, ocid string) error {
	logger.Logf(t, "Deleting image with OCID %s", ocid)

	configProvider := common.DefaultConfigProvider()
	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		return err
	}

	_, err = computeClient.DeleteImage(context.Background(), core.DeleteImageRequest{ImageId: &ocid})
	return err
}

// BaseImageOcid gets the OCID of the most recent image in the given compartment that has the given OS name and version.
func BaseImageOcid(t *testing.T, compartmentOcid string, osName string, osVersion string) string {
	ocid, err := BaseImageOcidE(t, compartmentOcid, osName, osVersion)
	if err != nil {
		t.Fatal(err)
	}
	return ocid
}

// BaseImageOcidE gets the OCID of the most recent image in the given compartment that has the given OS name and version.
func BaseImageOcidE(t *testing.T, compartmentOcid string, osName string, osVersion string) (string, error) {
	configProvider := common.DefaultConfigProvider()
	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		return "", err
	}

	request := core.ListImagesRequest{
		CompartmentId:          &compartmentOcid,
		OperatingSystem:        &osName,
		OperatingSystemVersion: &osVersion,
	}
	response, err := computeClient.ListImages(context.Background(), request)
	if err != nil {
		return "", err
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("No %s %s images found in the %s compartment", osName, osVersion, compartmentOcid)
	}

	mostRecentImage := mostRecentImage(response.Items)
	return *mostRecentImage.Id, nil

}

// Image sorting code borrowed from: https://github.com/hashicorp/packer/blob/7f4112ba229309cfc0ebaa10ded2abdfaf1b22c8/builder/amazon/common/step_source_ami_info.go
type imageSort []core.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	iTime := a[i].TimeCreated.Unix()
	jTime := a[j].TimeCreated.Unix()
	return iTime < jTime
}

// mostRecentImage returns the most recent image out of a slice of images.
func mostRecentImage(images []core.Image) core.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
