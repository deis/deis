package deis_controller_views_spec

import (
    . "/path/to/makefile"
    . "github.com/onsi/ginkgo"
)

/*
SETUP:
------
install deis cluster on vagrant
admin_user = register_a_deis_user()
logout
app_user = register_a_deis_user()
first_app = create_a_new_app()
logout
*/

var _ = Describe("AppPermsViewSet.list", func() {
  Context("when a admin-user tries to fetch it for an app he does not own", func() {
    /*
      login as admin_user
      output = `run deis auth:list --app=<app_name>`
      correct_json_response = {'users': usernames}
    */
    It("should allow him access", func() {
      Expect(output.header).To(Eq(200))
    })
    It("should return a json with the values in the format <correct_json_response>", func() {
      Expect(output).To(Contain(correct_json_response))
    })
  })
})

var _ = Describe("AppPermsViewSet.create", func() {
  Context("when a admin-user tries to fetch it for an app he does not own", func() {
    /*
      login as admin_user
      output = `run deis auth:create <new_user> --app=<app_name>`
      correct_json_response = "User <new_user> was granted access to <app_name>"
    */
    It("should allow him access", func() {
      Expect(output.header).To(Eq(200))
    })
    It("should return a json with the values in the format <correct_json_response>", func() {
      Expect(output).To(Contain(correct_json_response))
    })
  })
})

var _ = Describe("AppPermsViewSet.destroy", func() {
  Context("when a admin-user tries to fetch it for an app he does not own", func() {
    /*
      login as admin_user
      output = `run deis auth:destroy <username> --app=<app_name>`
      correct_json_response = User <username> was revoked access to <app_name>
    */
    It("should allow him access", func() {
      Expect(output.header).To(Eq(200))
    })
    It("should return a json with the values in the format <correct_json_response>", func() {
      Expect(output).To(Contain(correct_json_response))
    })
  })
})
