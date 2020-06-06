<template>
  <div class="login m-5">
    <h1>Debug Login</h1>
    <label for="username">Username: </label>
    <br />
    <input id="username" v-model="username" />
    <b-button variant="primary" v-on:click="login">Login</b-button>
  </div>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      username: "",
    };
  },
  methods: {
    login: function() {
      this.$store
        .dispatch("login", this.username)
        .then(() => {
          console.log("connection succeeded");
          this.$router.push("/app");
        })
        .catch((ex) => console.error(ex));
    },
  },
};
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
