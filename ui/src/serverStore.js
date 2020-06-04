class WSMessage {
  constructor(type, payload) {
    this.messageType = type;
    this.payload = payload;
  }
}

function createLoginRequestVerifier(socket, resolve, reject) {
  let handler = (ev) => {
    socket.removeEventListener("message", handler);
    console.log(ev);
    let data = JSON.parse(ev.data);
    if (data.messageType === "LoginSuccess") {
      resolve({ payload: data.payload, socket });
    } else if (data.messageType === "LoginFailure") {
      console.error(data.payload.message);
      reject("login failed");
    } else {
      reject("Unexpected MessageType" + data.messageType);
    }
  };

  socket.addEventListener("message", handler);
}

function sendLoginRequest(socket, username) {
  let msg = new WSMessage("LoginRequest", {
    loginID: username,
  });
  socket.send(JSON.stringify(msg));
  return new Promise((resolve, reject) => {
    createLoginRequestVerifier(socket, resolve, reject);
  });
}

const webSocketHandler = {
  socket: null,
  store: null,
  installFunc() {
    let handler = this;
    return (store) => {
      handler.store = store;
    };
  },
  setSocket(socket) {
    if (this.socket !== null) {
      // todo: support reconnections? Or just refresh
      console.error("socket being replaced!");
    }
    this.socket = socket;
    let handler = this;
    socket.addEventListener("message", (ev) => {
      handler.handleMessage(ev.data);
    });
  },
  handleMessage(msg) {
    console.debug(msg);
    switch (msg.type) {
      case "GameUpdate":
        this.state.commit("gameUpdate", msg);
        break;
      default:
        console.error("unknown websocket message type: " + msg.type);
        break;
    }
  },
  sendMessage(msg) {
    this.socket.send(msg);
  },
};

const store = {
  state: {
    username: null,
    playerID: "",
    games: [],
  },
  mutations: {
    setUser(state, { username, playerID }) {
      state.username = username;
      state.playerID = playerID;
    },
    initGames(state, games) {
      state.games = games;
    },
    gameUpdate(state, game) {
      let idx = state.games.find((g) => g.id == game.id);
      if (idx == -1) {
        state.games.push(game);
      } else {
        state.games[idx] = game;
      }
    },
  },
  plugins: [webSocketHandler.installFunc()],
  actions: {
    login(context, username) {
      let url = process.env.VUE_APP_SERVER_SOCKET_URL;
      console.log("connecting to " + url);
      let socket = new WebSocket(url);
      return new Promise((resolve, reject) => {
        socket.addEventListener("open", () => {
          sendLoginRequest(socket, username)
            .then(({ payload, socket }) => {
              console.log("successful connection");
              webSocketHandler.setSocket(socket);
              context.commit("setUser", payload);
              resolve();
            })
            .catch((reason) => {
              console.error(reason);
              reject(reason);
            });
        });
      });
    },
    playMove(context, { gameID, move }) {
      let message = new WSMessage("PlayMove", { gameID, move });
      webSocketHandler.sendMessage(message);
    },
    newGame(context, opponent) {
      let message = new WSMessage("NewGame", { opponent });
      webSocketHandler.sendMessage(message);
    },
    lookupOpponent(context, opponent) {
      let message = new WSMessage("LookupByUsername", { username: opponent });
      webSocketHandler.sendMessage(message);
    },
  },
  // getters: {
  //   getUsername(state) {
  //     return state.username;
  //   },
  // },
};

export default store;
