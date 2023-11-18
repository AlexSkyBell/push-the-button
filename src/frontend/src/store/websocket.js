const DEBUG_WS_SERVER = 'ws://localhost:18080'

const state = {
    websocket: null,
    wsReconnectAttempts: 0,
    established: false,
    reconnecting: false,
}

const getters = {
    getWsReconnectAttempts(state) {
        return state.wsReconnectAttempts
    },
    getIsReconnecting(state) {
        return state.reconnecting
    },
    getWebsocket(state) {
        return state.websocket
    },
    getEstablished(state) {
        return state.established
    }
}

const mutations = {
    increaseWsReconnectAttempts(state) {
        state.wsReconnectAttempts++
    },
    truncateAttempts(state) {
        state.wsReconnectAttempts = 0
    },
    setReconnecting(state, value) {
        state.reconnecting = value
    },
    setWebsocket(state, ws) {
        state.websocket = ws
    },
    setEstablished(state, newState) {
        state.established = newState
    }
}

const actions = {
    connectWebsocket({ commit, dispatch, getters }, vm) {
        const ws = new WebSocket(vm.$appSettings.websocket || DEBUG_WS_SERVER)
        ws.onopen = () => {
            commit('setWebsocket', ws)
            commit('setEstablished', true)
            commit('truncateAttempts')

            dispatch('sendWsMsg', { command: 'hello', payload: 'hello' })
        }
        ws.onclose = () => {
            commit('setWebsocket', null)
            commit('setEstablished', false)
            commit('increaseWsReconnectAttempts')
            if (getters.getWsReconnectAttempts <= 40) {
                setTimeout(function () {
                    dispatch('connectWebsocket', vm);
                }, 3000);
            }
        }
        ws.onerror = () => {
            dispatch('setReconnecting', true)
            ws.close()
        }
        ws.onmessage = (e) => {
            const response = JSON.parse(e.data)
            dispatch('handleResponse', response)
        }
    },
    sendWsMsg({ state, getters }, payload) {
        if (getters.getEstablished) state.websocket.send(JSON.stringify(payload))
    },
    setReconnecting(context, newState) {
        context.commit('setReconnecting', newState)
    },
    handleResponse({ commit, dispatch, getters }, payload) {
        dispatch('prepareMessage', payload)
    }
}

export default {
    state,
    getters,
    mutations,
    actions,
}
