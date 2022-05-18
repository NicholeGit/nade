package command

import (
	"sync"
	"testing"
	"time"

	"github.com/NicholeGit/nade/framework/contract"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func appStartAndStop(t *testing.T) {
	Convey("command app", t, func(c C) {
		var wg sync.WaitGroup
		Convey("app Start", func() {
			wg.Add(1)
			go func(rootCmd *cobra.Command) {
				defer wg.Done()
				rootCmd.ResetFlags()
				rootCmd.SetArgs([]string{"app", "start"})
				_, err := rootCmd.ExecuteContextC(gCtx)
				c.So(err, ShouldBeNil)
			}(gRootCmd)
		})
		Convey("app Stop", func() {
			kernel := gContainer.MustMake(contract.KernelKey).(contract.IKernel)
			app := kernel.GetServer("TestServer")
			testApp, ok := app.(*TestServer)
			So(ok, ShouldBeTrue)

			for !testApp.IsRun() { // 等待服务器启动成功过
				time.Sleep(time.Millisecond)
			}
			for testApp.IsRun() {
				// time.Sleep(time.Second)
				gRootCmd.ResetFlags()
				gRootCmd.SetArgs([]string{"app", "stop"})
				_, err := gRootCmd.ExecuteContextC(gCtx)
				So(err, ShouldBeNil)
				break
			}
			wg.Wait()
			So(testApp.IsRun(), ShouldBeFalse)
		})
	})
}
