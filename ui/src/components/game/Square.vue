<template>
  <div class="square">
    <div
      v-if="square.playable"
      v-on:click="play()"
      class="text-center square playable"
    />

    <!--TODO: Add icons-->
    <div
      v-else-if="square.owner == game.playerX"
      class="text-center square player-x"
    >
      X
    </div>
    <div
      v-else-if="square.owner == game.playerO"
      class="text-center square player-o"
    >
      O
    </div>
    <div v-else class="empty"></div>
  </div>
</template>

<script>
export default {
  name: "Square",
  props: {
    square: Object,
    game: Object,
  },
  computed: {
    move() {
      return {
        gameID: this.game.gameID,
        move: {
          playerID: this.$store.state.playerID,
          coordinate: this.square.coordinate,
        },
      };
    },
  },
  methods: {
    play() {
      this.$store.dispatch("playMove", this.move);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.playable {
  background-color: yellow;
}

.player-x {
  background-color: lightcoral;
}

.player-o {
  background-color: lightskyblue;
}

.square {
  width: 100%;
  height: 100%;
  font-size: x-large;
}
</style>
