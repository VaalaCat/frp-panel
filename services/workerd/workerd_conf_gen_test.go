package workerd

import (
	"testing"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestBuildCapfile(t *testing.T) {
	tests := []struct {
		name   string
		wokers []*pb.Worker
		expect func(t *testing.T, resp map[string]string)
	}{
		{
			name: "common case",
			wokers: []*pb.Worker{
				{
					WorkerId:  lo.ToPtr("test"),
					CodeEntry: lo.ToPtr("test/entry.js"),
					Socket: &pb.Socket{
						Address: lo.ToPtr("unix:/test/test.sock"),
					},
				},
				{
					WorkerId:  lo.ToPtr("test1"),
					CodeEntry: lo.ToPtr("test1/entry.js"),
					Socket: &pb.Socket{
						Address: lo.ToPtr("unix:/test1/test.sock"),
					},
				},
			},

			expect: func(t *testing.T, result map[string]string) {
				assert.Equal(t,
					`using Workerd = import "/workerd/workerd.capnp";

const config :Workerd.Config = (
  services = [
    (name = "test", worker = .vtestWorker),
  ],

  sockets = [
    (
      name = "test",
      address = "unix:/test/test.sock",
      http=(),
      service="test"
    ),
  ]
);

const vtestWorker :Workerd.Worker = (
  modules = [
    (name = "test/entry.js", esModule = embed "src/test/entry.js"),
  ],
  compatibilityDate = "2023-04-03",
);`, result["test"])

				assert.Equal(t, `using Workerd = import "/workerd/workerd.capnp";

const config :Workerd.Config = (
  services = [
    (name = "test1", worker = .vtest1Worker),
  ],

  sockets = [
    (
      name = "test1",
      address = "unix:/test1/test.sock",
      http=(),
      service="test1"
    ),
  ]
);

const vtest1Worker :Workerd.Worker = (
  modules = [
    (name = "test1/entry.js", esModule = embed "src/test1/entry.js"),
  ],
  compatibilityDate = "2023-04-03",
);`, result["test1"])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect(t, BuildCapfile(tt.wokers))
		})
	}
}
