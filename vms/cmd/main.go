// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"k8s.io/client-go/informers"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	istiocmd "istio.io/istio/pkg/cmd"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/security/pkg/cmd"
	"istio.io/istio/vms/pkg/k8s/vmcontrol"
	"istio.io/pkg/collateral"
	"istio.io/pkg/ctrlz"
	"istio.io/pkg/log"
	"istio.io/pkg/version"
)

const ()

var (
	opts = cliOptions{
		logOptions:   log.DefaultOptions(),
		ctrlzOptions: ctrlz.DefaultOptions(),
	}

	rootCmd = &cobra.Command{
		Use:   "Istio VM Controller",
		Short: "Istio VM Controller",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if !opts.enableController {
				log.Info("Istio VM Controller is not enabled, exiting")
				return
			}
			runController()
		},
	}
)

type cliOptions struct {
	kubeConfigFile string

	logOptions   *log.Options
	ctrlzOptions *ctrlz.Options

	// Whether to enable the VM controller
	enableController bool
}

func init() {
	flags := rootCmd.Flags()

	flags.BoolVar(&opts.enableController, "enable-controller", false, "Specifies whether to enable "+
		"the Istio VM Controller.")

	flags.StringVar(&opts.kubeConfigFile, "kube-config", "",
		"Specifies path to kubeconfig file. This must be specified when not running inside a Kubernetes pod.")

	rootCmd.AddCommand(version.CobraCommand())
	rootCmd.AddCommand(collateral.CobraCommand(rootCmd, &doc.GenManHeader{
		Title:   "VMControl: Istio VM Controller",
		Section: "VMControl: Istio VM Controller",
		Manual:  "VMControl: Istio VM Controller",
	}))
	rootCmd.AddCommand(cmd.NewProbeCmd())

	opts.logOptions.AttachCobraFlags(rootCmd)
	opts.ctrlzOptions.AttachCobraFlags(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(1)
	}
}

func runController() {
	_, _ = ctrlz.Run(opts.ctrlzOptions, nil)

	if err := log.Configure(opts.logOptions); err != nil {
		log.Errorf("failed to configure logging (%v)", err)
		os.Exit(1)
	}

	k8sClient, err := kube.CreateClientset(opts.kubeConfigFile, "")
	if err != nil {
		log.Errorf("could not create k8s clientset: %v", err)
		os.Exit(1)
	}

	stopCh := make(chan struct{})

	sharedInformers := informers.NewSharedInformerFactory(k8sClient, 5*time.Second)

	c, err := vmcontrol.NewController(sharedInformers)

	if err != nil {
		log.Errorf("failed to create VM controller: %v", err)
		os.Exit(1)
	}

	// Run the controller to manage VMs added to the mesh.
	c.Run(stopCh)

	istiocmd.WaitSignal(stopCh)
}
