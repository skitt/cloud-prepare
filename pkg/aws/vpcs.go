/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

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

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
)

func (ac *awsCloud) getVpcID() (string, error) {
	var err error
	var result *ec2.DescribeVpcsOutput

	if vpcID, exists := ac.cloudConfig[VPCIDKey]; exists {
		vpcIDStr, ok := vpcID.(string)
		if !ok || vpcIDStr == "" {
			return "", errors.New("VPC ID needs to be a valid non-empty string")
		}

		return vpcIDStr, nil
	}

	ownedFilters := ac.filterByCurrentCluster()
	vpcName := ac.withAWSInfo("{infraID}-vpc")

	for i := range ownedFilters {
		filters := []types.Filter{
			ac.filterByName(vpcName),
			ownedFilters[i],
		}

		result, err = ac.client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{Filters: filters})
		if err != nil {
			return "", errors.Wrap(err, "error describing AWS VPCs")
		}

		if len(result.Vpcs) != 0 {
			break
		}
	}

	if len(result.Vpcs) == 0 {
		return "", newNotFoundError("VPC %s", vpcName)
	}

	return *result.Vpcs[0].VpcId, nil
}
