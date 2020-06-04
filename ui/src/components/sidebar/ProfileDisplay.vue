<template>
  <div class="profile-display">
    {{ username }}
  </div>
</template>

<script>
export default {
  name: "ProfileDisplay",
  methods: {
    login: function() {
      this.$store
        .dispatch("login", this.username)
        .then(() => console.log("connection succeeded"))
        .catch((ex) => console.error(ex));
    },
  },
  computed: {
    username() {
      return this.$store.state.username;
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
