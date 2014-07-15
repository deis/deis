package deis_builder_spec

import (
    . "/path/to/makefile"
    . "github.com/onsi/ginkgo"
)

var _ = Describe("SSH ClientAliveInterval 45", func() {
  Context("When an intermediary is configured at 30 seconds timeout", func() {
    /*
      vagrant_firewall := setup_vagrant_firewall()
      deis_cluster_behind_firewall := setup_deis_cluster_behind_firewall()
      output_of_making_deis_builder_call := run_deis_builder_call()
    */
    It("should timeout with <error>", func() {
      Expect(output_of_making_deis_builder_call).To(Fail)
    })
  })
  Context("When an intermediary is configured at 60 seconds timeout", func() {
    /*
      vagrant_firewall := setup_vagrant_firewall()
      deis_cluster_behind_firewall := setup_deis_cluster_behind_firewall()
      output_of_making_deis_builder_call := run_deis_builder_call()
    */

    It("should be able to run a deis-builder push successfully", func() {
      Expect(output_of_making_deis_builder_call).To(Pass)
    })
  })
})
