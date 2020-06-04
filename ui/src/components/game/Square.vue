<template>
  <div class="square">
    <div v-if="square.playable" v-on:click="play()" class="square playable" />

    <!--TODO: Add icons-->
    <div v-else-if="square.owner == game.playerX" class="player-x"></div>
    <div v-else-if="square.owner == game.playerO" class="player-o"></div>
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
  background-color: red;
}

.player-o {
  background-color: blue;
}

.square {
  width: 100%;
  height: 100%;
}
</style>
