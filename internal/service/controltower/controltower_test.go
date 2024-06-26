// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package controltower_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccControlTower_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"LandingZone": {
			"basic":        testAccLandingZone_basic,
			"disappears":   testAccLandingZone_disappears,
			names.AttrTags: testAccLandingZone_tags,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
