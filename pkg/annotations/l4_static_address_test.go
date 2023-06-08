package annotations

import (
	"fmt"
	"testing"

	"google.golang.org/api/compute/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cloud-provider-gcp/providers/gce"
)

func TestAddressFromAnnotation(t *testing.T) {
	vals := gce.DefaultTestClusterValues()

	ipv4AddressString := "123.0.124.1"
	ipv4AddressNoVersionString := "123.0.111.1"
	ipv6AddressString := "0::1"

	ipv4AddressName := "ipv4-address"
	ipv4AddressNoVersionName := "ipv4-address-no-version"
	ipv6AddressName := "ipv6-address"

	ipv4Address := compute.Address{
		Name:      ipv4AddressName,
		Address:   ipv4AddressString,
		IpVersion: IPv4Version,
	}
	ipv4AddressNoVersion := compute.Address{
		Name:    ipv4AddressNoVersionName,
		Address: ipv4AddressNoVersionString,
	}
	ipv6Address := compute.Address{
		Name:      ipv6AddressName,
		Address:   ipv6AddressString,
		IpVersion: IPv6Version,
	}
	testCases := []struct {
		desc              string
		reservedAddresses []compute.Address
		annotationVal     string
		wantIPv4Address   string
		wantIPv6Address   string
	}{
		{
			desc: "Single existing IPv4 address",
			reservedAddresses: []compute.Address{
				ipv4Address,
			},
			annotationVal:   fmt.Sprintf("%s", ipv4Address.Name),
			wantIPv4Address: ipv4Address.Address,
		},
		{
			desc: "Single existing IPv6 address",
			reservedAddresses: []compute.Address{
				ipv6Address,
			},
			annotationVal:   fmt.Sprintf("%s", ipv6Address.Name),
			wantIPv6Address: ipv6Address.Address,
		},
		{
			desc: "Single existing IPv4 address with no IpVersion",
			reservedAddresses: []compute.Address{
				ipv4AddressNoVersion,
			},
			annotationVal:   fmt.Sprintf("%s", ipv4AddressNoVersion.Name),
			wantIPv4Address: ipv4AddressNoVersion.Address,
		},
		{
			desc:              "Many non-existing IPv4 and IPv6 addresses",
			reservedAddresses: []compute.Address{},
			annotationVal:     fmt.Sprintf("%s, %s,%s, %s", ipv4Address.Name, ipv4Address.Name, ipv6Address.Name, ipv4AddressNoVersion.Name),
			wantIPv4Address:   "",
			wantIPv6Address:   "",
		},
		{
			desc: "Repeated existing IPv4 addresses",
			reservedAddresses: []compute.Address{
				ipv4Address,
			},
			annotationVal:   fmt.Sprintf("%s, %s", ipv4Address.Name, ipv4Address.Name),
			wantIPv4Address: ipv4Address.Address,
		},
		{
			desc: "IPv4 and IPv6 addresses",
			reservedAddresses: []compute.Address{
				ipv4Address,
				ipv6Address,
			},
			annotationVal:   fmt.Sprintf("%s, %s", ipv4Address.Name, ipv6Address.Name),
			wantIPv4Address: ipv4Address.Address,
			wantIPv6Address: ipv6Address.Address,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			svc := &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						StaticL4AddressesAnnotationKey: tc.annotationVal,
					},
				},
			}
			fakeGCE := gce.NewFakeGCECloud(vals)
			for i := range tc.reservedAddresses {
				err := fakeGCE.ReserveRegionAddress(&tc.reservedAddresses[i], fakeGCE.Region())
				if err != nil {
					t.Fatalf("fakeGCE.ReserveRegionAddress(%v, %s) returned error %v", tc.reservedAddresses[i], fakeGCE.Region(), err)
				}
			}

			ipv4Addr, err := FromService(svc).IPv4AddressAnnotation(fakeGCE)
			if err != nil {
				t.Fatalf("IPv4AddressAnnotation(..., %s) returned error %v", tc.annotationVal, err)
			}
			if ipv4Addr != tc.wantIPv4Address {
				t.Errorf("IPv4AddressAnnotation(..., %s) returned %s, not equal to expected = %s", tc.annotationVal, ipv4Addr, tc.wantIPv4Address)
			}

			ipv6Addr, err := FromService(svc).IPv6AddressAnnotation(fakeGCE)
			if err != nil {
				t.Fatalf("IPv6AddressAnnotation(..., %s) returned error %v", tc.annotationVal, err)
			}
			if ipv6Addr != tc.wantIPv6Address {
				t.Errorf("IPv6AddressAnnotation(..., %s) returned %s, not equal to expected = %s", tc.annotationVal, ipv6Addr, tc.wantIPv6Address)
			}
			if ipv6Addr != tc.wantIPv6Address {
				t.Errorf("IPv6AddressAnnotation(..., %s) returned %s, not equal to expected = %s", tc.annotationVal, ipv6Addr, tc.wantIPv6Address)
			}
		})
	}
}
