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
        this.text =''
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
    msg: String
  },
     data() {
      return {
        text: '',
        domainList:'',
        name: ''
      }
    },
  methods: {
    getDomainInfo: function () {
      axios.get("http://localhost:9898/serverInfo/" + this.name)
              .then((response) => {
                this.text = "Success:"+ response.data.value;
              })
              .catch((error) => {
                  if (error.response) {
                    this.text = "Error: "+error.response.data + error.response.status + error.response.headers
                    // The request was made and the server responded with a status code
                    // that falls out of the range of 2xx
                    // console.log(error.response.data);
                    // console.log(error.response.status);
                    // console.log(error.response.headers);
                  } else if (error.request) {
                    // The request was made but no response was received
                    // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
                    // http.ClientRequest in node.js
                    this.text = "Error: "+error.request.data.value;
                  } else {
                    // Something happened in setting up the request that triggered an Error
                    this.text="Error: "+ error.message;
                  }
                })

    },

    getDomainList: function () {
      this.domainList = ''
      axios.get("http://localhost:9898/serverInfo/list")
              .then((response => this.domainList = response.data.text)
                      .catch( error => this.domainList = error.response.data))
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
