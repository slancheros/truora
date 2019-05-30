<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <b-tabs content-class="mt-3">
      <b-tab title="Find info for a Domain" active>
        <b-container>
          <b-form-group
                  id="fieldset-1"
                  description="Example: google.com"
                  label="Please enter domain's name"
                  label-for="input-horizontal"
                  :invalid-feedback="invalidFeedback"
                  :valid-feedback="validFeedback"
                  :state="state">
            <b-form-row>
              <b-col cols="100"></b-col>
              <b-col >
                <b-form-input id="input-1" v-model="name" :state="state" trim size="20"></b-form-input>
              </b-col>
              <b-col cols="50">  <b-button variant="success" v-on:click="getDomainInfo">Get Info</b-button></b-col>
            </b-form-row>
            <b-form-row>
              <p>

              </p>
            </b-form-row>
            <b-form-row>
              <b-form-textarea
                      id="textarea"
                      v-model="text"
                      placeholder="Domain Info..."
                      rows="30"
                      max-rows="50"
                      plaintext
              ></b-form-textarea>

            </b-form-row>
          </b-form-group>
        </b-container>
      </b-tab>
      <b-tab title="List Queried Domains">
        <b-container>
          <b-form-group
                  id="fieldset-2"
                  description="List of queried domains"
                  label="Queried domain's"
                  label-for="input-horizontal"
                  :state="state">
            <b-form-row>
              <b-col cols="50">  <b-button variant="success" v-on:click="getDomainList">List Domains</b-button></b-col>
            </b-form-row>
            <b-form-row>
              <p>

              </p>
            </b-form-row>
            <b-form-row>

              <b-form-textarea
                      id="listDomains_TA"
                      v-model="domainList"
                      placeholder="Domain Info..."
                      rows="30"
                      max-rows="50"
                      plaintext
              ></b-form-textarea>

            </b-form-row>
          </b-form-group>
        </b-container>
      </b-tab>
    </b-tabs>

  </div>

</template>

<script>
  import axios from 'axios';
export default {
  computed: {
    state() {
      return this.name.length >= 4 ? true : false
    },
    invalidFeedback() {
      if (this.name.length > 4) {
        return ''
      } else if (this.name.length > 0) {
        return 'Enter at least 4 characters'
      } else {
        return 'The domain\'s name is still empty'
      }
    },
    validFeedback() {
      return this.state === true ? 'Thank you' : ''
    }
  },
  name: 'DomainInfo',
  props: {
    msg: String,
    info:String
  },
  data: function () {
    return {
      text: '',
      domainList: '',
      name: '',
      isLoading: false,
    }
  },
  methods: {
    getDomainInfo: function () {
       axios.get("http://192.168.1.30:3344/serverInfo/" + this.name)
              .then(response =>
                this.text = JSON.stringify(response.data)
              ).catch(function (error) {
                this.text = "Error: " + error.data
              })

    },
    getDomainList: function () {
      return axios.get("http://192.168.1.30:3344/serverInfo/list")
              .then(response =>
                  this.domainList = JSON.stringify(response.data)
              ).catch(error =>
                  this.domainList = "Error: " + error.data
              )
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
