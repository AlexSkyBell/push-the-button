import { createStore } from 'vuex'
import message from './message'
import websocket from './websocket'


const state = {}
const getters = {}
const mutations = {}
const actions = {}

const store = createStore({
    state: state,
    getters: getters,
    mutations: mutations,
    actions: actions,
    modules: {
        message: message,
        websocket: websocket
    }
})

export default store

