<template>
  <div class="p-t-2">
    <b-card class="mt-5">
      <b-card-text>
        <b-form>
          <b-form-group label="Opponent Username" label-for="opponent">
            <b-form-input
              id="opponent"
              v-model="opponentUsername"
              :state="validated"
              required
              v-on:change="validateOpponent()"
            >
            </b-form-input>
            <!--<b-form-text v-if="isValidating"
              >Checking that user exists...</b-form-text
            >-->
            <b-form-invalid-feedback :state="!invalidated"
              >User does not exist.</b-form-invalid-feedback
            >
          </b-form-group>
          <b-form-group>
            <b-button
              :disabled="!validated"
              v-on:click="startGame()"
              variant="primary"
              >Start Game</b-button
            >
          </b-form-group>
        </b-form>
      </b-card-text>
    </b-card>
  </div>
</template>

<script>
export default {
  name: "NewGameButton",
  data() {
    return {
      opponentUsername: "",
      opponentUUID: null,
      isValidating: false,
    };
  },
  computed: {
    validated() {
      return !this.isValidating && this.opponentUUID != null;
    },
    invalidated() {
      return !this.isValidating && this.opponentUUID == null;
    },
  },
  methods: {
    validateOpponent() {
      this.isValidating = true;
      this.$store.dispatch("lookupOpponent", this.opponentUsername).then(() => {
        this.opponentUUID = this.$store.state.usernameMap[
          this.opponentUsername
        ];
        this.isValidating = false;
      });
    },
    startGame() {
      if (!this.validated) {
        return;
      }

      this.$store.dispatch("newGame", this.opponentUUID);
      this.opponentUsername = "";
      this.opponentUUID = null;
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
