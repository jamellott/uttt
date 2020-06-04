<template>
  <div class="login">
    <input v-model="username" />
    <button v-on:click="login">Login</button>
  </div>
</template>

<script>
export default {
  name: "NewGameMenu",
  data() {
    return {
      username: "",
    };
  },
  methods: {
    login: function() {
      this.$store
        .dispatch("login", this.username)
        .then(() => console.log("connection succeeded"))
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
